package acidtab_test

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"zgo.at/acidtab"
)

func TestTable(t *testing.T) {
	for _, f := range []func(){
		Example_basic,
		Example_options,
		Example_coloptions,
		Example_vertical,
		Example_chain,
		Example_format,
		Example_stringRows,
	} {
		fmt.Println("=> " +
			strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[2] +
			":\n",
		)
		f()
		fmt.Println()
	}
}

func Example_basic() {
	// Create a new table
	t := acidtab.New("Name", "Origin", "Job", "Speciality", "Alive")

	// Add rows to it
	t.Row("James Holden", "Montana", "Captain", "Tilting windmills", true)
	t.Row("Amos Burton", "Baltimore", "Mechanic", "Specific people skills", true)

	// And then print it:
	t.Horizontal(os.Stdout)

	// Output:
	//       Name      │   Origin    │    Job     │        Speciality        │  Alive
	// ────────────────┼─────────────┼────────────┼──────────────────────────┼─────────
	//   James Holden  │  Montana    │  Captain   │  Tilting windmills       │  true
	//   Amos Burton   │  Baltimore  │  Mechanic  │  Specific people skills  │  true
}

func Example_options() {
	t := acidtab.New("Name", "Origin", "Job", "Speciality", "Alive")

	t.Borders(acidtab.BordersHeavy)                 // Set different borders.
	t.Pad(" ")                                      // Pad cells with one space.
	t.Prefix(" ")                                   // Prefix every line with a space.
	t.Close(acidtab.CloseTop | acidtab.CloseBottom) // "Close" top and bottom.
	t.Header(false)                                 // Don't print the header.

	t.Row("Naomi Nagata", "Pallas", "Mechanic", "Spicy red food", true)
	t.Row("Alex Kamal", "Mars", "Pilot", "Cowboys", false)

	t.Horizontal(os.Stdout)

	// Output:
	//  ━━━━━━━━━━━━━━┳━━━━━━━━┳━━━━━━━━━━┳━━━━━━━━━━━━━━━━┳━━━━━━━
	//   Naomi Nagata ┃ Pallas ┃ Mechanic ┃ Spicy red food ┃ true
	//   Alex Kamal   ┃ Mars   ┃ Pilot    ┃ Cowboys        ┃ false
	//  ━━━━━━━━━━━━━━┻━━━━━━━━┻━━━━━━━━━━┻━━━━━━━━━━━━━━━━┻━━━━━━━
}

func Example_coloptions() {
	t := acidtab.New("Name", "Origin", "Job", "Speciality", "Alive")
	t.Close(acidtab.CloseLeft | acidtab.CloseRight)

	t.AlignCol(3, acidtab.Right) // Align column 3 and 4 (starts at 0)
	t.AlignCol(4, acidtab.Center)

	t.PrintCol(3, "%q") // Print column 3 as %q

	// Callback for column 4
	t.PrintFuncCol(4, func(v interface{}) string {
		if b, ok := v.(bool); ok {
			return map[bool]string{true: "yes", false: "no"}[b]
		}
		// Return a NULL byte to fall back to regular formatting.
		return "\x00"
	})

	t.Row("Joe Miller", "Ceres", "Cop", "Doors 'n corners", false)
	t.Row("Chrisjen Avasarala", "Earth", "Politician", "Insults", true)

	t.Horizontal(os.Stdout)

	// Output:
	// │         Name         │  Origin  │     Job      │      Speciality      │  Alive  │
	// ├──────────────────────┼──────────┼──────────────┼──────────────────────┼─────────┤
	// │  Joe Miller          │  Ceres   │  Cop         │  "Doors 'n corners"  │   no    │
	// │  Chrisjen Avasarala  │  Earth   │  Politician  │           "Insults"  │   yes   │
}

func Example_vertical() {
	t := acidtab.New("Name", "Origin", "Job", "Speciality", "Alive")
	t.Row("Prax Meng", "Ganymede", "Botanist", "Plant metaphors", true)
	t.Row("Klaes Ashford", "The belt", "Pirate", "Singing", "😢")
	t.Vertical(os.Stdout)

	// Output:
	//   Name        │  Prax Meng
	//   Origin      │  Ganymede
	//   Job         │  Botanist
	//   Speciality  │  Plant metaphors
	//   Alive       │  true
	// ──────────────┼───────────────────
	//   Name        │  Klaes Ashford
	//   Origin      │  The belt
	//   Job         │  Pirate
	//   Speciality  │  Singing
	//   Alive       │  😢
}

