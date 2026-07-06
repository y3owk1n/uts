// Package panel contains functions for rendering styled panels.
//
//nolint:mnd
package panel

import (
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/y3owk1n/uts/internal/ui/style"
)

// DefaultMaxWidth is the default maximum panel width.
const DefaultMaxWidth = 80

const (
	borderSize = 2 // Left + Right borders consume 2 columns total
	paddingX   = 2 // Horizontal padding is 2 cells on each side
	paddingY   = 1 // Vertical padding is 1 line top/bottom
)

func width() int {
	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || width <= 0 {
		return DefaultMaxWidth
	}

	if width < DefaultMaxWidth {
		return width
	}

	return DefaultMaxWidth
}

// Panel renders a bordered panel with the given content.
func Panel(palette style.Palette, content string) string {
	outer := width()

	// In lipgloss v2, Width sets the total outer width (including borders).
	// Content inside the panel has outer - borderSize (2) - paddingX * 2 (4) width.
	contentWidth := outer - borderSize - (paddingX * 2)

	styled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(palette.Subtle).
		Padding(paddingY, paddingX).
		Width(outer)

	bodyStyle := lipgloss.NewStyle().
		Foreground(palette.Text).
		Width(contentWidth)

	trimmedContent := strings.TrimRight(content, "\r\n")

	return styled.Render(bodyStyle.Render(trimmedContent)) + "\n"
}

// Section renders a section with title and content inside a bordered frame.
func Section(palette style.Palette, title, content string) string {
	outer := width()

	// In lipgloss v2, Width sets the total outer width (including borders).
	// Content inside the section has outer - borderSize (2) - paddingY * 2 (2) width.
	// Note: Section uses paddingY (1 cell) for horizontal padding.
	contentWidth := outer - borderSize - (paddingY * 2)

	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(palette.Subtle).
		Padding(paddingY, paddingY).
		Width(outer)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Primary).
		Width(contentWidth).
		Align(lipgloss.Left)

	bodyStyle := lipgloss.NewStyle().
		Foreground(palette.Text).
		Width(contentWidth)

	separator := lipgloss.NewStyle().
		Foreground(palette.Subtle).
		Render(strings.Repeat("─", contentWidth))

	trimmedContent := strings.TrimRight(content, "\r\n")

	rendered := titleStyle.Render(title) + "\n" +
		separator + "\n" +
		bodyStyle.Render(trimmedContent)

	return frame.Render(rendered) + "\n"
}
