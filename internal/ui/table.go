package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/richecr/holocron/internal/config"
)

func renderTable(headers []string, rows [][]string) string {

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = lipgloss.Width(h)
	}
	for _, row := range rows {
		for i, cell := range row {

			if w := plainWidth(cell); w > widths[i] {
				widths[i] = w
			}
		}
	}

	var b strings.Builder

	b.WriteString(renderRow(headers, widths, true))
	b.WriteString("\n")

	for _, row := range rows {
		b.WriteString(renderRow(row, widths, false))
		b.WriteString("\n")
	}

	return strings.TrimRight(b.String(), "\n")
}

func renderRow(cells []string, widths []int, isHeader bool) string {
	parts := make([]string, len(cells))
	for i, cell := range cells {
		pad := widths[i] - plainWidth(cell)
		if pad < 0 {
			pad = 0
		}
		if isHeader {
			parts[i] = HeaderStyle.Render(cell + strings.Repeat(" ", pad))
		} else {
			parts[i] = cell + strings.Repeat(" ", pad)
		}
	}
	line := strings.Join(parts, "  ")

	return strings.TrimRight(line, " ")
}

func plainWidth(s string) int {
	return lipgloss.Width(stripANSI(s))
}

func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		switch {
		case r == '\x1b':
			inEscape = true
		case inEscape && (r == 'm'):
			inEscape = false
		case inEscape:

		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func sortedNames(scripts map[string]config.ScriptConfig) []string {
	return ScriptNames(scripts)
}

func ScriptNames(scripts map[string]config.ScriptConfig) []string {
	names := make([]string, 0, len(scripts))
	for name := range scripts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func PlatformBadges(script config.ScriptConfig) string {
	if len(script.Platforms) == 0 {
		return PlatformBadge("any")
	}

	order := []string{"linux", "darwin", "windows"}
	badges := make([]string, 0, len(script.Platforms))

	for _, osName := range order {
		if cmd, ok := script.Platforms[osName]; ok && cmd != "" {
			badges = append(badges, PlatformBadge(osName))
		}
	}

	for osName, cmd := range script.Platforms {
		if cmd == "" {
			continue
		}
		if osName != "linux" && osName != "darwin" && osName != "windows" {
			badges = append(badges, PlatformBadge(osName))
		}
	}

	if len(badges) == 0 {
		return PlatformBadge("any")
	}

	return strings.Join(badges, " ")
}

func RenderScriptTable(scripts map[string]config.ScriptConfig) string {
	names := sortedNames(scripts)
	rows := make([][]string, 0, len(names))

	for _, name := range names {
		script := scripts[name]
		rows = append(rows, []string{
			ScriptNameStyle.Render(name),
			PlatformBadges(script),
			TagBadges(script.Tags),
			script.Description,
		})
	}

	return renderTable([]string{"NAME", "PLATFORMS", "TAGS", "DESCRIPTION"}, rows)
}

func TagBadges(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	badges := make([]string, 0, len(tags))
	for _, t := range tags {
		badges = append(badges, badge(t, ColorGold))
	}
	return strings.Join(badges, " ")
}

func FormatPath(path string) string {
	return MutedStyle.Render(path)
}

func FormatScriptCommand(cmd string) string {
	return MutedStyle.Render(cmd)
}

func PromptConfirm(command string) string {
	return WarningStyle.Render(fmt.Sprintf("Run %s ? [y/N]", command))
}
