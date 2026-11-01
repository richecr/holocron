package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/richecr/holocron/internal/config"
)

func ResolveCommand(script config.ScriptConfig) (string, error) {
	currentOS := runtime.GOOS

	if cmd, ok := script.Platforms[currentOS]; ok && cmd != "" {
		return cmd, nil
	}

	if script.Command != "" {
		return script.Command, nil
	}

	return "", fmt.Errorf("no command available for platform: %s", currentOS)
}

// ResolveRunChain returns the ordered list of script names to execute for
// name: pre_run scripts first (depth-first, each script at most once), the
// requested script last. It rejects unknown references and cycles.
func ResolveRunChain(scripts map[string]config.ScriptConfig, name string) ([]string, error) {
	if _, ok := scripts[name]; !ok {
		return nil, fmt.Errorf("script not found: %s", name)
	}

	var order []string
	visited := map[string]bool{}

	var visit func(node, parent string, stack []string) error
	visit = func(node, parent string, stack []string) error {
		for i, s := range stack {
			if s == node {
				cycle := append(stack[i:], node)
				return fmt.Errorf("pre_run cycle: %s", strings.Join(cycle, " -> "))
			}
		}
		if visited[node] {
			return nil
		}

		script, ok := scripts[node]
		if !ok {
			return fmt.Errorf("pre_run of %q references unknown script %q", parent, node)
		}

		for _, pre := range script.PreRun {
			if err := visit(pre, node, append(stack, node)); err != nil {
				return err
			}
		}

		visited[node] = true
		order = append(order, node)
		return nil
	}

	if err := visit(name, "", nil); err != nil {
		return nil, err
	}
	return order, nil
}

func selectShell(settings config.Settings, goos, override string) (string, string, error) {
	shell := override
	if shell == "" {
		shell = settings.DefaultShell[goos]
	}

	if shell != "" {
		switch shell {
		case "cmd":
			return "cmd.exe", "/c", nil
		case "powershell", "powershell.exe":
			return "powershell.exe", "-NoProfile", nil
		case "sh", "bash", "zsh":
			return shell, "-c", nil
		default:
			if override != "" {
				return "", "", fmt.Errorf("invalid shell %q (accepted: sh, bash, zsh, powershell, cmd)", override)
			}
			// Invalid shell in config falls through to the OS default;
			// hc doctor reports it.
		}
	}

	if goos == "windows" {
		return "powershell.exe", "-NoProfile", nil
	}
	return "/bin/sh", "-c", nil
}

func buildInvocation(shell, flag, commandStr string) []string {
	if shell == "powershell.exe" {
		return []string{flag, "-Command", commandStr}
	}
	return []string{flag, commandStr}
}

func BuildCommand(script config.ScriptConfig, extraArgs []string) (string, error) {
	commandStr, err := ResolveCommand(script)
	if err != nil {
		return "", err
	}

	if len(extraArgs) > 0 {
		commandStr = commandStr + " " + strings.Join(extraArgs, " ")
	}

	return commandStr, nil
}

func Run(script config.ScriptConfig, settings config.Settings, extraArgs []string, shellOverride string) error {
	commandStr, err := BuildCommand(script, extraArgs)
	if err != nil {
		return err
	}

	shell, flag, err := selectShell(settings, runtime.GOOS, shellOverride)
	if err != nil {
		return err
	}
	args := buildInvocation(shell, flag, commandStr)

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if script.Dir != "" {
		cmd.Dir = script.Dir
	}

	if len(script.Env) > 0 {
		env := os.Environ()
		for k, v := range script.Env {
			env = append(env, k+"="+v)
		}
		cmd.Env = env
	}

	return cmd.Run()
}
