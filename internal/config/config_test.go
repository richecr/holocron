package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigCreatesDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.toml")

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to be created: %v", err)
	}

	if cfg.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", cfg.Version, "1.0.0")
	}
	if cfg.Settings.ConfirmBeforeRun {
		t.Error("ConfirmBeforeRun = true, want false")
	}
	if len(cfg.Settings.DefaultShell) != 0 {
		t.Errorf("DefaultShell = %v, want empty", cfg.Settings.DefaultShell)
	}
	if len(cfg.Scripts) != 0 {
		t.Errorf("Scripts = %v, want empty", cfg.Scripts)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")

	want := &Config{
		Version: "1.0.0",
		Settings: Settings{
			ConfirmBeforeRun: true,
			DefaultShell:     map[string]string{"linux": "bash"},
		},
		Scripts: map[string]ScriptConfig{
			"build:java": {
				Description: "Compiles the Java project",
				Command:     "./mvnw clean package",
				Tags:        []string{"java", "build"},
			},
			"docker:up": {
				Description: "Starts core services",
				Platforms: map[string]string{
					"linux":   "docker compose up -d",
					"windows": "docker compose up -d",
				},
				Tags: []string{"docker"},
			},
		},
	}

	if err := SaveConfig(path, want); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	got, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if got.Version != want.Version {
		t.Errorf("Version = %q, want %q", got.Version, want.Version)
	}
	if got.Settings.ConfirmBeforeRun != want.Settings.ConfirmBeforeRun {
		t.Errorf("ConfirmBeforeRun = %v, want %v", got.Settings.ConfirmBeforeRun, want.Settings.ConfirmBeforeRun)
	}
	if got.Settings.DefaultShell["linux"] != "bash" {
		t.Errorf("DefaultShell[linux] = %q, want %q", got.Settings.DefaultShell["linux"], "bash")
	}

	build, ok := got.Scripts["build:java"]
	if !ok {
		t.Fatal("missing script build:java")
	}
	if build.Command != "./mvnw clean package" {
		t.Errorf("Command = %q, want %q", build.Command, "./mvnw clean package")
	}
	if len(build.Tags) != 2 || build.Tags[0] != "java" {
		t.Errorf("Tags = %v, want [java build]", build.Tags)
	}

	docker, ok := got.Scripts["docker:up"]
	if !ok {
		t.Fatal("missing script docker:up")
	}
	if docker.Platforms["linux"] != "docker compose up -d" {
		t.Errorf("Platforms[linux] = %q", docker.Platforms["linux"])
	}
	if docker.Platforms["windows"] != "docker compose up -d" {
		t.Errorf("Platforms[windows] = %q", docker.Platforms["windows"])
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}
	if path == "" {
		t.Fatal("GetConfigPath() returned empty path")
	}
	if filepath.Base(path) != "config.toml" {
		t.Errorf("Base(path) = %q, want config.toml", filepath.Base(path))
	}
	if filepath.Base(filepath.Dir(path)) != "holocron" {
		t.Errorf("parent dir = %q, want holocron", filepath.Base(filepath.Dir(path)))
	}
}
