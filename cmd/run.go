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

var runCmd = &cobra.Command{
	Use:   "run <script> [args...]",
	Short: "Run a script from the catalog",
	Long: `Run a script from the catalog.

Extra arguments after the script name (optionally separated by --) are
appended to the resolved command, e.g.:

  hc run build:java -- --debug`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]
		extraArgs := args[1:]

		path, err := config.GetConfigPath()
		if err != nil {
			fatal("failed to resolve config path: %v", err)
		}

		cfg, err := config.LoadConfig(path)
		if err != nil {
			fatal("failed to load config: %v", err)
		}

		script, ok := cfg.Scripts[scriptName]
		if !ok {
			fatal("script not found: %s", scriptName)
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
	},
}

func init() {

	runCmd.Flags().SetInterspersed(false)
}
