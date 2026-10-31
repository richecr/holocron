package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/ui"
)

var removeCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a script from the catalog",
	Example: `  hc remove docker:up`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, path := loadCatalog()

		if _, ok := cfg.Scripts[name]; !ok {
			fatal("%s", notFoundError("script", name, ui.ScriptNames(cfg.Scripts)))
		}

		delete(cfg.Scripts, name)

		if err := saveCatalog(path, cfg); err != nil {
			fatal("failed to save config: %v", err)
		}

		fmt.Printf("%s %s\n",
			ui.WarningStyle.Render("Removed"),
			ui.ScriptNameStyle.Render(name),
		)
	},
}
