package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/ui"
)

var rootDry bool

var rootCmd = &cobra.Command{
	Use:   "hc",
	Short: "Holocron is a cross-platform script manager",
	Long: ui.Banner() + `

Holocron (hc) keeps your frequently used commands in a single TOML catalog
(~/.config/holocron/config.toml) and runs them on Linux, macOS and Windows
with the right shell for each platform.

Scripts can be run directly: "hc docker:up" is the same as
"hc run docker:up" (built-in commands take precedence over scripts).`,
	Example: `  hc add docker:up "docker compose up -d" --tags docker
  hc docker:up
  hc list
  hc show docker:up
  hc remove docker:up`,
	Args:              cobra.ArbitraryArgs,
	ValidArgsFunction: scriptNameCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		runScript(args[0], stripDashDash(args[1:]), rootDry)
	},
}

func scriptNameCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	cfg, _ := loadCatalog()
	var matches []string
	for _, name := range ui.ScriptNames(cfg.Scripts) {
		if strings.HasPrefix(name, toComplete) {
			matches = append(matches, name)
		}
	}
	return matches, cobra.ShellCompDirectiveNoFileComp
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render(fmt.Sprintf(format, args...)))
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.Flags().BoolVar(&rootDry, "dry", false, "print the resolved command without running it")
	rootCmd.Flags().SetInterspersed(false)
}