func Example_chain() {
	acidtab.New("Name", "Origin", "Job", "Speciality", "Alive").
		Close(acidtab.CloseTop|acidtab.CloseBottom).
		Prefix(" ").
		Pad(" ").
		PrintCol(1, "%q").
		Rows(
			"Adolphus Murtry", "Earth", "Security", "General twattery", false,
			"Fred Johnson", "Earth", "Colonol", "Beltalowda", false,
		).
		Vertical(os.Stdout)

	// Output:
	//  ────────────┬──────────────────
	//   Name       │ Adolphus Murtry
	//   Origin     │ "Earth"
	//   Job        │ Security
	//   Speciality │ General twattery
	//   Alive      │ false
	//  ────────────┼──────────────────
	//   Name       │ Fred Johnson
	//   Origin     │ "Earth"
	//   Job        │ Colonol
	//   Speciality │ Beltalowda
	//   Alive      │ false
	//  ────────────┴──────────────────
}

func Example_format() {
	bold := func(s string) string { return "\x1b[1m" + s + "\x1b[0m" }

	t := acidtab.New(bold("Name"), bold("Origin"), bold("Job"), bold("Speciality"), bold("Alive")).
		Close(acidtab.CloseAll).
		AlignCol(4, acidtab.Center).
		PrintFuncCol(4, func(v interface{}) string {
			if b, ok := v.(bool); ok {
				return map[bool]string{
					true:  "\x1b[32m ✔ \x1b[0m",
					false: "\x1b[31m✘\x1b[0m",
				}[b]
			}
			return "\x00"
		})

	t.Rows(
		"James Holden", "Montana 🌎", "Captain 🚀", "Tilting windmills", true,
		"Amos Burton", "Baltimore 🌎", "Mechanic 🔧", "Specific people skills", true,
		"Naomi Nagata", "Pallas 🌌", "Mechanic 💻", "Spicy red food", true,
		"Alex Kamal", "Mars 🔴", "Pilot 🎧", "Cowboys", false,
		"Joe Miller", "Ceres 🌌", "Cop 👮", "Doors 'n corners", true,
		"Chrisjen Avasarala", "Earth 🌏", "Politician 🖕", "Insults", true,
		"Prax Meng", "Ganymede 🌌", "Botanist 🌻", "Plant metaphors", true,
		"Klaes Ashford", "The belt 🌌", "Pirate 🕱", "Singing", "😢",
		"Adolphus Murtry", "Earth 🌎", "Security 💂", "General twattery", false,
		"Fred Johnson", "Earth 🌎", "Colonol 🎖", "Beltalowda", false)

	t.Horizontal(os.Stdout)
	// Output:
	// ┌──────────────────────┬────────────────┬─────────────────┬──────────────────────────┬─────────┐
	// │         [1mName[0m         │     [1mOrigin[0m     │       [1mJob[0m       │        [1mSpeciality[0m        │  [1mAlive[0m  │
	// ├──────────────────────┼────────────────┼─────────────────┼──────────────────────────┼─────────┤
	// │  James Holden        │  Montana 🌎    │  Captain 🚀     │  Tilting windmills       │   [32m ✔ [0m   │
	// │  Amos Burton         │  Baltimore 🌎  │  Mechanic 🔧    │  Specific people skills  │   [32m ✔ [0m   │
	// │  Naomi Nagata        │  Pallas 🌌     │  Mechanic 💻    │  Spicy red food          │   [32m ✔ [0m   │
	// │  Alex Kamal          │  Mars 🔴       │  Pilot 🎧       │  Cowboys                 │    [31m✘[0m    │
	// │  Joe Miller          │  Ceres 🌌      │  Cop 👮         │  Doors 'n corners        │   [32m ✔ [0m   │
	// │  Chrisjen Avasarala  │  Earth 🌏      │  Politician 🖕  │  Insults                 │   [32m ✔ [0m   │
	// │  Prax Meng           │  Ganymede 🌌   │  Botanist 🌻    │  Plant metaphors         │   [32m ✔ [0m   │
	// │  Klaes Ashford       │  The belt 🌌   │  Pirate 🕱       │  Singing                 │   😢    │
	// │  Adolphus Murtry     │  Earth 🌎      │  Security 💂    │  General twattery        │    [31m✘[0m    │
	// │  Fred Johnson        │  Earth 🌎      │  Colonol 🎖      │  Beltalowda              │    [31m✘[0m    │
	// └──────────────────────┴────────────────┴─────────────────┴──────────────────────────┴─────────┘
}

