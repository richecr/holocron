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
	addPre       []string
)

var addCmd = &cobra.Command{
	Use:   `add <name> "<command>"`,
	Short: "Add or update a script in the catalog",
	Long: `Add or update a script in the catalog.

Without --platform the command is stored as a generic fallback (command).
With --platform (repeatable) it is stored as a platform-specific command
under platforms["<os>"].

Calling add for an existing script merges: platform commands are added or
updated per OS; description, tags and pre_run only change when the flags
are passed.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		command := args[1]

		for _, platform := range addPlatforms {
			normalized := strings.ToLower(strings.TrimSpace(platform))
			if !config.IsKnownPlatform(normalized) {
				fatal("unknown platform %q (accepted: linux, darwin, windows)", platform)
			}
		}

		cfg, path := loadCatalog()

		for _, pre := range addPre {
			if pre == "" {
				continue
			}
			if pre != name {
				if _, ok := cfg.Scripts[pre]; !ok {
					fmt.Println(ui.WarningStyle.Render(fmt.Sprintf(
						"warning: pre_run script %q does not exist yet (hc doctor will flag it)", pre)))
				}
			}
		}

		cfg.Scripts[name] = mergeScript(cfg.Scripts[name], command, addDesc, addPlatforms, addTags, addPre)

		if err := saveCatalog(path, cfg); err != nil {
			fatal("failed to save config: %v", err)
		}

		saved := cfg.Scripts[name]
		fmt.Printf("%s %s %s\n",
			ui.SuccessStyle.Render("Saved"),
			ui.ScriptNameStyle.Render(name),
			ui.PlatformBadges(saved),
		)
	},
}

func mergeScript(existing config.ScriptConfig, command, desc string, platforms, tags, preRun []string) config.ScriptConfig {
	script := existing

	if desc != "" {
		script.Description = desc
	}
	if len(tags) > 0 {
		script.Tags = tags
	}

	if preRun != nil {
		if len(preRun) == 1 && preRun[0] == "" {
			script.PreRun = nil
		} else {
			script.PreRun = preRun
		}
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
	addCmd.Flags().StringSliceVar(&addPre, "pre", nil, "scripts to run before this one, in order (e.g. test,clean); pass \"\" to clear")

	if err := addCmd.RegisterFlagCompletionFunc("platform",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"linux", "darwin", "windows"}, cobra.ShellCompDirectiveNoFileComp
		}); err != nil {
		panic(err)
	}

	if err := addCmd.RegisterFlagCompletionFunc("pre",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, _ := loadCatalog()
			return ui.ScriptNames(cfg.Scripts), cobra.ShellCompDirectiveNoFileComp
		}); err != nil {
		panic(err)
	}
}
