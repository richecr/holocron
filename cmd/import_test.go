package cmd

import (
	"testing"

	"github.com/richecr/holocron/internal/config"
)

func TestMergeImportedScripts(t *testing.T) {
	existing := map[string]config.ScriptConfig{
		"keep":     {Command: "keep-cmd"},
		"replace":  {Command: "old-cmd"},
		"conflict": {Command: "existing-cmd"},
	}
	incoming := map[string]config.ScriptConfig{
		"replace":  {Command: "new-cmd"},
		"conflict": {Command: "imported-cmd"},
		"added":    {Command: "added-cmd"},
	}

	merged := mergeImportedScripts(existing, incoming)

	if len(merged) != 4 {
		t.Fatalf("merged size = %d, want 4", len(merged))
	}
	if merged["keep"].Command != "keep-cmd" {
		t.Errorf("keep = %q, want existing value preserved", merged["keep"].Command)
	}
	if merged["replace"].Command != "new-cmd" {
		t.Errorf("replace = %q, want new-cmd", merged["replace"].Command)
	}
	if merged["conflict"].Command != "imported-cmd" {
		t.Errorf("conflict = %q, want imported-cmd", merged["conflict"].Command)
	}
	if merged["added"].Command != "added-cmd" {
		t.Errorf("added = %q, want added-cmd", merged["added"].Command)
	}

	if _, ok := existing["added"]; ok {
		t.Error("existing map was mutated")
	}
}

func TestMergeImportedScriptsEmptyExisting(t *testing.T) {
	merged := mergeImportedScripts(nil, map[string]config.ScriptConfig{"a": {Command: "x"}})

	if len(merged) != 1 || merged["a"].Command != "x" {
		t.Errorf("merged = %v", merged)
	}
}