func Example_stringRows() {
	bold := func(s string) string { return "\x1b[1m" + s + "\x1b[0m" }

	t := acidtab.New(bold("Name"), bold("Origin"), bold("Job"), bold("Speciality"), bold("Alive")).
		Close(acidtab.CloseAll).
		AlignCol(4, acidtab.Center).
		PrintFuncCol(4, func(v interface{}) string {
			if b, ok := v.(bool); ok {
				return map[bool]string{
					true:  "\x1b[32m ✔ \x1b[0m",
					false: "\x1b[31m✘\x1b[0m",
				}[b]
			}
			return "\x00"
		})

	t.StringRows("|", "\n", `
		James Holden       | Montana 🌎   | Captain 🚀    | Tilting windmills      | true
		Amos Burton        | Baltimore 🌎 | Mechanic 🔧   | Specific people skills | true
		Naomi Nagata       | Pallas 🌌    | Mechanic 💻   | Spicy red food         | true
		Alex Kamal         | Mars 🔴      | Pilot 🎧      | Cowboys                | false
		Joe Miller         | Ceres 🌌     | Cop 👮        | Doors 'n corners       | true
		Chrisjen Avasarala | Earth 🌏     | Politician 🖕 | Insults                | true
		Prax Meng          | Ganymede 🌌  | Botanist 🌻   | Plant metaphors        | true
		Klaes Ashford      | The belt 🌌  | Pirate 🕱      | Singing                | 😢
		Adolphus Murtry    | Earth 🌎     | Security 💂   | General twattery       | false
		Fred Johnson       | Earth 🌎     | Colonol 🎖    | Beltalowda             | false
	`)

	t.Horizontal(os.Stdout)
	// Output:
	// ┌──────────────────────┬────────────────┬─────────────────┬──────────────────────────┬─────────┐
	// │         [1mName[0m         │     [1mOrigin[0m     │       [1mJob[0m       │        [1mSpeciality[0m        │  [1mAlive[0m  │
	// ├──────────────────────┼────────────────┼─────────────────┼──────────────────────────┼─────────┤
	// │  James Holden        │  Montana 🌎    │  Captain 🚀     │  Tilting windmills       │  true   │
	// │  Amos Burton         │  Baltimore 🌎  │  Mechanic 🔧    │  Specific people skills  │  true   │
	// │  Naomi Nagata        │  Pallas 🌌     │  Mechanic 💻    │  Spicy red food          │  true   │
	// │  Alex Kamal          │  Mars 🔴       │  Pilot 🎧       │  Cowboys                 │  false  │
	// │  Joe Miller          │  Ceres 🌌      │  Cop 👮         │  Doors 'n corners        │  true   │
	// │  Chrisjen Avasarala  │  Earth 🌏      │  Politician 🖕  │  Insults                 │  true   │
	// │  Prax Meng           │  Ganymede 🌌   │  Botanist 🌻    │  Plant metaphors         │  true   │
	// │  Klaes Ashford       │  The belt 🌌   │  Pirate 🕱       │  Singing                 │   😢    │
	// │  Adolphus Murtry     │  Earth 🌎      │  Security 💂    │  General twattery        │  false  │
	// │  Fred Johnson        │  Earth 🌎      │  Colonol 🎖      │  Beltalowda              │  false  │
	// └──────────────────────┴────────────────┴─────────────────┴──────────────────────────┴─────────┘
}
