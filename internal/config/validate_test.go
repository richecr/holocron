package config

import "testing"

func TestValidateCleanConfig(t *testing.T) {
	cfg := &Config{
		Settings: Settings{
			DefaultShell: map[string]string{"linux": "bash"},
		},
		Scripts: map[string]ScriptConfig{
			"docker:up": {Command: "docker compose up -d", Tags: []string{"docker"}},
			"clean": {
				Platforms: map[string]string{
					"linux":   "rm -rf dist/",
					"windows": "rmdir /s /q dist",
				},
			},
		},
	}

	if issues := Validate(cfg); len(issues) != 0 {
		t.Errorf("Validate() = %v, want no issues", issues)
	}
}

func TestValidateInvalidShell(t *testing.T) {
	cfg := &Config{
		Settings: Settings{
			DefaultShell: map[string]string{"linux": "fish"},
		},
		Scripts: map[string]ScriptConfig{},
	}

	issues := Validate(cfg)
	if len(issues) != 1 {
		t.Fatalf("Validate() = %v, want 1 issue", issues)
	}
	if issues[0].Severity != "error" {
		t.Errorf("severity = %q, want error", issues[0].Severity)
	}
}

func TestValidateScriptWithoutCommand(t *testing.T) {
	cfg := &Config{
		Scripts: map[string]ScriptConfig{
			"broken": {Description: "no command at all"},
		},
	}

	issues := Validate(cfg)
	if len(issues) != 1 {
		t.Fatalf("Validate() = %v, want 1 issue", issues)
	}
	if issues[0].Severity != "error" {
		t.Errorf("severity = %q, want error", issues[0].Severity)
	}
}

func TestValidateEmptyPlatformCommand(t *testing.T) {
	cfg := &Config{
		Scripts: map[string]ScriptConfig{
			"clean": {Platforms: map[string]string{"linux": ""}},
		},
	}

	issues := Validate(cfg)
	if len(issues) != 1 {
		t.Fatalf("Validate() = %v, want 1 issue", issues)
	}
	if issues[0].Severity != "error" {
		t.Errorf("severity = %q, want error", issues[0].Severity)
	}
}

func TestValidateUnknownPlatformIsWarning(t *testing.T) {
	cfg := &Config{
		Scripts: map[string]ScriptConfig{
			"clean": {Platforms: map[string]string{"windwos": "rmdir dist"}},
		},
	}

	issues := Validate(cfg)
	if len(issues) != 1 {
		t.Fatalf("Validate() = %v, want 1 issue", issues)
	}
	if issues[0].Severity != "warning" {
		t.Errorf("severity = %q, want warning", issues[0].Severity)
	}
}
