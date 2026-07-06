package table

import (
	"charm.land/lipgloss/v2"
	ltable "charm.land/lipgloss/v2/table"
	"github.com/y3owk1n/uts/internal/ui/style"
)

const (
	noCurrentRow  = -1
	columnPadding = 2
	defaultWrap   = false
)

// Table represents a formatted table.
type Table struct {
	headers    []string
	rows       [][]string
	currentRow int
	width      int
	wrap       bool
}

// New creates a new Table with the given headers.
func New(headers ...string) *Table {
	return &Table{
		headers:    append([]string(nil), headers...),
		currentRow: noCurrentRow,
		wrap:       defaultWrap,
	}
}

// Row adds a row to the table.
func (t *Table) Row(values ...string) *Table {
	copied := append([]string(nil), values...)
	t.rows = append(t.rows, copied)

	return t
}

// Current sets the current (highlighted) row index.
func (t *Table) Current(idx int) *Table {
	t.currentRow = idx

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
	cellPadding := lipgloss.NewStyle().Padding(0, columnPadding)

	headerStyle := cellPadding.
		Bold(true).
		Foreground(palette.Subtle)

	currentStyle := cellPadding.Bold(true)

	bodyStyle := cellPadding

	borderStyle := lipgloss.NewStyle().Foreground(palette.Subtle)

	tbl := ltable.New().
		Headers(t.headers...).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderColumn(false).
		BorderRow(false).
		BorderHeader(true).
		BorderStyle(borderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case ltable.HeaderRow:
				return headerStyle
			case t.currentRow:
				return currentStyle
			default:
				return bodyStyle
			}
		}).
		Wrap(t.wrap)

	if t.width > 0 {
		tbl.Width(t.width)
	}

	for _, row := range t.rows {
		tbl.Row(row...)
	}

	return tbl.String()
}
