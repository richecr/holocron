package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var (
	listTag  string
	listJSON bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scripts in the catalog",
	Example: `  hc list
  hc list --tag docker
  hc list --json | jq '.[].name'`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := loadCatalog()

		scripts := cfg.Scripts
		if listTag != "" {
			scripts = filterByTag(scripts, listTag)
		}

		if len(scripts) == 0 {
			if listJSON {
				fmt.Println("[]")
				return
			}
			fmt.Println(ui.WarningStyle.Render("No scripts found."))
			return
		}

		if listJSON {
			if err := printJSON(scripts); err != nil {
				fatal("failed to encode JSON: %v", err)
			}
			return
		}

		fmt.Println(ui.RenderScriptTable(scripts))
	},
}

type jsonScript struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Command     string            `json:"command,omitempty"`
	Platforms   map[string]string `json:"platforms,omitempty"`
}

func printJSON(scripts map[string]config.ScriptConfig) error {
	out := make([]jsonScript, 0, len(scripts))
	for _, name := range ui.ScriptNames(scripts) {
		s := scripts[name]
		out = append(out, jsonScript{
			Name:        name,
			Description: s.Description,
			Tags:        s.Tags,
			Command:     s.Command,
			Platforms:   s.Platforms,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
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

func init() {
	listCmd.Flags().StringVar(&listTag, "tag", "", "filter scripts by tag")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output the catalog as JSON")

	if err := listCmd.RegisterFlagCompletionFunc("tag",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, _ := loadCatalog()
			seen := map[string]bool{}
			var tags []string
			for _, script := range cfg.Scripts {
				for _, t := range script.Tags {
					if !seen[t] {
						seen[t] = true
						tags = append(tags, t)
					}
				}
			}
			return tags, cobra.ShellCompDirectiveNoFileComp
		}); err != nil {
		panic(err)
	}
}
