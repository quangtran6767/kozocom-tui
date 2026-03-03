package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	appDirName = "kozocom-tui"
	tokenFile  = "token"
	BaseURL    = "http://localhost:8000"
)

// Build tokenpath, os independent
func tokenPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, appDirName, tokenFile), nil
}

// Load token from tokenPath
func LoadToken() (string, error) {
	path, err := tokenPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// Save token to tokenPath
func SaveToken(token string) error {
	path, err := tokenPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(token), 0600)
}
