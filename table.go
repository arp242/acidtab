// Package acidtab prints aligned tables.
package acidtab

import (
	"fmt"
	"strings"

	"zgo.at/termtext"
)

type (
	Close        uint8  // Which sides of the table to "close".
	Align        uint8  // Alignment for columns.
	FormatAs     string // How to print a value; fmt format string (e.g. "%q", "%#v", etc.)
	FormatAsFunc func(v interface{}) string

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

	printAs  []FormatAs // Printf format verb; defaults to %v
	printAsF []FormatAsFunc
	align    []Align

	err error
}

// New creates a new table with the given headers.
func New(header ...string) *Table {
	t := &Table{pad: "  ", borders: BordersDefault}
	return t.Header(true, header...)
}

func (t Table) String() string {
	b := new(strings.Builder)
	t.Horizontal(b)
	return b.String()
}

// Close sets which sides of the table to close.
func (t *Table) Close(close Close) *Table { t.close = close; return t }

// Pad sets the left/right padding for every cell.
func (t *Table) Pad(pad string) *Table { t.pad = pad; return t }

// Prefix sets the prefix to use before the left border.
func (t *Table) Prefix(prefix string) *Table { t.prefix = prefix; return t }

// Borders sets the characters to use for borders
func (t *Table) Borders(borders Borders) *Table { t.borders = borders; return t }

// AlignCol sets the alignment for column n.
//
// The default is right-aligned for numbers, and left-aligned for everything
// else.
func (t *Table) AlignCol(n int, a Align) *Table {
	if t.checkN(n, "AlignCol") {
		t.align[n] = a
	}
	return t
}

// FormatCol sets how to format column n, as fmt format string (e.g. "%q",
// "%#v", etc.)
func (t *Table) FormatCol(n int, p FormatAs) *Table {
	if t.checkN(n, "FormatCol") {
		t.printAs[n] = p
	}
	return t
}

// FormatColFunc sets a callback function to print a cell.
func (t *Table) FormatColFunc(n int, p FormatAsFunc) *Table {
	if t.checkN(n, "FormatColFunc") {
		t.printAsF[n] = p
	}
	return t
}

func (t *Table) checkN(n int, f string) bool {
	if n > len(t.header)-1 {
		t.err = fmt.Errorf("%s: cannot set column %d as there are only %d columns", f, n, len(t.header))
		return false
	}
	return true
}

// Header sets the header.
//
// The show parameter controls if the header is printed.
//
// If the list of headers is given then it will use this as the headers,
// overriding any previously set headers.
func (t *Table) Header(show bool, header ...string) *Table {
	t.pHeader = show

	if len(header) > 0 {
		switch {
		// Set new header.
		case len(t.header) == 0:
			t.widths = make([]int, len(header))
			t.printAs = make([]FormatAs, len(header))
			t.printAsF = make([]FormatAsFunc, len(header))
			t.align = make([]Align, len(header))
			for i := range header {
				t.printAs[i] = "%v"
			}
		// Grow header.
		case len(header) > len(t.header):
			grow := len(header) - len(t.header)
			t.widths = append(t.widths, make([]int, grow)...)
			t.printAs = append(t.printAs, make([]FormatAs, grow)...)
			t.printAsF = append(t.printAsF, make([]FormatAsFunc, grow)...)
			t.align = append(t.align, make([]Align, grow)...)
			t.printAs[2] = "%v"
			t.widths[2] = 4
			//for i:=grow; i++ {
			//	t.printAs[i] = "%v"
			//}
		// Fill with empty header, so printing is still correct.
		case len(header) < len(t.header):
			//header = append(header, make([]string, len(t.header)-len(header))...)
		}

		t.header = header
		for i := range header {
			if l := termtext.Width(header[i]); l > t.widths[i] {
				t.widths[i] = l
			}
		}
	}

	return t
}

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

// Grow the rows allocation by n.
func (t *Table) Grow(n int) {
	if len(t.rows) == 0 {
		t.rows = make([][]string, 0, n)
		return
	}
	r := make([][]string, len(t.rows), cap(t.rows)+n)
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
			"Rows: number of cells (%d) not a multitude of number of headers (%d)",
			len(r), l)
		return t
	}
	t.Grow(len(r) / l)
	for ; len(r) > 0; r = r[l:] {
		t.Row(r[:l]...)
	}
	return t
}

// RowsFromString adds multiple rows from a single string. Columns are separated
// by colDelim, and rows by rowDelim.
//
// Leading and trailing whitespace will be removed, as will all whitespace
// surrounding the colDelim. All items will be added as strings, but you can
// still parse/format things with FormatColFunc().
//
// For example:
//
//   t.Rows("|", "\n", `
//       row1 | row1
//       row2 | row2
//   `)
//
// If header is set the first row will be used as the header, overriding any
// header that was given with New().
//
// The biggest advantage is that it looks a bit nicer if you want to print out
// loads of static data, since gofmt doesn't format it too well. It's also a bit
// easier to write.
//
// Another use case is to feed output from a program or function.
func (t *Table) RowsFromString(colDelim, rowDelim string, header bool, rows string) *Table {
	r := strings.Split(strings.TrimSpace(rows), rowDelim)
	if header {
		t = t.Header(t.pHeader, strings.Split(r[0], colDelim)...)
		r = r[1:]
	}

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

			if isNumber(r[i]) {
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

func isNumber(i interface{}) bool {
	switch i.(type) {
	default:
		return false
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128:
		return true
	}
}

func (t *Table) stringRow(colDelim, row string) {
	cols := strings.Split(row, colDelim)
	r := make([]interface{}, 0, len(cols))
	for _, c := range cols {
		r = append(r, strings.TrimSpace(c))
	}
	t.Row(r...)
}
