package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "hc",
	Short: "Holocron is a cross-platform script manager",
	Long: `Holocron (hc) is a cross-platform script manager.

It keeps your frequently used commands in a single TOML catalog
(~/.config/holocron/config.toml) and runs them on Linux, macOS and
Windows with the right shell for each platform.`,
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
	rootCmd.AddCommand(configCmd)
}
