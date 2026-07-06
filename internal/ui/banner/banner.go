package banner

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/y3owk1n/uts/internal/ui/style"
)

func Logo(palette style.Palette, version string) string {
	mark := lipgloss.NewStyle().
		Foreground(palette.Primary).
		Render("▌")

	wordmark := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Primary).
		Padding(0, 0, 0, 1).
		Render("uts")

	versionStyle := lipgloss.NewStyle().
		Foreground(palette.Muted).
		Render(" v" + version)

	tagline := lipgloss.NewStyle().
		Foreground(palette.Muted).
		Padding(0, 0, 0, 2).
		Render("All-in-one utility toolkit")

	return mark + wordmark + versionStyle + "\n" + tagline + "\n"
}

func Header(palette style.Palette, text string) string {
	styled := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Text).
		Render(text)

	sep := lipgloss.NewStyle().
		Foreground(palette.Border).
		Render(stringsRepeat("─", lipgloss.Width(text)))

	return styled + "\n" + sep + "\n"
}

func Mark(palette style.Palette) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Primary).
		Render("uts")
}

func stringsRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	out := make([]byte, 0, len(s)*count)
	for range count {
		out = append(out, s...)
	}
	return string(out)
}
