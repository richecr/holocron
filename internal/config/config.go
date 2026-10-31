package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Version  string                  `toml:"version"`
	Settings Settings                `toml:"settings"`
	Scripts  map[string]ScriptConfig `toml:"scripts"`
}

type Settings struct {
	ConfirmBeforeRun bool              `toml:"confirm_before_run"`
	DefaultShell     map[string]string `toml:"default_shell"`
}

type ScriptConfig struct {
	Description string            `toml:"description,omitempty"`
	Command     string            `toml:"command,omitempty"`
	Tags        []string          `toml:"tags,omitempty"`
	Platforms   map[string]string `toml:"platforms,omitempty"`
	Env         map[string]string `toml:"env,omitempty"`
	Dir         string            `toml:"dir,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Version: "1.0.0",
		Settings: Settings{
			ConfirmBeforeRun: false,
			DefaultShell:     map[string]string{},
		},
		Scripts: map[string]ScriptConfig{},
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			if err := SaveConfig(path, cfg); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return Decode(data)
}

func Decode(data []byte) (*Config, error) {
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Settings.DefaultShell == nil {
		cfg.Settings.DefaultShell = map[string]string{}
	}
	if cfg.Scripts == nil {
		cfg.Scripts = map[string]ScriptConfig{}
	}

	return cfg, nil
}

func Encode(cfg *Config) ([]byte, error) {
	return toml.Marshal(cfg)
}

func SaveConfig(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := Encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
