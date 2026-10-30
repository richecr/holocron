package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var (
	addDesc      string
	addPlatforms []string
	addTags      []string
)

var addCmd = &cobra.Command{
	Use:   `add <name> "<command>"`,
	Short: "Add or update a script in the catalog",
	Long: `Add or update a script in the catalog.

Without --platform the command is stored as a generic fallback (command).
With --platform (repeatable) it is stored as a platform-specific command
under platforms["<os>"].

Calling add for an existing script merges: platform commands are added or
updated per OS; description and tags only change when the flags are passed.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		command := args[1]

		path, err := config.GetConfigPath()
		if err != nil {
			fatal("failed to resolve config path: %v", err)
		}

		cfg, err := config.LoadConfig(path)
		if err != nil {
			fatal("failed to load config: %v", err)
		}

		cfg.Scripts[name] = mergeScript(cfg.Scripts[name], command, addDesc, addPlatforms, addTags)

		if err := config.SaveConfig(path, cfg); err != nil {
			fatal("failed to save config: %v", err)
		}

		fmt.Printf("%s %s\n",
			ui.SuccessStyle.Render("Saved"),
			ui.ScriptNameStyle.Render(name),
		)
	},
}

func mergeScript(existing config.ScriptConfig, command, desc string, platforms, tags []string) config.ScriptConfig {
	script := existing

	if desc != "" {
		script.Description = desc
	}
	if len(tags) > 0 {
		script.Tags = tags
	}

	if len(platforms) > 0 {
		if script.Platforms == nil {
			script.Platforms = map[string]string{}
		}
		for _, platform := range platforms {
			script.Platforms[strings.ToLower(strings.TrimSpace(platform))] = command
		}
	} else {
		script.Command = command
	}

	return script
}

func init() {
	addCmd.Flags().StringVar(&addDesc, "desc", "", "script description")
	addCmd.Flags().StringSliceVar(&addPlatforms, "platform", nil, "target platform (linux, darwin, windows); repeatable")
	addCmd.Flags().StringSliceVar(&addTags, "tags", nil, "comma-separated tags (e.g. java,build)")
}
