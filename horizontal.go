package acidtab

import (
	"io"

	"zgo.at/termtext"
)

type writer interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteByte(byte) error
	WriteRune(rune) (int, error)
}

type wrapWriter struct {
	w io.Writer
}

func (w wrapWriter) Write(b []byte) (int, error)       { return w.w.Write(b) }
func (w wrapWriter) WriteString(s string) (int, error) { return w.Write([]byte(s)) }
func (w wrapWriter) WriteByte(b byte) error            { _, err := w.Write([]byte{b}); return err }
func (w wrapWriter) WriteRune(r rune) (int, error)     { return w.Write([]byte(string(r))) }

func getWriter(w io.Writer) writer {
	if ww, ok := w.(writer); ok {
		return ww
	}
	return wrapWriter{w}
}

func fillRunes(r rune, n int) string {
	d := make([]rune, n)
	for i := range d {
		d[i] = r
	}
	return string(d)
}
func fillBytes(r byte, n int) string {
	d := make([]byte, n)
	for i := range d {
		d[i] = r
	}
	return string(d)
}

func (t Table) Horizontal(w io.Writer) {
	b := getWriter(w)
	padStr := fillRunes(t.borders.Line, termtext.Width(t.pad))

	if t.close&CloseTop != 0 {
		t.horiLine(b, padStr,
			t.borders.LineTop, t.borders.TopLeft, t.borders.TopRight)
	}
	if t.pHeader {
		t.horiRow(b, t.header, true)
		t.horiLine(b, padStr,
			t.borders.Cross, t.borders.BarRight, t.borders.BarLeft)
	}

	for _, r := range t.rows {
		if len(r) < len(t.header) {
			m := make([]string, len(t.header))
			copy(m, r)
			r = m
		}
		t.horiRow(b, r, false)
	}

	if t.close&CloseBottom != 0 {
		t.horiLine(b, padStr,
			t.borders.LineBottom, t.borders.BottomLeft, t.borders.BottomRight)
	}
}

func (t Table) horiRow(b writer, row []string, alwaysCenter bool) {
	b.WriteString(t.prefix)
	if t.close&CloseLeft != 0 {
		b.WriteRune(t.borders.Bar)
	}
	for i := range row {
		/// In case the header was set to something larger later on.
		if i > len(t.header)-1 {
			continue
		}

		b.WriteString(t.pad)
		align := fillBytes(' ', t.widths[i]-termtext.Width(row[i]))
		a := t.align[i]
		if alwaysCenter {
			a = Center
		}
		switch a {
		case Auto:
			// TODO: this is set in a different location, and not correct after
			// increasing the header size.
			fallthrough
		case Left:
			b.WriteString(row[i])
			if t.close&CloseRight != 0 || i != len(row)-1 {
				b.WriteString(align)
			}
		case Right:
			b.WriteString(align)
			b.WriteString(row[i])
		case Center:
			l := len(align)
			align = align[:l/2]
			b.WriteString(align)
			b.WriteString(row[i])
			if t.close&CloseRight != 0 || i != len(row)-1 {
				if l%2 == 1 {
					b.WriteByte(' ')
				}
				b.WriteString(align)
			}
		}
		if t.close&CloseRight != 0 || i != len(row)-1 {
			b.WriteString(t.pad)
			b.WriteRune(t.borders.Bar)
		}
	}
	b.WriteByte('\n')
}

func (t Table) horiLine(b writer, padStr string, cross, first, last rune) {
	b.WriteString(t.prefix)
	if t.close&CloseLeft != 0 {
		b.WriteRune(first)
	}
	for i := range t.header {
		b.WriteString(padStr)
		b.WriteString(fillRunes(t.borders.Line, t.widths[i]))
		b.WriteString(padStr)
		if i < len(t.header)-1 {
			b.WriteRune(cross)
		} else if t.close&CloseRight != 0 {
			b.WriteRune(last)
		}
	}
	b.WriteByte('\n')
}
