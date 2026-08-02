package healthcheck

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseWhitelistedChars(t *testing.T) {
	got := parseWhitelistedChars("你, 好, a, invalid, ,你")
	if len(got) != 2 {
		t.Fatalf("parseWhitelistedChars() = %v", got)
	}
	if _, ok := got['你']; !ok {
		t.Fatal("missing whitelisted rune")
	}
}

func TestNonASCIIFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "source.txt"), []byte("hello 你\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := NonASCIIFiles(dir, "你"); err != nil {
		t.Fatalf("whitelisted NonASCIIFiles() error = %v", err)
	}
	if err := NonASCIIFiles(dir, ""); err == nil || !strings.Contains(err.Error(), "source.txt") {
		t.Fatalf("NonASCIIFiles() error = %v", err)
	}
}

func TestForbiddenCheck(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.out\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "result.out"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ForbiddenCheck("."); err == nil || !strings.Contains(err.Error(), "result.out") {
		t.Fatalf("ForbiddenCheck() error = %v", err)
	}
}

func TestVerifyFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "protected.txt")
	content := []byte("original")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	checksum := hex.EncodeToString(sum[:])
	if err := VerifyFiles(dir, "protected.txt", checksum); err != nil {
		t.Fatal(err)
	}
	if err := VerifyFiles(dir, "protected.txt", "bad"); err == nil || !strings.Contains(err.Error(), "altered") {
		t.Fatalf("VerifyFiles(altered) error = %v", err)
	}
	if err := VerifyFiles(dir, "one,two", checksum); err == nil || !strings.Contains(err.Error(), "do not match") {
		t.Fatalf("VerifyFiles(mismatch) error = %v", err)
	}
}
