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

func selectShell(settings config.Settings, goos string) (string, string) {
	shell := settings.DefaultShell[goos]

	switch shell {
	case "cmd":
		return "cmd.exe", "/c"
	case "powershell", "powershell.exe":
		return "powershell.exe", "-NoProfile"
	case "sh", "bash", "zsh":
		return shell, "-c"
	}

	if goos == "windows" {
		return "powershell.exe", "-NoProfile"
	}
	return "/bin/sh", "-c"
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

func Run(script config.ScriptConfig, settings config.Settings, extraArgs []string) error {
	commandStr, err := BuildCommand(script, extraArgs)
	if err != nil {
		return err
	}

	shell, flag := selectShell(settings, runtime.GOOS)
	args := buildInvocation(shell, flag, commandStr)

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
