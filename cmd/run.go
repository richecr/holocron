package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/executor"
	"github.com/richecr/holocron/internal/ui"
)

var runDry bool
var runShell string

var runCmd = &cobra.Command{
	Use:   "run <script> [args...]",
	Short: "Run a script from the catalog",
	Long: `Run a script from the catalog.

The "run" keyword is optional: "hc build:java" and "hc run build:java"
are equivalent.

Extra arguments after the script name (optionally separated by --) are
appended to the resolved command, e.g.:

  hc run build:java -- --debug`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runScript(args[0], stripDashDash(args[1:]), runDry, runShell)
	},
}

func runScript(scriptName string, extraArgs []string, dry bool, shellOverride string) {
	cfg, _ := loadCatalog()

	if _, ok := cfg.Scripts[scriptName]; !ok {
		fatal("%s", notFoundError("script", scriptName, ui.ScriptNames(cfg.Scripts)))
	}

	chain, err := executor.ResolveRunChain(cfg.Scripts, scriptName)
	if err != nil {
		fatal("%v", err)
	}

	main := cfg.Scripts[scriptName]

	if dry {
		if len(chain) > 1 {
			fmt.Println(ui.MutedStyle.Render("chain: " + strings.Join(chain, " -> ")))
		}
		resolved, err := executor.BuildCommand(main, extraArgs)
		if err != nil {
			fatal("%v", err)
		}
		fmt.Println(ui.FormatScriptCommand(resolved))
		return
	}

	if cfg.Settings.ConfirmBeforeRun {
		resolved, err := executor.ResolveCommand(main)
		if err != nil {
			fatal("%v", err)
		}
		if len(chain) > 1 {
			fmt.Println(ui.MutedStyle.Render("chain: " + strings.Join(chain, " -> ")))
		}
		fmt.Println(ui.PromptConfirm(ui.FormatScriptCommand(resolved)))
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println(ui.WarningStyle.Render("Aborted."))
			return
		}
	}

	for _, name := range chain[:len(chain)-1] {
		fmt.Println(ui.MutedStyle.Render("pre: " + name))
		runOne(cfg.Scripts[name], cfg.Settings, nil, shellOverride)
	}

	runOne(main, cfg.Settings, extraArgs, shellOverride)
}

func runOne(script config.ScriptConfig, settings config.Settings, extraArgs []string, shellOverride string) {
	if err := executor.Run(script, settings, extraArgs, shellOverride); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fatal("%v", err)
	}
}

func stripDashDash(args []string) []string {
	if len(args) > 0 && args[0] == "--" {
		return args[1:]
	}
	return args
}

func init() {
	runCmd.Flags().BoolVar(&runDry, "dry", false, "print the resolved command without running it")
	runCmd.Flags().StringVar(&runShell, "shell", "", "shell for this run (sh, bash, zsh, powershell, cmd); overrides the config")
	runCmd.Flags().SetInterspersed(false)
}
