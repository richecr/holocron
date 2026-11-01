package executor

import (
	"strings"
	"testing"

	"github.com/richecr/holocron/internal/config"
)

func TestResolveRunChain(t *testing.T) {
	tests := []struct {
		name    string
		scripts map[string]config.ScriptConfig
		target  string
		want    []string
		wantErr string
	}{
		{
			name: "no pre_run is just the script",
			scripts: map[string]config.ScriptConfig{
				"build": {Command: "make"},
			},
			target: "build",
			want:   []string{"build"},
		},
		{
			name: "linear pre_run runs first",
			scripts: map[string]config.ScriptConfig{
				"deploy": {Command: "ship", PreRun: []string{"test"}},
				"test":   {Command: "go test", PreRun: []string{"build"}},
				"build":  {Command: "make"},
			},
			target: "deploy",
			want:   []string{"build", "test", "deploy"},
		},
		{
			name: "diamond dependency runs shared pre once",
			scripts: map[string]config.ScriptConfig{
				"a": {Command: "a", PreRun: []string{"b", "c"}},
				"b": {Command: "b", PreRun: []string{"d"}},
				"c": {Command: "c", PreRun: []string{"d"}},
				"d": {Command: "d"},
			},
			target: "a",
			want:   []string{"d", "b", "c", "a"},
		},
		{
			name: "unknown target",
			scripts: map[string]config.ScriptConfig{
				"build": {Command: "make"},
			},
			target:  "nope",
			wantErr: "script not found: nope",
		},
		{
			name: "unknown pre_run reference",
			scripts: map[string]config.ScriptConfig{
				"deploy": {Command: "ship", PreRun: []string{"ghost"}},
			},
			target:  "deploy",
			wantErr: `pre_run of "deploy" references unknown script "ghost"`,
		},
		{
			name: "cycle is rejected",
			scripts: map[string]config.ScriptConfig{
				"a": {Command: "a", PreRun: []string{"b"}},
				"b": {Command: "b", PreRun: []string{"a"}},
			},
			target:  "a",
			wantErr: "pre_run cycle: a -> b -> a",
		},
		{
			name: "self cycle is rejected",
			scripts: map[string]config.ScriptConfig{
				"a": {Command: "a", PreRun: []string{"a"}},
			},
			target:  "a",
			wantErr: "pre_run cycle: a -> a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveRunChain(tt.scripts, tt.target)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("ResolveRunChain() = %v, want error %q", got, tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Errorf("ResolveRunChain() error = %q, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ResolveRunChain() error = %v", err)
			}
			if strings.Join(got, ",") != strings.Join(tt.want, ",") {
				t.Errorf("ResolveRunChain() = %v, want %v", got, tt.want)
			}
		})
	}
}
