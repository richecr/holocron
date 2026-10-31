package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the catalog for configuration problems",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := loadCatalog()

		issues := config.Validate(cfg)
		if len(issues) == 0 {
			fmt.Println(ui.SuccessStyle.Render("OK") + " " +
				ui.MutedStyle.Render(fmt.Sprintf("%d scripts, no issues found", len(cfg.Scripts))))
			return
		}

		errors := 0
		for _, issue := range issues {
			style := ui.WarningStyle
			if issue.Severity == "error" {
				style = ui.ErrorStyle
				errors++
			}
			fmt.Printf("%s %s\n", style.Render(issue.Severity), issue.Message)
		}

		if errors > 0 {
			fmt.Println(ui.WarningStyle.Render(fmt.Sprintf("%d error(s), %d warning(s)", errors, len(issues)-errors)))
			os.Exit(1)
		}
	},
}
