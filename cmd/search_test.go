package cmd

import (
	"testing"

	"github.com/richecr/holocron/internal/config"
)

func TestSearchScripts(t *testing.T) {
	scripts := map[string]config.ScriptConfig{
		"build:java": {
			Description: "Build the Java project",
			Tags:        []string{"java", "build"},
			Platforms:   map[string]string{"linux": "./mvnw clean package"},
		},
		"docker:up": {
			Description: "Start core services",
			Command:     "docker compose up -d",
			Tags:        []string{"docker"},
		},
		"deploy": {
			Description: "Ship it",
			Command:     "echo deploying",
			Dir:         "./deploy",
		},
	}

	tests := []struct {
		term string
		want []string
	}{
		{"mvnw", []string{"build:java"}},
		{"docker", []string{"docker:up"}},
		{"compose", []string{"docker:up"}},
		{"JAVA", []string{"build:java"}},
		{"./deploy", []string{"deploy"}},
		{"build", []string{"build:java"}}, // matches name and description
		{"ship", []string{"deploy"}},      // matches description
		{"nothing-matches", nil},
	}

	for _, tt := range tests {
		t.Run(tt.term, func(t *testing.T) {
			matches := searchScripts(scripts, tt.term)
			if len(matches) != len(tt.want) {
				t.Fatalf("searchScripts(%q) = %v, want %v", tt.term, keys(matches), tt.want)
			}
			for _, w := range tt.want {
				if _, ok := matches[w]; !ok {
					t.Errorf("searchScripts(%q) missing %q", tt.term, w)
				}
			}
		})
	}
}

func keys(m map[string]config.ScriptConfig) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
