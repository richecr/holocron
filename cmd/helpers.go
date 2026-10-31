package cmd

import (
	"github.com/richecr/holocron/internal/config"
)

func loadCatalog() (*config.Config, string) {
	path, err := config.GetConfigPath()
	if err != nil {
		fatal("failed to resolve config path: %v", err)
	}

	cfg, err := config.LoadConfig(path)
	if err != nil {
		fatal("failed to load config: %v", err)
	}

	return cfg, path
}

func saveCatalog(path string, cfg *config.Config) error {
	return config.SaveConfig(path, cfg)
}

func suggestName(input string, candidates []string) (string, bool) {
	best := ""
	bestDist := 3

	for _, c := range candidates {
		d := levenshtein(input, c)
		if d < bestDist {
			bestDist = d
			best = c
		}
	}

	if best == "" {
		return "", false
	}
	return best, true
}

func notFoundError(kind, name string, candidates []string) string {
	if s, ok := suggestName(name, candidates); ok {
		return kind + " not found: " + name + " (did you mean " + s + "?)"
	}
	return kind + " not found: " + name
}

func levenshtein(a, b string) int {
	if a == b {
		return 0
	}

	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)

	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(prev[j]+1, min(curr[j-1]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}

	return prev[lb]
}
