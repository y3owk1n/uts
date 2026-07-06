// Package panel contains functions for rendering styled panels.
//
//nolint:mnd
package panel

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/y3owk1n/uts/internal/ui/style"
)

// DefaultMaxWidth is the default maximum panel width.
const DefaultMaxWidth = 80

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
	frameWidth := outer - 2

	styled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(palette.Border).
		Padding(0, 1).
		Width(frameWidth)

	return styled.Render(content) + "\n"
}

// Section renders a section with title and content inside a bordered frame.
func Section(palette style.Palette, title, content string) string {
	outer := width()
	frameWidth := outer - 2
	contentWidth := frameWidth - 2

	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(palette.Border).
		Padding(0, 1).
		Width(frameWidth)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Primary).
		Width(contentWidth)

	separator := lipgloss.NewStyle().
		Foreground(palette.Border).
		Render(strings.Repeat("─", contentWidth))

	rendered := titleStyle.Render(title) + "\n" +
		separator + "\n" +
		content

	return frame.Render(rendered) + "\n"
}
