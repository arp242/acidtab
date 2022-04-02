package acidtab

import (
	"io"

	"zgo.at/termtext"
)

// Vertical prints the table as vertical.
//
// Data is always left-aligned, and Header(false) has no effect.
func (t Table) Vertical(w io.Writer) {
	b := getWriter(w)

	// We calculate this data when rows are added for horizontal tables; need to
	// do different width calculations for vertical tables.
	var (
		padWidth    = termtext.Width(t.pad)
		alignHeader = make([]string, len(t.header))
		headerWidth int
		valueWidth  int
	)
	for i := range t.header {
		if l := termtext.Width(t.header[i]); l > headerWidth {
			headerWidth = l
		}
	}
	for i := range t.header {
		alignHeader[i] = fillBytes(' ', headerWidth-termtext.Width(t.header[i]))
	}
	for _, w := range t.widths {
		if w > valueWidth {
			valueWidth = w
		}
	}
	var (
		padStr    = fillRunes(t.borders.Line, padWidth)
		valueStr  = fillRunes(t.borders.Line, valueWidth)
		headerStr = fillRunes(t.borders.Line, headerWidth)
	)

	/// Write the actual table.
	if t.close&CloseTop != 0 {
		t.vertLine(b, padStr, headerStr, valueStr,
			t.borders.LineTop, t.borders.TopLeft, t.borders.TopRight)
	}

	for i := range t.rows {
		if i > 0 {
			t.vertLine(b, padStr, headerStr, valueStr,
				t.borders.Cross, t.borders.BarRight, t.borders.BarLeft)
		}
		for j := range t.header {
			/// Write header.
			b.WriteString(t.prefix)
			if t.close&CloseLeft != 0 {
				b.WriteRune(t.borders.Bar)
			}
			b.WriteString(t.pad)
			b.WriteString(t.header[j])
			b.WriteString(t.pad)
			b.WriteString(alignHeader[j])
			b.WriteRune(t.borders.Bar)

			/// Write data.
			str := ""
			if len(t.rows[i])-1 >= j { /// In case the header size changed.
				str = t.rows[i][j]
			}

			b.WriteString(t.pad)
			b.WriteString(str)
			if t.close&CloseRight != 0 {
				b.WriteString(fillBytes(' ', valueWidth-termtext.Width(str)))
				b.WriteString(t.pad)
				b.WriteRune(t.borders.Bar)
			}
			b.WriteByte('\n')
		}
	}

	if t.close&CloseBottom != 0 {
		t.vertLine(b, padStr, headerStr, valueStr,
			t.borders.LineBottom, t.borders.BottomLeft, t.borders.BottomRight)
	}
}

func (t Table) vertLine(b writer, padStr, headerStr, valueStr string, cross, first, last rune) {
	b.WriteString(t.prefix)
	if t.close&CloseLeft != 0 {
		b.WriteRune(first)
	}
	b.WriteString(padStr)
	b.WriteString(headerStr)
	b.WriteString(padStr)
	b.WriteRune(cross)
	b.WriteString(padStr)
	b.WriteString(valueStr)
	b.WriteString(padStr)
	if t.close&CloseRight != 0 {
		b.WriteRune(last)
	}
	b.WriteByte('\n')
}
