package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/richecr/holocron/internal/config"
	"github.com/richecr/holocron/internal/ui"
)

var configCmd = &cobra.Command{
	Use:   "config [path|edit]",
	Short: "Show or edit the holocron configuration",
	Run: func(cmd *cobra.Command, args []string) {
		printConfigPath()
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the absolute path of config.toml",
	Run: func(cmd *cobra.Command, args []string) {
		printConfigPath()
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open config.toml in your default editor",
	Run: func(cmd *cobra.Command, args []string) {
		_, path := loadCatalog()

		editor := resolveEditor()
		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fatal("failed to open editor %q: %v", editor, err)
		}
	},
}

func init() {
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configEditCmd)
}

func printConfigPath() {
	path, err := config.GetConfigPath()
	if err != nil {
		fatal("failed to resolve config path: %v", err)
	}
	fmt.Println(ui.FormatPath(path))
}

func resolveEditor() string {
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "nano"
}
