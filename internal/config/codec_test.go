package config

import "testing"

func TestEncodeDecodeRoundTrip(t *testing.T) {
	original := &Config{
		Version: "1.0.0",
		Settings: Settings{
			ConfirmBeforeRun: true,
			DefaultShell:     map[string]string{"linux": "bash"},
		},
		Scripts: map[string]ScriptConfig{
			"build:java": {
				Description: "Build Java",
				Tags:        []string{"java", "build"},
				Platforms:   map[string]string{"linux": "./mvnw package"},
				Env:         map[string]string{"DEBUG": "1"},
				Dir:         "./backend",
			},
		},
	}

	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	decoded, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	got := decoded.Scripts["build:java"]
	want := original.Scripts["build:java"]

	if got.Dir != want.Dir {
		t.Errorf("Dir = %q, want %q", got.Dir, want.Dir)
	}
	if got.Env["DEBUG"] != "1" {
		t.Errorf("Env = %v, want DEBUG=1", got.Env)
	}
	if got.Command != want.Command {
		t.Errorf("Command = %q, want %q", got.Command, want.Command)
	}
	if decoded.Settings.DefaultShell["linux"] != "bash" {
		t.Errorf("DefaultShell = %v", decoded.Settings.DefaultShell)
	}
	if !decoded.Settings.ConfirmBeforeRun {
		t.Error("ConfirmBeforeRun not preserved")
	}
}

func TestDecodeNormalizesNilMaps(t *testing.T) {
	cfg, err := Decode([]byte(`version = "1.0.0"`))
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if cfg.Scripts == nil {
		t.Error("Scripts is nil")
	}
	if cfg.Settings.DefaultShell == nil {
		t.Error("DefaultShell is nil")
	}
}

func TestDecodeInvalidTOML(t *testing.T) {
	if _, err := Decode([]byte(`this is not toml =`)); err == nil {
		t.Fatal("Decode() = nil error, want error")
	}
}

func TestIsKnownPlatform(t *testing.T) {
	for _, p := range []string{"linux", "darwin", "windows"} {
		if !IsKnownPlatform(p) {
			t.Errorf("IsKnownPlatform(%q) = false, want true", p)
		}
	}
	for _, p := range []string{"windwos", "freebsd", ""} {
		if IsKnownPlatform(p) {
			t.Errorf("IsKnownPlatform(%q) = true, want false", p)
		}
	}
}
