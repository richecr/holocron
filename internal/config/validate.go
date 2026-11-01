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

func IsKnownPlatform(name string) bool {
	return knownPlatforms[name]
}

func ValidShells() []string {
	return []string{"sh", "bash", "zsh", "powershell", "cmd"}
}

func IsValidShell(name string) bool {
	return validShells[name]
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

		for _, pre := range script.PreRun {
			if _, ok := cfg.Scripts[pre]; !ok {
				issues = append(issues, Issue{
					Severity: "error",
					Message:  fmt.Sprintf("scripts.%q.pre_run references unknown script %q", name, pre),
				})
			}
		}
	}

	issues = append(issues, detectPreRunCycles(cfg.Scripts)...)

	return issues
}

func detectPreRunCycles(scripts map[string]ScriptConfig) []Issue {
	var issues []Issue
	reported := map[string]bool{}

	var walk func(node string, stack []string) bool
	walk = func(node string, stack []string) bool {
		for i, s := range stack {
			if s == node {
				cycle := canonicalCycle(append(stack[i:], node))
				key := strings.Join(cycle, ">")
				if !reported[key] {
					reported[key] = true
					issues = append(issues, Issue{
						Severity: "error",
						Message:  fmt.Sprintf("scripts: pre_run cycle: %s", strings.Join(cycle, " -> ")),
					})
				}
				return true
			}
		}
		script, ok := scripts[node]
		if !ok {
			return false
		}
		for _, pre := range script.PreRun {
			if walk(pre, append(stack, node)) {
				return true
			}
		}
		return false
	}

	for name := range scripts {
		walk(name, nil)
	}

	return issues
}

// canonicalCycle rotates a cycle so it starts at its smallest element,
// letting the dedupe key treat a -> b -> a and b -> a -> b as the same cycle.
func canonicalCycle(cycle []string) []string {
	core := cycle[:len(cycle)-1]
	minIdx := 0
	for i, s := range core {
		if s < core[minIdx] {
			minIdx = i
		}
	}
	out := make([]string, 0, len(cycle))
	out = append(out, core[minIdx:]...)
	out = append(out, core[:minIdx]...)
	return append(out, core[minIdx])
}

func FormatIssues(issues []Issue) string {
	var b strings.Builder
	for _, issue := range issues {
		b.WriteString(issue.Severity + ": " + issue.Message + "\n")
	}
	return b.String()
}
