package cmd

import (
	"testing"

	"github.com/richecr/holocron/internal/config"
)

func TestMergeScriptNew(t *testing.T) {
	got := mergeScript(config.ScriptConfig{}, "echo hi", "desc", nil, []string{"util"})

	if got.Command != "echo hi" {
		t.Errorf("Command = %q, want %q", got.Command, "echo hi")
	}
	if got.Description != "desc" {
		t.Errorf("Description = %q, want %q", got.Description, "desc")
	}
	if len(got.Tags) != 1 || got.Tags[0] != "util" {
		t.Errorf("Tags = %v, want [util]", got.Tags)
	}
	if got.Platforms != nil {
		t.Errorf("Platforms = %v, want nil", got.Platforms)
	}
}

func TestMergeScriptPlatformsMerges(t *testing.T) {
	s := mergeScript(config.ScriptConfig{}, "rm -rf dist/", "", []string{"linux", "darwin"}, nil)
	s = mergeScript(s, "rmdir /s /q dist", "", []string{"windows"}, nil)

	if len(s.Platforms) != 3 {
		t.Fatalf("Platforms = %v, want 3 entries", s.Platforms)
	}
	if s.Platforms["linux"] != "rm -rf dist/" {
		t.Errorf("Platforms[linux] = %q", s.Platforms["linux"])
	}
	if s.Platforms["darwin"] != "rm -rf dist/" {
		t.Errorf("Platforms[darwin] = %q", s.Platforms["darwin"])
	}
	if s.Platforms["windows"] != "rmdir /s /q dist" {
		t.Errorf("Platforms[windows] = %q", s.Platforms["windows"])
	}
}

func TestMergeScriptKeepsExisting(t *testing.T) {
	existing := config.ScriptConfig{
		Description: "original desc",
		Command:     "old",
		Tags:        []string{"a", "b"},
	}

	got := mergeScript(existing, "new", "", nil, nil)

	if got.Description != "original desc" {
		t.Errorf("Description = %q, want %q", got.Description, "original desc")
	}
	if got.Command != "new" {
		t.Errorf("Command = %q, want %q", got.Command, "new")
	}
	if len(got.Tags) != 2 {
		t.Errorf("Tags = %v, want [a b]", got.Tags)
	}
}

func TestMergeScriptReplacesFlags(t *testing.T) {
	existing := config.ScriptConfig{
		Description: "old desc",
		Command:     "old",
		Tags:        []string{"old-tag"},
		Platforms:   map[string]string{"linux": "old-linux"},
	}

	got := mergeScript(existing, "new-linux", "new desc", []string{"linux"}, []string{"new-tag"})

	if got.Description != "new desc" {
		t.Errorf("Description = %q, want %q", got.Description, "new desc")
	}
	if got.Tags == nil || got.Tags[0] != "new-tag" {
		t.Errorf("Tags = %v, want [new-tag]", got.Tags)
	}
	if got.Platforms["linux"] != "new-linux" {
		t.Errorf("Platforms[linux] = %q, want %q", got.Platforms["linux"], "new-linux")
	}
}

func TestMergeScriptNormalizesPlatform(t *testing.T) {
	got := mergeScript(config.ScriptConfig{}, "cmd", "", []string{" Linux "}, nil)

	if _, ok := got.Platforms["linux"]; !ok {
		t.Errorf("Platforms = %v, want key %q", got.Platforms, "linux")
	}
}
