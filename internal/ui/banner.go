package ui

import "github.com/charmbracelet/lipgloss"

func Banner() string {
	line := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray)).
		Render("◆──────────────────────────────◆")

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorJedi)).
		Render(" H O L O C R O N ")

	sub := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGold)).
		Render("the Force-powered script manager")

	return line + "\n  " + title + " " + sub + "\n" + line
}
