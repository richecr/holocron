package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/executor"
	"github.com/richecr/holocron/internal/ui"
)

var runDry bool

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
		runScript(args[0], stripDashDash(args[1:]), runDry)
	},
}

func runScript(scriptName string, extraArgs []string, dry bool) {
	cfg, _ := loadCatalog()

	script, ok := cfg.Scripts[scriptName]
	if !ok {
		fatal("%s", notFoundError("script", scriptName, ui.ScriptNames(cfg.Scripts)))
	}

	if dry {
		resolved, err := executor.BuildCommand(script, extraArgs)
		if err != nil {
			fatal("%v", err)
		}
		fmt.Println(ui.FormatScriptCommand(resolved))
		return
	}

	resolved, err := executor.ResolveCommand(script)
	if err != nil {
		fatal("%v", err)
	}

	if cfg.Settings.ConfirmBeforeRun {
		fmt.Println(ui.PromptConfirm(ui.FormatScriptCommand(resolved)))
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println(ui.WarningStyle.Render("Aborted."))
			return
		}
	}

	if err := executor.Run(script, cfg.Settings, extraArgs); err != nil {
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
	runCmd.Flags().SetInterspersed(false)
}
