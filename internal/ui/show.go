package ui

import (
	"strings"

	"github.com/richecr/holocron/internal/config"
)

func FormatScriptDetail(name string, script config.ScriptConfig, goos, resolved string) string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(name))
	b.WriteString("\n\n")

	b.WriteString(detailRow("Description", script.Description))
	b.WriteString(detailRow("Tags", strings.Join(script.Tags, ", ")))

	b.WriteString("\n")
	b.WriteString(HeaderStyle.Render("Platforms"))
	b.WriteString("\n")

	if script.Command != "" {
		b.WriteString("  " + PlatformBadge("any") + "  " + FormatScriptCommand(script.Command) + "\n")
	}

	for _, osName := range []string{"linux", "darwin", "windows"} {
		if cmd, ok := script.Platforms[osName]; ok && cmd != "" {
			b.WriteString("  " + PlatformBadge(osName) + "  " + FormatScriptCommand(cmd) + "\n")
		}
	}
	for osName, cmd := range script.Platforms {
		if cmd == "" {
			continue
		}
		if osName != "linux" && osName != "darwin" && osName != "windows" {
			b.WriteString("  " + PlatformBadge(osName) + "  " + FormatScriptCommand(cmd) + "\n")
		}
	}

	b.WriteString("\n")

	if resolved == "" {
		b.WriteString(WarningStyle.Render("On this OS (" + goos + "): no command available"))
	} else {
		b.WriteString(detailRow("On this OS ("+goos+")", resolved))
	}

	return strings.TrimRight(b.String(), "\n")
}

func detailRow(label, value string) string {
	const labelWidth = 12
	pad := labelWidth - len(label)
	if pad < 0 {
		pad = 0
	}
	padded := label + strings.Repeat(" ", pad)

	if value == "" {
		return "  " + MutedStyle.Render(padded) + "\n"
	}
	return "  " + MutedStyle.Render(padded) + "  " + value + "\n"
}
