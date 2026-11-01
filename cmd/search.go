package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var searchJSON bool

var searchCmd = &cobra.Command{
	Use:   "search <term>",
	Short: "Search scripts by name, description, tag or command",
	Long: `Search the catalog for scripts matching a term (case-insensitive).

A script matches when the term appears in its name, description, any tag,
its generic command, any per-platform command or its dir.`,
	Example: `  hc search mvnw
  hc search docker
  hc search --json build`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := strings.ToLower(strings.Join(args, " "))

		cfg, _ := loadCatalog()
		matches := searchScripts(cfg.Scripts, term)

		if len(matches) == 0 {
			fmt.Println(ui.WarningStyle.Render("No scripts match ") + ui.ScriptNameStyle.Render(term))
			return
		}

		if searchJSON {
			if err := printJSON(matches); err != nil {
				fatal("failed to encode JSON: %v", err)
			}
			return
		}

		fmt.Println(ui.RenderScriptTable(matches))
	},
}

func searchScripts(scripts map[string]config.ScriptConfig, term string) map[string]config.ScriptConfig {
	term = strings.ToLower(term)
	matches := make(map[string]config.ScriptConfig)
	for name, script := range scripts {
		if scriptMatches(script, term) || strings.Contains(strings.ToLower(name), term) {
			matches[name] = script
		}
	}
	return matches
}

func scriptMatches(script config.ScriptConfig, term string) bool {
	if strings.Contains(strings.ToLower(script.Description), term) {
		return true
	}
	if strings.Contains(strings.ToLower(script.Dir), term) {
		return true
	}
	for _, tag := range script.Tags {
		if strings.Contains(strings.ToLower(tag), term) {
			return true
		}
	}
	if strings.Contains(strings.ToLower(script.Command), term) {
		return true
	}
	for _, cmdStr := range script.Platforms {
		if strings.Contains(strings.ToLower(cmdStr), term) {
			return true
		}
	}
	return false
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "output matches as JSON")
}
