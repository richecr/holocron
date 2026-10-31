package ui

import (
	"strings"
	"testing"

	"github.com/richecr/holocron/internal/config"
)

func TestPlatformBadges(t *testing.T) {
	tests := []struct {
		name   string
		script config.ScriptConfig
		want   []string
	}{
		{
			name:   "generic command shows any badge",
			script: config.ScriptConfig{Command: "echo hi"},
			want:   []string{"any"},
		},
		{
			name: "platform map shows per-OS badges in order",
			script: config.ScriptConfig{
				Platforms: map[string]string{
					"windows": "cls",
					"linux":   "clear",
					"darwin":  "clear",
				},
			},
			want: []string{"linux", "darwin", "win"},
		},
		{
			name: "single platform shows only its badge",
			script: config.ScriptConfig{
				Platforms: map[string]string{"linux": "htop"},
			},
			want: []string{"linux"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PlatformBadges(tt.script)

			if tt.name == "platform map shows per-OS badges in order" {
				li := strings.Index(got, "linux")
				da := strings.Index(got, "darwin")
				wi := strings.Index(got, "win")
				if li >= da || da >= wi {
					t.Errorf("badge order wrong: %q", got)
				}
			}
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("PlatformBadges() = %q, missing %q", got, w)
				}
			}
		})
	}
}

func TestRenderScriptTable(t *testing.T) {
	scripts := map[string]config.ScriptConfig{
		"docker:up": {Command: "docker compose up -d", Description: "Docker up", Tags: []string{"docker"}},
		"build:java": {
			Description: "Build Java",
			Platforms:   map[string]string{"linux": "./mvnw package"},
		},
	}

	out := RenderScriptTable(scripts)

	for _, want := range []string{"NAME", "PLATFORMS", "TAGS", "DESCRIPTION", "build:java", "docker:up", "Build Java", "Docker up", "docker"} {
		if !strings.Contains(out, want) {
			t.Errorf("RenderScriptTable() missing %q in:\n%s", want, out)
		}
	}

	if strings.Index(out, "build:java") > strings.Index(out, "docker:up") {
		t.Errorf("scripts not sorted alphabetically:\n%s", out)
	}
}

func TestFormatScriptDetail(t *testing.T) {
	script := config.ScriptConfig{
		Description: "Build Java",
		Command:     "./mvnw package",
		Tags:        []string{"java", "build"},
		Platforms:   map[string]string{"linux": "make java"},
		Env:         map[string]string{"DEBUG": "1", "APP": "backend"},
		Dir:         "./backend",
	}

	out := FormatScriptDetail("build:java", script, "linux", "make java")

	for _, want := range []string{"build:java", "Build Java", "java, build", "make java", "./mvnw package", "On this OS (linux)", "./backend", "APP=backend", "DEBUG=1"} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatScriptDetail() missing %q in:\n%s", want, out)
		}
	}
}

func TestFormatScriptDetailNoCommand(t *testing.T) {
	script := config.ScriptConfig{
		Platforms: map[string]string{"windows": "rmdir dist"},
	}

	out := FormatScriptDetail("clean", script, "linux", "")

	if !strings.Contains(out, "no command available") {
		t.Errorf("FormatScriptDetail() missing warning in:\n%s", out)
	}
}
