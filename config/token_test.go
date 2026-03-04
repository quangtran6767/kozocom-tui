package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/quangtran6767/kozocom-tui/testutil"
)

// tokenFilePath returns the actual token file path after redirect.
// Must be called AFTER redirectConfigDir.
func tokenFilePath(t *testing.T) string {
	t.Helper()
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir: %v", err)
	}
	return filepath.Join(configDir, appDirName, tokenFile)
}

// writeTokenFile manually writes a token file for test setup.
func writeTokenFile(t *testing.T, content string) {
	t.Helper()
	path := tokenFilePath(t)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
}

// -- Tests ---
func TestLoadToken_FileNotExist(t *testing.T) {
	testutil.RedirectConfigDir(t)

	token, err := LoadToken()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if token != "" {
		t.Fatalf("expected empty token, got: %q", token)
	}
}

func TestLoadToken_ReturnsToken(t *testing.T) {
	testutil.RedirectConfigDir(t)
	writeTokenFile(t, "my-secrect-token")

	token, err := LoadToken()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if token != "my-secrect-token" {
		t.Fatalf("expected token 'my-secrect-token', got: %q", token)
	}
}

func TestSaveToken_CreatesDirectoryIfNotExist(t *testing.T) {
	testutil.RedirectConfigDir(t)
	err := SaveToken("new-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	path := tokenFilePath(t)
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		t.Fatalf("expected directory to be created, but it doesn't exist")
	}
}

func TestSaveToken_WritesCorrectContent(t *testing.T) {
	testutil.RedirectConfigDir(t)
	err := SaveToken("super-secret")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	content, err := os.ReadFile(tokenFilePath(t))
	if err != nil {
		t.Fatalf("failed to read token file: %v", err)
	}
	if string(content) != "super-secret" {
		t.Fatalf("expected %q, got: %q", "super-secret", string(content))
	}
}

func TestSaveToken_FilePermission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permission check not applicable on Windows")
	}
	testutil.RedirectConfigDir(t)
	_ = SaveToken("some-token")
	info, err := os.Stat(tokenFilePath(t))
	if err != nil {
		t.Fatalf("failed to stat token file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected permission 0600, got: %o", info.Mode().Perm())
	}
}

func TestSaveAndLoadToken(t *testing.T) {
	testutil.RedirectConfigDir(t)
	want := "round-trip-token"
	if err := SaveToken(want); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}
	got, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken failed: %v", err)
	}
	if got != want {
		t.Fatalf("round-trip mismatch: got %q, want %q", got, want)
	}
}
