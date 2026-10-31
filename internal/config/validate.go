package config

import (
	"fmt"
	"strings"
)

type Issue struct {
	Severity string // "error" or "warning"
	Message  string
}

var validShells = map[string]bool{
	"sh":         true,
	"bash":       true,
	"zsh":        true,
	"powershell": true,
	"cmd":        true,
}

var knownPlatforms = map[string]bool{
	"linux":   true,
	"darwin":  true,
	"windows": true,
}

func Validate(cfg *Config) []Issue {
	var issues []Issue

	for osName, shell := range cfg.Settings.DefaultShell {
		if !validShells[shell] {
			issues = append(issues, Issue{
				Severity: "error",
				Message: fmt.Sprintf("settings.default_shell.%s: invalid shell %q (accepted: sh, bash, zsh, powershell, cmd)",
					osName, shell),
			})
		}
	}

	for name, script := range cfg.Scripts {
		if script.Command == "" && len(script.Platforms) == 0 {
			issues = append(issues, Issue{
				Severity: "error",
				Message:  fmt.Sprintf("scripts.%q: no command and no platforms defined", name),
			})
		}

		if script.Command == "" && len(script.Platforms) > 0 {
			for osName, cmdStr := range script.Platforms {
				if cmdStr == "" {
					issues = append(issues, Issue{
						Severity: "error",
						Message:  fmt.Sprintf("scripts.%q.platforms.%s: empty command", name, osName),
					})
				}
			}
		}

		for osName := range script.Platforms {
			if !knownPlatforms[osName] {
				issues = append(issues, Issue{
					Severity: "warning",
					Message: fmt.Sprintf("scripts.%q.platforms: unknown platform %q (expected linux, darwin or windows)",
						name, osName),
				})
			}
		}
	}

	return issues
}

func FormatIssues(issues []Issue) string {
	var b strings.Builder
	for _, issue := range issues {
		b.WriteString(issue.Severity + ": " + issue.Message + "\n")
	}
	return b.String()
}
