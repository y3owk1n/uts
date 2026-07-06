package message

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/y3owk1n/uts/internal/ui/style"
)

type Icons struct {
	Info    string
	Success string
	Warn    string
	Error   string
	Step    string
	Bullet  string
	Arrow   string
}

func DefaultIcons() Icons {
	return Icons{
		Info:    "ℹ",
		Success: "✓",
		Warn:    "⚠",
		Error:   "✖",
		Step:    "▸",
		Bullet:  "•",
		Arrow:   "→",
	}
}

type Printer struct {
	palette style.Palette
	types   style.Type
	icons   Icons
	out     io.Writer
	errOut  io.Writer
}

func New(palette style.Palette, types style.Type, icons Icons, out, errOut io.Writer) *Printer {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	return &Printer{
		palette: palette,
		types:   types,
		icons:   icons,
		out:     out,
		errOut:  errOut,
	}
}

func Default() *Printer {
	palette := style.Default()
	types := style.Types(palette)
	return New(palette, types, DefaultIcons(), os.Stdout, os.Stderr)
}

func (p *Printer) Infof(format string, args ...any) {
	p.line(p.out, p.icons.Info, fmt.Sprintf(format, args...), p.palette.Accent)
}

func (p *Printer) Successf(format string, args ...any) {
	p.line(p.out, p.icons.Success, fmt.Sprintf(format, args...), p.palette.Success)
}

func (p *Printer) Warnf(format string, args ...any) {
	p.line(p.out, p.icons.Warn, fmt.Sprintf(format, args...), p.palette.Warning)
}

func (p *Printer) Errorf(format string, args ...any) {
	p.line(p.errOut, p.icons.Error, fmt.Sprintf(format, args...), p.palette.Error)
}

func (p *Printer) Stepf(format string, args ...any) {
	p.line(p.out, p.icons.Step, fmt.Sprintf(format, args...), p.palette.Primary)
}

func (p *Printer) Bulletf(format string, args ...any) {
	styled := lipgloss.NewStyle().
		Foreground(p.palette.Muted).
		Render("  " + p.icons.Bullet + " " + fmt.Sprintf(format, args...))
	fmt.Fprintln(p.out, styled)
}

func (p *Printer) Mutedf(format string, args ...any) {
	fmt.Fprintln(p.out, p.types.Muted.Render(fmt.Sprintf(format, args...)))
}

func (p *Printer) Pair(key, value string) {
	styledKey := p.types.Key.Render(key)
	styledVal := p.types.Code.Render(value)
	fmt.Fprintf(p.out, "%s  %s\n", styledKey, styledVal)
}

func (p *Printer) PairLine(key, value string) string {
	return p.types.Key.Render(key) + "  " + p.types.Code.Render(value) + "\n"
}

func (p *Printer) Highlight(text string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(p.palette.Primary).Render(text)
}

func (p *Printer) Dim(text string) string {
	return p.types.Muted.Render(text)
}

func (p *Printer) Success(text string) string {
	return lipgloss.NewStyle().Foreground(p.palette.Success).Render(text)
}

func (p *Printer) Warn(text string) string {
	return lipgloss.NewStyle().Foreground(p.palette.Warning).Render(text)
}

func (p *Printer) Error(text string) string {
	return lipgloss.NewStyle().Foreground(p.palette.Error).Render(text)
}

func (p *Printer) Accent(text string) string {
	return lipgloss.NewStyle().Foreground(p.palette.Accent).Render(text)
}

func (p *Printer) Text(text string) string {
	return lipgloss.NewStyle().Foreground(p.palette.Text).Render(text)
}

func (p *Printer) SuccessRow(text string) string {
	return p.styledIconLine(p.icons.Success, p.palette.Success) + p.types.Body.Render(text) + "\n"
}

func (p *Printer) WarnRow(text string) string {
	return p.styledIconLine(p.icons.Warn, p.palette.Warning) + p.types.Body.Render(text) + "\n"
}

func (p *Printer) ErrorRow(text string) string {
	return p.styledIconLine(p.icons.Error, p.palette.Error) + p.types.Body.Render(text) + "\n"
}

func (p *Printer) Detail(text string) string {
	return p.types.Muted.Render("    "+text) + "\n"
}

func (p *Printer) styledIconLine(icon string, color lipgloss.AdaptiveColor) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(icon) + " "
}

func (p *Printer) line(writer io.Writer, icon string, message string, iconColor lipgloss.AdaptiveColor) {
	styledIcon := lipgloss.NewStyle().Bold(true).Foreground(iconColor).Render(icon)
	styledMsg := p.types.Body.Render(message)
	fmt.Fprintf(writer, "%s %s\n", styledIcon, styledMsg)
}
