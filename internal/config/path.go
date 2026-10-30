package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigPath() (string, error) {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "holocron", "config.toml"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "holocron", "config.toml"), nil
		}
		return filepath.Join(home, "AppData", "Roaming", "holocron", "config.toml"), nil
	}

	return filepath.Join(home, ".config", "holocron", "config.toml"), nil
}
