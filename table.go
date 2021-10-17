// Package acidtab prints aligned tables.
package acidtab

import (
	"fmt"
	"strings"

	"zgo.at/termtext"
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
	BordersDouble  = Borders{'═', '║', '╬', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩'}
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

// Table defines a table to print.
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

	err error
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

// Setters.

func (t *Table) Close(close Close) *Table                 { t.close = close; return t }
func (t *Table) Pad(pad string) *Table                    { t.pad = pad; return t }
func (t *Table) Prefix(prefix string) *Table              { t.prefix = prefix; return t }
func (t *Table) Borders(borders Borders) *Table           { t.borders = borders; return t }
func (t *Table) Header(on bool) *Table                    { t.pHeader = on; return t }
func (t *Table) AlignCol(n int, a Align) *Table           { t.align[n] = a; return t }
func (t *Table) PrintCol(n int, p PrintAs) *Table         { t.printAs[n] = p; return t }
func (t *Table) PrintFuncCol(n int, p PrintAsFunc) *Table { t.printAsF[n] = p; return t }

// Error returns any error that may have happened when setting the data.
//
// The print functions will never set an error.
func (t Table) Error() error {
	return t.err
}

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
// the number of headers; it will set an error if it's not.
//
// For example:
//
//   t.Rows(
//       "row1", "row1",
//       "row2", "row2",)
func (t *Table) Rows(r ...interface{}) *Table {
	l := len(t.header)
	if len(r)%l != 0 {
		t.err = fmt.Errorf(
			"Rows: number of cells (%d) not a multitide of number of headers (%d)",
			len(r), l)
		return t
	}
	t.Grow(len(r) / l)
	for ; len(r) > 0; r = r[l:] {
		t.Row(r[:l]...)
	}
	return t
}

// StringRows adds multiple rows from a single string. Columns are separated by
// colDelim, and rows by rowDelim.
//
// Leading and trailing whitespace will be removed, as will all whitespace
// surrounding the colDelim. All items will be added as strings, but you can
// still parse/format things with PrintFuncCol().
//
// For example:
//
//   t.Rows("|", "\n", `
//       row1 | row1
//       row2 | row2
//   `)
//
// The biggest advantage is that it looks a bit nicer if you want to print out
// loads of static data, since gofmt doesn't format it too well. It's also a bit
// easier to write.
func (t *Table) StringRows(colDelim, rowDelim, rows string) *Table {
	r := strings.Split(strings.TrimSpace(rows), rowDelim)
	t.Grow(len(r) * len(t.header))
	for _, rr := range r {
		t.stringRow(colDelim, rr)
	}
	return t
}

// Row adds a new row.
//
// Remaining columns will be filled with spaces if the number of values is lower
// than the numbers of headers. It will set an error if the number of values is
// greater.
func (t *Table) Row(r ...interface{}) *Table {
	if len(r) > len(t.header) {
		t.err = fmt.Errorf(
			"Row: adding row %d: too many values (%d); there are only %d headers",
			len(t.rows), len(r), len(t.header))
		return t
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

func (t *Table) stringRow(colDelim, row string) {
	cols := strings.Split(row, colDelim)
	r := make([]interface{}, 0, len(cols))
	for _, c := range cols {
		r = append(r, strings.TrimSpace(c))
	}
	t.Row(r...)
}
