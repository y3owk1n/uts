// Package banner contains functions for rendering styled banners.
//
//nolint:mnd
package banner

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/y3owk1n/uts/internal/ui/style"
)

// Logo renders the application logo with version.
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
		Render("One CLI for every format")

	return mark + wordmark + versionStyle + "\n" + tagline + "\n"
}

// Header renders a section header with underline.
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

// Mark renders the application wordmark.
func Mark(palette style.Palette) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Primary).
		Render("uts")
}

func stringsRepeat(str string, count int) string {
	if count <= 0 {
		return ""
	}

	out := make([]byte, 0, len(str)*count)
	for range count {
		out = append(out, str...)
	}

	return string(out)
}
