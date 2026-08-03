package conf

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestConfLogValueRedactsSandboxToken(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	conf := &Conf{
		Name:         "test",
		SandboxToken: "top-secret-token",
	}
	logger.Info("config", "conf", conf)
	got := output.String()
	if strings.Contains(got, "top-secret-token") {
		t.Fatalf("configuration log exposed sandbox token: %s", got)
	}
	if !strings.Contains(got, "[REDACTED]") || !strings.Contains(got, `"Name":"test"`) {
		t.Fatalf("configuration log lost expected diagnostic fields: %s", got)
	}
	if conf.SandboxToken != "top-secret-token" {
		t.Fatalf("logging mutated the runtime configuration: %q", conf.SandboxToken)
	}
}

func TestGetSHA256(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input")
	content := []byte("joj3")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	wantSum := sha256.Sum256(content)
	got, err := GetSHA256(path)
	if err != nil || got != hex.EncodeToString(wantSum[:]) {
		t.Fatalf("GetSHA256() = %q, %v", got, err)
	}
	if _, err := GetSHA256(path + ".missing"); !os.IsNotExist(err) {
		t.Fatalf("missing GetSHA256() error = %v", err)
	}
}

func TestGetConfPath(t *testing.T) {
	root := t.TempDir()
	scopedDir := filepath.Join(root, "course")
	if err := os.Mkdir(scopedDir, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{filepath.Join(scopedDir, "conf.json"), filepath.Join(root, "fallback.json")} {
		if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	got, _, commit, err := GetConfPath(root, "conf.json", "fallback.json", "test(course): run", "")
	if err != nil || got != filepath.Join(scopedDir, "conf.json") || commit.Scope != "course" {
		t.Fatalf("GetConfPath(scoped) = %q, %+v, %v", got, commit, err)
	}
	got, _, _, err = GetConfPath(root, "conf.json", "fallback.json", "not conventional", "")
	if err != nil || got != filepath.Join(root, "fallback.json") {
		t.Fatalf("GetConfPath(fallback) = %q, %v", got, err)
	}
	if _, _, err = parseMsg(root, "conf.json", "test(../escape): run", ""); err == nil || !strings.Contains(err.Error(), "invalid scope") {
		t.Fatalf("parseMsg(traversal) error = %v", err)
	}
	got, _, commit, err = GetConfPath(root, "conf.json", "fallback.json", "ignored", "missing-tag")
	if !os.IsNotExist(err) || got != filepath.Join(root, "missing-tag", "conf.json") || commit.Scope != "missing-tag" {
		t.Fatalf("GetConfPath(tag) = %q, %+v, %v", got, commit, err)
	}
	if err := os.Remove(filepath.Join(root, "fallback.json")); err != nil {
		t.Fatal(err)
	}
	got, _, _, err = GetConfPath(root, "conf.json", "fallback.json", "invalid", "")
	if !os.IsNotExist(err) || got != filepath.Join(root, "fallback.json") {
		t.Fatalf("GetConfPath(missing fallback) = %q, %v", got, err)
	}
}

func TestParseConventionalCommit(t *testing.T) {
	tests := []struct {
		name    string
		commit  string
		want    *ConventionalCommit
		wantErr bool
	}{
		{
			name:   "Simple feat commit",
			commit: "feat: add new feature",
			want: &ConventionalCommit{
				Type:        "feat",
				Description: "add new feature",
			},
			wantErr: false,
		},
		{
			name:   "Commit with scope",
			commit: "fix(core): resolve memory leak",
			want: &ConventionalCommit{
				Type:        "fix",
				Scope:       "core",
				Description: "resolve memory leak",
			},
			wantErr: false,
		},
		{
			name:   "Breaking change commit",
			commit: "feat(api)!: redesign user authentication",
			want: &ConventionalCommit{
				Type:        "feat",
				Scope:       "api",
				Description: "redesign user authentication",
			},
			wantErr: false,
		},
		{
			name:   "Commit with body",
			commit: "docs: update README\n\nAdd installation instructions and improve examples",
			want: &ConventionalCommit{
				Type:        "docs",
				Description: "update README",
				Body:        "Add installation instructions and improve examples",
			},
			wantErr: false,
		},
		{
			name:   "Commit with body and group",
			commit: "docs: update README [group]\n\nAdd installation instructions and improve examples",
			want: &ConventionalCommit{
				Type:        "docs",
				Description: "update README [group]",
				Group:       "group",
				Body:        "Add installation instructions and improve examples",
			},
			wantErr: false,
		},
		{
			name:   "Commit with body and empty group",
			commit: "docs: update README []\n\nAdd installation instructions and improve examples",
			want: &ConventionalCommit{
				Type:        "docs",
				Description: "update README []",
				Group:       "",
				Body:        "Add installation instructions and improve examples",
			},
			wantErr: false,
		},
		{
			name:   "Full commit with body and footer",
			commit: "feat(auth)!: implement OAuth2\n\nThis commit adds OAuth2 support to the authentication system.\n\nBREAKING CHANGE: Previous authentication tokens are no longer valid.",
			want: &ConventionalCommit{
				Type:        "feat",
				Scope:       "auth",
				Description: "implement OAuth2",
				Body:        "This commit adds OAuth2 support to the authentication system.",
				Footer:      "BREAKING CHANGE: Previous authentication tokens are no longer valid.",
			},
			wantErr: false,
		},
		{
			name:    "Invalid commit format",
			commit:  "This is not a valid conventional commit",
			want:    &ConventionalCommit{},
			wantErr: true,
		},
		{
			name:   "Multi-line body",
			commit: "feat(h1/e2): group (#86)\n\nReviewed-on: https://focs.gc.sjtu.edu.cn/git/test/test/pulls/86\nReviewed-by: foo <foo@sjtu.edu.cn>\nReviewed-by: bar <bar@sjtu.edu.cn>\nReviewed-by: nobody <nobody@sjtu.edu.cn>\n",
			want: &ConventionalCommit{
				Type:        "feat",
				Scope:       "h1/e2",
				Description: "group (#86)",
				Body:        "Reviewed-on: https://focs.gc.sjtu.edu.cn/git/test/test/pulls/86\nReviewed-by: foo <foo@sjtu.edu.cn>\nReviewed-by: bar <bar@sjtu.edu.cn>\nReviewed-by: nobody <nobody@sjtu.edu.cn>",
				Footer:      "",
			},
			wantErr: false,
		},
		{
			name:   "Multi-line body with footer",
			commit: "feat(h1/e2): group (#86)\n\nReviewed-on: https://focs.gc.sjtu.edu.cn/git/test/test/pulls/86\nReviewed-by: foo <foo@sjtu.edu.cn>\nReviewed-by: bar <bar@sjtu.edu.cn>\nReviewed-by: nobody <nobody@sjtu.edu.cn>\n\nFooter here\n",
			want: &ConventionalCommit{
				Type:        "feat",
				Scope:       "h1/e2",
				Description: "group (#86)",
				Body:        "Reviewed-on: https://focs.gc.sjtu.edu.cn/git/test/test/pulls/86\nReviewed-by: foo <foo@sjtu.edu.cn>\nReviewed-by: bar <bar@sjtu.edu.cn>\nReviewed-by: nobody <nobody@sjtu.edu.cn>",
				Footer:      "Footer here",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConventionalCommit(tt.commit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConventionalCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConventionalCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchGroupsUsesExactTokens(t *testing.T) {
	tests := []struct {
		name  string
		group string
		want  []string
	}{
		{name: "exact not substring", group: "cpp", want: []string{"cpp"}},
		{name: "case insensitive", group: "CPP", want: []string{"cpp"}},
		{name: "comma and space", group: "cpp, lint", want: []string{"cpp", "lint"}},
		{name: "semicolon", group: "cpp;lint", want: []string{"cpp", "lint"}},
		{name: "pipe", group: "cpp|lint", want: []string{"cpp", "lint"}},
		{name: "tab", group: "cpp\tlint", want: []string{"cpp", "lint"}},
		{name: "duplicate token", group: "cpp,cpp", want: []string{"cpp"}},
		{name: "all", group: "ALL", want: []string{"c", "cpp", "lint"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Conf{Stages: []ConfStage{
				{Name: "short", Groups: []string{"c"}},
				{Name: "cpp", Groups: []string{"cpp"}},
				{Name: "lint", Groups: []string{"lint"}},
			}}
			got := MatchGroups(conf, &ConventionalCommit{Group: tt.group})
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("MatchGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}
