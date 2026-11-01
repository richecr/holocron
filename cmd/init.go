package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var initYes bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter catalog (example scripts)",
	Long: `Create a starter catalog in the holocron config file.

If the config file does not exist yet it is created; if it exists but has
no scripts, example scripts are added. Scripts already present are never
touched. The examples are harmless echo commands you can edit or remove.

The examples demonstrate the three patterns: a generic command, a
per-platform command and pre_run chaining.`,
	Example: `  hc init`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, path := loadCatalog()

		if len(cfg.Scripts) > 0 {
			fatal("catalog already has %d script(s): %s (use `hc list` or `hc config edit`)",
				len(cfg.Scripts), path)
		}

		if !initYes && !confirm("Add example scripts to the empty catalog?") {
			fmt.Println(ui.WarningStyle.Render("Aborted."))
			return
		}

		for name, script := range exampleScripts() {
			cfg.Scripts[name] = script
		}

		if err := saveCatalog(path, cfg); err != nil {
			fatal("failed to save config: %v", err)
		}

		fmt.Printf("%s %d example script(s) in %s\n",
			ui.SuccessStyle.Render("Initialized"),
			len(cfg.Scripts),
			ui.FormatPath(path),
		)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  %-28s %s\n", "hc list", "see the catalog")
		fmt.Printf("  %-28s %s\n", "hc hello", "run the first example")
		fmt.Printf("  %-28s %s\n", `hc add my:script "echo hi"`, "add your own")
	},
}

func exampleScripts() map[string]config.ScriptConfig {
	return map[string]config.ScriptConfig{
		"hello": {
			Description: "Sanity check: prints a greeting",
			Command:     "echo \"Hello from Holocron!\"",
			Tags:        []string{"example"},
		},
		"clean": {
			Description: "Remove build output (per-platform)",
			Tags:        []string{"example"},
			Platforms: map[string]string{
				"linux":   "rm -rf dist/",
				"darwin":  "rm -rf dist/",
				"windows": "rmdir /s /q dist",
			},
		},
		"deploy": {
			Description: "Example of pre_run: runs 'clean' first",
			Command:     "echo \"deploying... (try 'hc --dry deploy' to see the chain)\"",
			Tags:        []string{"example"},
			PreRun:      []string{"clean"},
		},
	}
}

func confirm(prompt string) bool {
	fmt.Println(ui.WarningStyle.Render(prompt + " [y/N]"))
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

func init() {
	initCmd.Flags().BoolVar(&initYes, "yes", false, "do not ask for confirmation")
}
