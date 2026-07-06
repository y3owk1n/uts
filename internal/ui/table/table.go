package table

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/y3owk1n/uts/internal/ui/style"
)

type Table struct {
	headers []string
	rows    [][]string
	current int
	width   int
	wrap    bool
}

func New(headers ...string) *Table {
	return &Table{
		headers: headers,
		width:   80,
	}
}

func (t *Table) Row(values ...string) *Table {
	t.rows = append(t.rows, values)
	return t
}

func (t *Table) Current(idx int) *Table {
	t.current = idx
	return t
}

func (t *Table) Width(w int) *Table {
	t.width = w
	return t
}

func (t *Table) Wrap(b bool) *Table {
	t.wrap = b
	return t
}

func (t *Table) Render(palette style.Palette) string {
	if len(t.headers) == 0 || len(t.rows) == 0 {
		return ""
	}

	colWidths := make([]int, len(t.headers))
	for i, h := range t.headers {
		colWidths[i] = len(h)
	}
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	line := strings.Repeat("─", t.width)
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(palette.Primary)
	b.WriteString(headerStyle.Render(strings.Join(t.headers, "  ")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(palette.Border).Render(line))
	b.WriteString("\n")

	for i, row := range t.rows {
		cells := make([]string, len(row))
		for j, cell := range row {
			if j < len(colWidths) {
				cells[j] = fmtCell(cell, colWidths[j])
			} else {
				cells[j] = cell
			}
		}
		rowStr := strings.Join(cells, "  ")

		if i == t.current {
			b.WriteString(lipgloss.NewStyle().
				Foreground(palette.Primary).
				Bold(true).
				Render("▸ " + rowStr))
		} else {
			b.WriteString("  " + rowStr)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func fmtCell(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
