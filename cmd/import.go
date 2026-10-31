package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import scripts from a TOML file into the catalog",
	Long: `Import scripts from a TOML file (or stdin with "-") into the catalog.

Imported scripts replace scripts with the same name; scripts that exist
only in the current catalog are kept. Settings are never modified — the
file's [settings] section is ignored. The imported scripts are validated
before anything is saved; errors abort the import.

Typical round trip:

  hc export > holocron.toml        # on machine A
  hc import holocron.toml          # on machine B`,
	Example: `  hc import holocron.toml
  curl -s https://example.com/team.toml | hc import -`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var err error

		if args[0] == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(args[0])
		}
		if err != nil {
			fatal("failed to read import source: %v", err)
		}

		incoming, err := config.Decode(data)
		if err != nil {
			fatal("%v", err)
		}

		if len(incoming.Scripts) == 0 {
			fatal("no scripts found in import source")
		}

		issues := config.Validate(&config.Config{Scripts: incoming.Scripts})
		for _, issue := range issues {
			if issue.Severity == "error" {
				fatal("invalid import: %s", issue.Message)
			}
			fmt.Println(ui.WarningStyle.Render(issue.Severity), issue.Message)
		}

		cfg, path := loadCatalog()
		cfg.Scripts = mergeImportedScripts(cfg.Scripts, incoming.Scripts)

		if err := saveCatalog(path, cfg); err != nil {
			fatal("failed to save config: %v", err)
		}

		fmt.Printf("%s %d script(s) from %s\n",
			ui.SuccessStyle.Render("Imported"),
			len(incoming.Scripts),
			ui.MutedStyle.Render(args[0]),
		)
	},
}

func mergeImportedScripts(existing, incoming map[string]config.ScriptConfig) map[string]config.ScriptConfig {
	merged := make(map[string]config.ScriptConfig, len(existing)+len(incoming))
	for name, script := range existing {
		merged[name] = script
	}
	for name, script := range incoming {
		merged[name] = script
	}
	return merged
}
