package healthcheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func newCommitRepo(t *testing.T, message, email string) string {
	t.Helper()
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatal(err)
	}
	signature := &object.Signature{Name: "Student", Email: email, When: time.Unix(1, 0)}
	if _, err := worktree.Commit(message, &git.CommitOptions{Author: signature, Committer: signature}); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestCommitChecks(t *testing.T) {
	valid := newCommitRepo(t, "feat: valid", "student@example.edu")
	if err := NonASCIIMsg(valid); err != nil {
		t.Fatal(err)
	}
	if err := AuthorEmailCheck(valid, []string{"example.edu"}, filepath.Join(valid, "missing.csv")); err != nil {
		t.Fatal(err)
	}

	invalidMessage := newCommitRepo(t, "feat: 测试", "student@invalid.test")
	if err := NonASCIIMsg(invalidMessage); err == nil || !strings.Contains(err.Error(), "测试") {
		t.Fatalf("NonASCIIMsg() error = %v", err)
	}
	if err := AuthorEmailCheck(invalidMessage, []string{"example.edu"}, filepath.Join(invalidMessage, "missing.csv")); err == nil || !strings.Contains(err.Error(), "allowed domains") {
		t.Fatalf("AuthorEmailCheck(domain) error = %v", err)
	}
}

func TestAuthorEmailCheckActorCSV(t *testing.T) {
	repo := newCommitRepo(t, "feat: valid", "student@example.edu")
	csvPath := filepath.Join(t.TempDir(), "actors.csv")
	if err := os.WriteFile(csvPath, []byte("name,id,student\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := AuthorEmailCheck(repo, []string{"example.edu"}, csvPath); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(csvPath, []byte("name,id,other\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := AuthorEmailCheck(repo, []string{"example.edu"}, csvPath); err == nil || !strings.Contains(err.Error(), "not stored") {
		t.Fatalf("AuthorEmailCheck(actor) error = %v", err)
	}
}

func TestCommitChecksRejectNonRepository(t *testing.T) {
	if err := NonASCIIMsg(t.TempDir()); err == nil || !strings.Contains(err.Error(), "opening git repo") {
		t.Fatalf("NonASCIIMsg() error = %v", err)
	}
	if err := AuthorEmailCheck(t.TempDir(), nil, ""); err == nil || !strings.Contains(err.Error(), "opening git repo") {
		t.Fatalf("AuthorEmailCheck() error = %v", err)
	}
}
