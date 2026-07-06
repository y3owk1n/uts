package panel

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/y3owk1n/uts/internal/ui/style"
)

const DefaultMaxWidth = 80

func width() int {
	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || w <= 0 {
		return DefaultMaxWidth
	}
	if w < DefaultMaxWidth {
		return w
	}
	return DefaultMaxWidth
}

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
