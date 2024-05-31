package acidtab

import "testing"

var benchTable = func() *Table {
	t := New("Name", "Origin", "Job", "Speciality", "Alive").
		Header(true).
		Prefix("  ").
		AlignCol(4, Center).
		FormatColFunc(4, func(v any) string {
			switch vv := v.(type) {
			default:
				return "\x00"
			case bool:
				if vv {
					return "\x1b[32m ✔ \x1b[0m"
				}
				return "\x1b[31m✘\x1b[0m"
			}
		})
	return t
}

var benchRows = []any{
	"James Holden", "Montana", "Captain", "Tilting windmills", true,
	"Amos Burton", "Baltimore", "Mechanic", "Specific people skills", true,
	"Naomi Nagata", "Pallas", "Mechanic", "Spicy red food", true,
	"Alex Kamal", "Mars", "Pilot", "Cowboys", false,
	"Joe Miller", "Ceres", "Detective", "Doors 'n corners", true,
	"Chrisjen Avasarala", "Earth", "Politician", "Insults", true,
	"Prax Meng", "Ganymede", "Botanist", "Plant metaphors", true,
	"Klaes Ashford", "The belt", "Pirate", "Singing", "😢",
	"Adolphus Murtry", "Earth", "Security", "General twattery", false,
	"Fred Johnson", "Earth", "Colonol", "Beltalowda", false,
}

type blackhole struct{}

func (blackhole) Write([]byte) (int, error) { return 0, nil }

func BenchmarkHorizontal(b *testing.B) {
	t := benchTable()
	buf := new(blackhole)

	b.Run("0", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Horizontal(buf)
		}
	})
	t.Rows(benchRows...)
	b.Run("10", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Horizontal(buf)
		}
	})
	t.Rows(benchRows...)
	b.Run("20", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Horizontal(buf)
		}
	})
	t.Rows(benchRows...)
	t.Rows(benchRows...)
	t.Rows(benchRows...)
	b.Run("50", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Horizontal(buf)
		}
	})

	for len(t.rows) < 1000 {
		t.Rows(benchRows...)
	}
	b.Run("1000", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Horizontal(buf)
		}
	})
}

func BenchmarkVertical(b *testing.B) {
	t := benchTable()
	buf := new(blackhole)

	b.Run("0", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Vertical(buf)
		}
	})
	t.Rows(benchRows...)
	b.Run("10", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Vertical(buf)
		}
	})
	t.Rows(benchRows...)
	b.Run("20", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Vertical(buf)
		}
	})
	t.Rows(benchRows...)
	t.Rows(benchRows...)
	t.Rows(benchRows...)
	b.Run("50", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Vertical(buf)
		}
	})

	for len(t.rows) < 1000 {
		t.Rows(benchRows...)
	}
	b.Run("1000", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			t.Vertical(buf)
		}
	})
}
