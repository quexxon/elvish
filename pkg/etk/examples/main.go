// Example terminal apps implemented using the Etk framework.
package main

import (
	"flag"
	"fmt"
	"os"

	"src.elv.sh/pkg/cli"
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/must"
	"src.elv.sh/pkg/ui"
)

var example = flag.String("example", "", "the example to run")

func main() {
	flag.Parse()

	etk.Notify = func(ui.Text) {}
	var f etk.Comp
	switch *example {
	// 7GUIs
	case "counter":
		f = Counter
	case "temperature":
		f = Temperature
	case "flight":
		f = Flight

	case "codearea":
		f = etk.WithStates(etk.CodeArea,
			"prompt", ui.T("~> "),
			"abbr", func(y func(a, f string)) { y("foo", "lorem") })
	case "wizard":
		f = Wizard
	case "todo":
		f = Todo
	case "preso":
		content := must.OK1(os.ReadFile(flag.Args()[0]))
		pages := parsePreso(string(content))
		f = etk.WithStates(Preso, "pages", pages)
	case "hiernav":
		f = etk.WithStates(HierNav, "data", hierData)
	case "life":
		f = etk.WithStates(Life,
			"name", "pentadecathlon",
			"history", []Board{pentadecathlon})
	default:
		fmt.Println("unknown example:", *example)
		return
	}
	etk.Run(cli.NewTTY(os.Stdin, os.Stdout),
		etk.WithStates(Wrapper, "inner-comp", f), binding)
}

func binding(ev term.Event, c etk.Context, tag string, f etk.React) etk.Action {
	if tag == "codearea" {
		action := f(ev)
		if action == etk.Unused {
			bufferVar := etk.BindState[etk.CodeBuffer](c, "buffer")
			switch ev {
			case term.K(ui.Left):
				bufferVar.Swap(makeMove(moveDotLeft))
			case term.K(ui.Right):
				bufferVar.Swap(makeMove(moveDotRight))
			case term.K(ui.Home):
				bufferVar.Swap(makeMove(moveDotSOL))
			case term.K(ui.End):
				bufferVar.Swap(makeMove(moveDotEOL))
			default:
				return etk.Unused
			}
			return etk.Consumed
		}
		return action
	}
	return f(ev)
}
