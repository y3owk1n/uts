package style

import "github.com/charmbracelet/lipgloss"

// KeyColumnWidth is the default width for key columns.
const KeyColumnWidth = 14

// Type defines a set of text styles for different message types.
type Type struct {
	Title   lipgloss.Style
	Section lipgloss.Style
	Body    lipgloss.Style
	Muted   lipgloss.Style
	Subtle  lipgloss.Style
	Code    lipgloss.Style
	Key     lipgloss.Style
}

// Types returns a Type set styled with the given palette.
func Types(palette Palette) Type {
	return Type{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(palette.Primary).
			MarginBottom(1),
		Section: lipgloss.NewStyle().
			Bold(true).
			Foreground(palette.Text).
			MarginTop(1).
			MarginBottom(1),
		Body: lipgloss.NewStyle().
			Foreground(palette.Text),
		Muted: lipgloss.NewStyle().
			Foreground(palette.Muted),
		Subtle: lipgloss.NewStyle().
			Foreground(palette.Subtle).
			Italic(true),
		Code: lipgloss.NewStyle().
			Foreground(palette.Accent),
		Key: lipgloss.NewStyle().
			Foreground(palette.Muted).
			Width(KeyColumnWidth).
			Align(lipgloss.Right),
	}
}
