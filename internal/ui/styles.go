package ui

import "github.com/charmbracelet/lipgloss"

const (
	ColorJedi  = "#00d7ff"
	ColorSith  = "#ff005f"
	ColorGold  = "#ffaf00"
	ColorGray  = "#585858"
	ColorGreen = "#00ff87"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorJedi))

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorJedi))

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorSith))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGold))

	ScriptNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorGold))

	MutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGray))

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorJedi))
)

func badge(label string, color string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color(color)).
		Padding(0, 1).
		Render(label)
}

func PlatformBadge(platform string) string {
	switch platform {
	case "linux":
		return badge("linux", ColorGreen)
	case "darwin":
		return badge("darwin", ColorJedi)
	case "windows", "win":
		return badge("win", ColorGreen)
	case "any":
		return badge("any", ColorGray)
	default:
		return badge(platform, ColorGray)
	}
}
