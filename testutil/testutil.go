package testutil

import (
	"runtime"
	"testing"
)

func RedirectConfigDir(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()

	switch runtime.GOOS {
	case "windows":
		t.Setenv("AppData", tmpDir)
	case "darwin":
		t.Setenv("HOME", tmpDir)
	default:
		t.Setenv("XDG_CONFIG_HOME", tmpDir)
	}
}
