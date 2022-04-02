package acidtab

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func trim(s string) string {
	lines := strings.Split(strings.Trim(s, "\t\n"), "\n")
	for i := range lines {
		lines[i] = strings.TrimLeft(lines[i], "\t")
	}
	return strings.Join(lines, "\n") + "\n"
}

func test(t *testing.T, f func(io.Writer), want string) {
	t.Helper()
	want = trim(want)
	have := new(bytes.Buffer)
	f(have)
	if have.String() != want {
		t.Errorf("\nwant:\n%[1]s\nhave:\n%[2]s\nwant: %[1]q\nhave: %[2]q", want, have.String())
	}
}

func errorContains(have error, want string) bool {
	if have == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(have.Error(), want)
}

func TestHeader(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		tbl := New("one", "two").Close(CloseLeft|CloseRight).Rows("aa1", "aa2", "bb1", "bb2").Header(true, "1", "2").Row("cc1", "cc2")
		if err := tbl.Error(); err != nil {
			t.Fatal(err)
		}

		test(t, tbl.Horizontal, `
			│   1   │   2   │
			├───────┼───────┤
			│  aa1  │  aa2  │
			│  bb1  │  bb2  │
			│  cc1  │  cc2  │
		`)
		test(t, tbl.Vertical, `
			│  1  │  aa1  │
			│  2  │  aa2  │
			├─────┼───────┤
			│  1  │  bb1  │
			│  2  │  bb2  │
			├─────┼───────┤
			│  1  │  cc1  │
			│  2  │  cc2  │
		`)
	})

	t.Run("shrink", func(t *testing.T) {
		tbl := New("one", "two").Close(CloseLeft|CloseRight).Rows("aa1", "aa2", "bb1", "bb2").Header(true, "one").Row("cc1")
		if err := tbl.Error(); err != nil {
			t.Fatal(err)
		}

		test(t, tbl.Horizontal, `
			│  one  │
			├───────┤
			│  aa1  │
			│  bb1  │
			│  cc1  │
		`)
		test(t, tbl.Vertical, `
			│  one  │  aa1  │
			├───────┼───────┤
			│  one  │  bb1  │
			├───────┼───────┤
			│  one  │  cc1  │
		`)
	})

	t.Run("grow", func(t *testing.T) {
		tbl := New("one", "two").Close(CloseLeft|CloseRight).
			Rows("aa1", "aa2", "bb1", "bb2").
			Header(true, "one", "two", "three").
			Row("cc1", "cc2", "cc3")
		if err := tbl.Error(); err != nil {
			t.Fatal(err)
		}

		test(t, tbl.Horizontal, `
			│  one  │  two  │  three  │
			├───────┼───────┼─────────┤
			│  aa1  │  aa2  │         │
			│  bb1  │  bb2  │         │
			│  cc1  │  cc2  │  cc3    │
		`)
		test(t, tbl.Vertical, `
			│  one    │  aa1    │
			│  two    │  aa2    │
			│  three  │         │
			├─────────┼─────────┤
			│  one    │  bb1    │
			│  two    │  bb2    │
			│  three  │         │
			├─────────┼─────────┤
			│  one    │  cc1    │
			│  two    │  cc2    │
			│  three  │  cc3    │
		`)
	})
}

func TestErrors(t *testing.T) {
	tests := []struct {
		tbl     *Table
		wantErr string
	}{
		{New("one", "two").Close(CloseLeft | CloseRight).Rows("aa1"), "not a multitude"},
		{New("one", "two").Close(CloseLeft|CloseRight).Row("aa1", "aa2", "aa3"), "too many values"},
		{New("asd").AlignCol(99, Center), "cannot set column 99 as there are only 1 columns"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if !errorContains(tt.tbl.Error(), tt.wantErr) {
				t.Errorf("wrong error\nwant: %s\nhave: %s", tt.wantErr, tt.tbl.Error())
			}
		})
	}
}

func TestWidthAndClose(t *testing.T) {
	bold := func(s string) string { return "\x1b[1m" + s + "\x1b[0m" }
	tbl := New(bold("Name"), bold("Origin"), bold("Job"), bold("Alive")).
		AlignCol(3, Center).
		FormatColFunc(3, func(v interface{}) string {
			if b, ok := v.(bool); ok {
				return map[bool]string{
					true:  "\x1b[32m ✔ \x1b[0m",
					false: "\x1b[31m✘\x1b[0m",
				}[b]
			}
			return "\x00"
		}).
		Rows("James Holden", "Montana 🌎", "Captain 🚀", true)
	if tbl.Error() != nil {
		t.Fatal(tbl.Error())
	}

	if tbl.Width() != 56 {
		t.Error(tbl.Width())
	}

	tbl = tbl.Close(CloseLeft)
	if tbl.Width() != 57 {
		t.Error(tbl.Width())
	}

	tbl = tbl.Close(CloseLeft | CloseRight)
	if tbl.Width() != 58 {
		t.Error(tbl.Width())
	}

	test(t, tbl.Horizontal, ""+
		"│      \x1b[1mName\x1b[0m      │    \x1b[1mOrigin\x1b[0m    │     \x1b[1mJob\x1b[0m      │  \x1b[1mAlive\x1b[0m  │\n"+
		"├────────────────┼──────────────┼──────────────┼─────────┤\n"+
		"│  James Holden  │  Montana 🌎  │  Captain 🚀  │   \x1b[32m ✔ \x1b[0m   │")

	test(t, tbl.Vertical, ""+
		"│  \x1b[1mName\x1b[0m    │  James Holden  │\n"+
		"│  \x1b[1mOrigin\x1b[0m  │  Montana 🌎    │\n"+
		"│  \x1b[1mJob\x1b[0m     │  Captain 🚀    │\n"+
		"│  \x1b[1mAlive\x1b[0m   │  \x1b[32m ✔ \x1b[0m           │\n")
}

