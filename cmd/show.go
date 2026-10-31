package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/executor"
	"github.com/richecr/holocron/internal/ui"
)

var showCmd = &cobra.Command{
	Use:     "show <name>",
	Short:   "Show details of a script",
	Example: `  hc show docker:up`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, _ := loadCatalog()

		script, ok := cfg.Scripts[name]
		if !ok {
			fatal("%s", notFoundError("script", name, ui.ScriptNames(cfg.Scripts)))
		}

		resolved, err := executor.ResolveCommand(script)
		if err != nil {
			resolved = ""
		}

		fmt.Println(ui.FormatScriptDetail(name, script, runtime.GOOS, resolved))
	},
}
