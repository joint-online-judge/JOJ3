package conf

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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

func TestGetConfPath(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Setup layout: tmpDir/alpha/conf.json, tmpDir/conf.json
	alphaDir := filepath.Join(tmpDir, "alpha")
	if err := os.MkdirAll(alphaDir, 0755); err != nil {
		t.Fatalf("failed to create alpha dir: %v", err)
	}
	alphaConfPath := filepath.Join(alphaDir, "conf.json")
	if err := os.WriteFile(alphaConfPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write alpha conf: %v", err)
	}

	rootConfPath := filepath.Join(tmpDir, "conf.json")
	if err := os.WriteFile(rootConfPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write root conf: %v", err)
	}

	t.Run("Exact Scope Match", func(t *testing.T) {
		gotPath, _, cc, err := GetConfPath(tmpDir, "conf.json", "conf.json", "", "alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotPath != alphaConfPath {
			t.Errorf("got path %s, want %s", gotPath, alphaConfPath)
		}
		if cc.Scope != "alpha" {
			t.Errorf("got scope %s, want alpha", cc.Scope)
		}
	})

	t.Run("Prefix Scope Match for Floating Release", func(t *testing.T) {
		gotPath, _, cc, err := GetConfPath(tmpDir, "conf.json", "conf.json", "", "alpha1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotPath != alphaConfPath {
			t.Errorf("got path %s, want %s", gotPath, alphaConfPath)
		}
		if cc.Scope != "alpha" {
			t.Errorf("got scope %s, want alpha", cc.Scope)
		}
	})

	t.Run("Fallback to Root Conf when Scope Not Found", func(t *testing.T) {
		gotPath, _, cc, err := GetConfPath(tmpDir, "conf.json", "conf.json", "", "beta")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotPath != rootConfPath {
			t.Errorf("got path %s, want %s", gotPath, rootConfPath)
		}
		if cc.Scope != "beta" {
			t.Errorf("got scope %s, want beta", cc.Scope)
		}
	})

	t.Run("Fallback to Root Conf with Empty Fallback Conf Name", func(t *testing.T) {
		gotPath, _, _, err := GetConfPath(tmpDir, "conf.json", "", "", "gamma")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotPath != rootConfPath {
			t.Errorf("got path %s, want %s", gotPath, rootConfPath)
		}
	})

	t.Run("No Conf Found Error", func(t *testing.T) {
		emptyTmpDir := t.TempDir()
		_, _, _, err := GetConfPath(emptyTmpDir, "conf.json", "conf.json", "", "alpha")
		if err == nil {
			t.Fatalf("expected error when no conf found, got nil")
		}
	})
}