func TestGrow(t *testing.T) {
	tbl := New("asd")

	test := func(want string) {
		t.Helper()
		if have := fmt.Sprint(cap(tbl.rows), len(tbl.rows), tbl.rows); have != want {
			t.Errorf("\nhave: %s\nwant: %s", have, want)
		}
	}

	test("0 0 []")

	tbl.Grow(8)
	test("8 0 []")

	tbl.Row("zxc")
	test("8 1 [[zxc]]")

	tbl.Grow(8)
	test("16 1 [[zxc]]")
}

func TestStringRows(t *testing.T) {
	tbl := New("one", "two", "three").Close(CloseLeft|CloseRight).StringRows("\x00", "\n", false,
		"1\x002\x003\n4\x005\x006")
	test(t, tbl.Horizontal, `
		│  one  │  two  │  three  │
		├───────┼───────┼─────────┤
		│  1    │  2    │  3      │
		│  4    │  5    │  6      │
	`)

	// TODO: obscure bug here: the width of the last column is too wide. This is
	// because it calculated the width for "three" before. To reset this we need
	// to scan all the rows; meh.
	tbl = New("one", "two", "three").Close(CloseLeft|CloseRight).StringRows("\x00", "\n", true,
		"1\x002\x003\n4\x005\x006")
	test(t, tbl.Horizontal, `
		│   1   │   2   │    3    │
		├───────┼───────┼─────────┤
		│  4    │  5    │  6      │
	`)
}

func TestAlign(t *testing.T) {
	tbl := New("int", "float", "int64", "uint", "-- complex --", "forceleft").Close(CloseLeft|CloseRight).
		AlignCol(5, Left).
		Rows(1, 1.1, int64(-2), 3, complex(5, 6), 9)

	test(t, tbl.Horizontal, `
		│  int  │  float  │  int64  │  uint  │  -- complex --  │  forceleft  │
		├───────┼─────────┼─────────┼────────┼─────────────────┼─────────────┤
		│    1  │    1.1  │     -2  │     3  │         (5+6i)  │  9          │
	`)

	// Doesn't align on purpose.
	// TODO: this is too wide though
	test(t, tbl.Vertical, `
		│  int            │  1              │
		│  float          │  1.1            │
		│  int64          │  -2             │
		│  uint           │  3              │
		│  -- complex --  │  (5+6i)         │
		│  forceleft      │  9              │
	`)
}

func TestFormatAs(t *testing.T) {
	tbl := New("s").Close(CloseLeft|CloseRight).FormatCol(0, "%q").Row("asd")

	test(t, tbl.Horizontal, `
		│    s    │
		├─────────┤
		│  "asd"  │
	`)

	test(t, tbl.Vertical, `
		│  s  │  "asd"  │
	`)
}

func TestFormatAsFunc(t *testing.T) {
	tbl := New("f1", "f2", "f3", "f4", "f5", "n1", "n2", "n3", "n4").Close(CloseLeft|CloseRight).
		FormatColFunc(0, FormatAsFloat(2)).
		FormatColFunc(1, FormatAsFloat(6)).
		FormatColFunc(2, FormatAsFloat(3)).
		FormatColFunc(3, FormatAsFloat(0)).
		FormatColFunc(4, FormatAsFloat(0)).
		FormatColFunc(5, FormatAsNum).
		FormatColFunc(6, FormatAsNum).
		FormatColFunc(7, FormatAsNum).
		FormatColFunc(8, FormatAsNum).
		FormatColFunc(9, FormatAsNum).
		Row(1.5, 1.5, 0.8, 1.4, 1.6, 1234, uint64(123456789), 12341.123131, int16(-9999))

	test(t, tbl.Horizontal, `
		│   f1   │     f2     │   f3   │  f4  │  f5  │   n1    │      n2       │    n3    │    n4    │
		├────────┼────────────┼────────┼──────┼──────┼─────────┼───────────────┼──────────┼──────────┤
		│  1.50  │  1.500000  │  .800  │   1  │   2  │  1,234  │  123,456,789  │  12,341  │  -9,999  │
	`)

}
