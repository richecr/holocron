package cmd

import "testing"

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"abc", "abc", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"build", "buld", 1},
		{"docker:up", "docker:up", 0},
		{"kitten", "sitting", 3},
	}

	for _, tt := range tests {
		if got := levenshtein(tt.a, tt.b); got != tt.want {
			t.Errorf("levenshtein(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSuggestName(t *testing.T) {
	candidates := []string{"build:java", "docker:up", "clean"}

	tests := []struct {
		input  string
		want   string
		wantOK bool
	}{
		{"buld:java", "build:java", true},
		{"docker:up", "docker:up", true},
		{"totally-different", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got, ok := suggestName(tt.input, candidates)
		if ok != tt.wantOK {
			t.Errorf("suggestName(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
		}
		if ok && got != tt.want {
			t.Errorf("suggestName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNotFoundError(t *testing.T) {
	candidates := []string{"build:java"}

	withSuggestion := notFoundError("script", "buld:java", candidates)
	if withSuggestion != "script not found: buld:java (did you mean build:java?)" {
		t.Errorf("withSuggestion = %q", withSuggestion)
	}

	plain := notFoundError("script", "xyzzy", candidates)
	if plain != "script not found: xyzzy" {
		t.Errorf("plain = %q", plain)
	}
}
