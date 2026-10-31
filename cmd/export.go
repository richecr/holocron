package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
)

var exportTag string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the catalog as TOML to stdout",
	Long: `Export the catalog (settings and scripts) as TOML to stdout.

Pipe it to a file to back up your catalog or version it in git:

  hc export > holocron.toml

Use --tag to export only the scripts carrying that tag.`,
	Example: `  hc export > holocron.toml
  hc export --tag docker`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := loadCatalog()

		out := *cfg
		if exportTag != "" {
			out.Scripts = filterByTag(cfg.Scripts, exportTag)
		}

		data, err := config.Encode(&out)
		if err != nil {
			fatal("failed to encode config: %v", err)
		}

		if _, err := os.Stdout.Write(data); err != nil {
			fatal("failed to write: %v", err)
		}
	},
}

func init() {
	exportCmd.Flags().StringVar(&exportTag, "tag", "", "export only scripts with this tag")
	if err := exportCmd.RegisterFlagCompletionFunc("tag",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return tagNames(), cobra.ShellCompDirectiveNoFileComp
		}); err != nil {
		panic(err)
	}
}

func tagNames() []string {
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
	return tags
}
