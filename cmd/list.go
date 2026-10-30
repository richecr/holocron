package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var listTag string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scripts in the catalog",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := config.GetConfigPath()
		if err != nil {
			fatal("failed to resolve config path: %v", err)
		}

		cfg, err := config.LoadConfig(path)
		if err != nil {
			fatal("failed to load config: %v", err)
		}

		scripts := cfg.Scripts
		if listTag != "" {
			scripts = filterByTag(scripts, listTag)
		}

		if len(scripts) == 0 {
			fmt.Println(ui.WarningStyle.Render("No scripts found."))
			return
		}

		fmt.Println(ui.RenderScriptTable(scripts))
	},
}

func init() {
	listCmd.Flags().StringVar(&listTag, "tag", "", "filter scripts by tag")
}

func filterByTag(scripts map[string]config.ScriptConfig, tag string) map[string]config.ScriptConfig {
	filtered := make(map[string]config.ScriptConfig)
	for name, script := range scripts {
		if slices.Contains(script.Tags, tag) {
			filtered[name] = script
		}
	}
	return filtered
}
