// Package acidtab prints aligned tables.
package acidtab

import (
	"fmt"
	"strings"

	"arp242.net/termtext"
)

type (
	Close       uint8  // Which sides of the table to "close".
	Align       uint8  // Alignment for columns.
	PrintAs     string // How to print a value.
	PrintAsFunc func(v interface{}) string

	// Borders to use.
	Borders struct {
		Line, Bar, Cross                           rune
		TopLeft, TopRight, BottomLeft, BottomRight rune
		BarRight, BarLeft, LineTop, LineBottom     rune
	}
)

// Which sides to close.
const (
	CloseBottom Close = 1 << iota
	CloseTop
	CloseLeft
	CloseRight
	CloseAll Close = CloseBottom | CloseTop | CloseLeft | CloseRight
)

// Characters to use to draw the borders.
var (
	BordersDefault = Borders{'─', '│', '┼', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴'}
	BordersHeavy   = Borders{'━', '┃', '╋', '┏', '┓', '┗', '┛', '┣', '┫', '┳', '┻'}
	BordersASCII   = Borders{'-', '|', '+', '+', '+', '+', '+', '+', '+', '+', '+'}
	BordersSpace   = Borders{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
)

// Column alignment.
const (
	Auto Align = iota
	Left
	Right
	Center
)

type Table struct {
	header []string
	rows   [][]string
	widths []int

	close   Close   // Which sides to close?
	borders Borders // Border characters to use.
	pad     string  // Padding between columns, before and after.
	prefix  string  // Print before every line.
	pHeader bool    // Print header?

	printAs  []PrintAs // Printf format verb; defaults to %v
	printAsF []PrintAsFunc
	align    []Align
}

// New creates a new table with the given headers.
func New(header ...string) *Table {
	t := &Table{
		header:   header,
		widths:   make([]int, len(header)),
		printAs:  make([]PrintAs, len(header)),
		printAsF: make([]PrintAsFunc, len(header)),
		align:    make([]Align, len(header)),
		pad:      "  ",
		borders:  BordersDefault,
		pHeader:  true,
	}
	for i := range header {
		t.printAs[i] = "%v"
		t.widths[i] = termtext.Width(header[i])
	}
	return t
}

func (t Table) String() string {
	b := new(strings.Builder)
	t.Horizontal(b)
	return b.String()
}

func (t *Table) Close(close Close) *Table                 { t.close = close; return t }
func (t *Table) Pad(pad string) *Table                    { t.pad = pad; return t }
func (t *Table) Prefix(prefix string) *Table              { t.prefix = prefix; return t }
func (t *Table) Borders(borders Borders) *Table           { t.borders = borders; return t }
func (t *Table) Header(on bool) *Table                    { t.pHeader = on; return t }
func (t *Table) AlignCol(n int, a Align) *Table           { t.align[n] = a; return t }
func (t *Table) PrintCol(n int, p PrintAs) *Table         { t.printAs[n] = p; return t }
func (t *Table) PrintFuncCol(n int, p PrintAsFunc) *Table { t.printAsF[n] = p; return t }

// Width gets the display width of the table, including any padding characters.
//
// The width may grow if more rows are added.
func (t *Table) Width() int {
	p := termtext.Width(t.prefix) + termtext.Width(t.pad)*2 + 1 // 1 for the bar character
	var w int
	for _, c := range t.widths {
		w += c + p
	}
	if t.close&CloseLeft != 0 {
		w++
	}
	if t.close&CloseRight == 0 {
		w--
	}
	return w
}

// Grow the rows allocation.
func (t *Table) Grow(n int) {
	if len(t.rows) == 0 {
		t.rows = make([][]string, 0, n)
		return
	}
	r := make([][]string, len(t.rows), len(t.rows)+n)
	copy(r, t.rows)
	t.rows = r
}

// Rows adds multiple rows; the number of values should be an exact multitude of
// the number of headers.
//
// For example:
//
//   t.Rows(
//       "row1", "row1",
//       "row2", "row2",)
func (t *Table) Rows(r ...interface{}) *Table {
	l := len(t.header)
	if len(r)%l != 0 {
		panic("wrong number")
	}
	t.Grow(len(r) / l)
	for ; len(r) > 0; r = r[l:] {
		t.Row(r[:l]...)
	}
	return t
}

// Row adds a new row.
//
// The number of values can be lower than the number of headers; the remaining
// cells will be filled with spaces.
//
// If the number of values is greater it will panic.
func (t *Table) Row(r ...interface{}) *Table {
	if len(r) > len(t.header) {
		panic(fmt.Sprintf("table.Row: too many values (%d); there are only %d headers", len(r), len(t.header)))
	}

	if len(t.rows) == 0 {
		for i := range r {
			if t.align[i] != Auto {
				continue
			}

			if _, ok := r[i].(int); ok {
				t.align[i] = Right
			} else {
				t.align[i] = Left
			}
		}
	}

	row := make([]string, len(r))
	for i := range r {
		f := t.printAsF[i]
		if f != nil {
			row[i] = f(r[i])
			if row[i] == "\x00" {
				row[i] = fmt.Sprintf(string(t.printAs[i]), r[i])
			}
		} else {
			row[i] = fmt.Sprintf(string(t.printAs[i]), r[i])
		}
		if l := termtext.Width(row[i]); l > t.widths[i] {
			t.widths[i] = l
		}
	}
	t.rows = append(t.rows, row)
	return t
}
