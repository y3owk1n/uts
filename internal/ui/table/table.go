//nolint:mnd
package table

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/y3owk1n/uts/internal/ui/style"
)

// Table represents a formatted table.
type Table struct {
	headers []string
	rows    [][]string
	current int
	width   int
	wrap    bool
}

// New creates a new Table with the given headers.
func New(headers ...string) *Table {
	return &Table{
		headers: headers,
		width:   80,
	}
}

// Row adds a row to the table.
func (t *Table) Row(values ...string) *Table {
	t.rows = append(t.rows, values)

	return t
}

// Current sets the current (highlighted) row index.
func (t *Table) Current(idx int) *Table {
	t.current = idx

	return t
}

// Width sets the table width.
func (t *Table) Width(width int) *Table {
	t.width = width

	return t
}

// Wrap sets whether to wrap cell content.
func (t *Table) Wrap(wrap bool) *Table {
	t.wrap = wrap

	return t
}

// Render renders the table with the given palette.
func (t *Table) Render(palette style.Palette) string {
	if len(t.headers) == 0 || len(t.rows) == 0 {
		return ""
	}

	colWidths := make([]int, len(t.headers))
	for idx, h := range t.headers {
		colWidths[idx] = len(h)
	}

	for _, row := range t.rows {
		for idx, cell := range row {
			if idx < len(colWidths) && len(cell) > colWidths[idx] {
				colWidths[idx] = len(cell)
			}
		}
	}

	line := strings.Repeat("─", t.width)

	var buf strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(palette.Primary)
	buf.WriteString(headerStyle.Render(strings.Join(t.headers, "  ")))
	buf.WriteString("\n")
	buf.WriteString(lipgloss.NewStyle().Foreground(palette.Border).Render(line))
	buf.WriteString("\n")

	for idx, row := range t.rows {
		cells := make([]string, len(row))
		for jdx, cell := range row {
			if jdx < len(colWidths) {
				cells[jdx] = fmtCell(cell, colWidths[jdx])
			} else {
				cells[jdx] = cell
			}
		}

		rowStr := strings.Join(cells, "  ")

		if idx == t.current {
			buf.WriteString(lipgloss.NewStyle().
				Foreground(palette.Primary).
				Bold(true).
				Render("▸ " + rowStr))
		} else {
			buf.WriteString("  " + rowStr)
		}

		buf.WriteString("\n")
	}

	return buf.String()
}

func fmtCell(str string, width int) string {
	if len(str) >= width {
		return str
	}

	return str + strings.Repeat(" ", width-len(str))
}
