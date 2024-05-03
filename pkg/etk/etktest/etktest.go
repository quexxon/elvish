// Package etktest provides facilities for testing Etk components.
package etktest

import (
	"fmt"
	"strings"

	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/eval/vals"
	"src.elv.sh/pkg/must"
	"src.elv.sh/pkg/ui"
	"src.elv.sh/pkg/wcwidth"
)

type renderOpts struct {
	Width, Height int
}

func (opts *renderOpts) SetDefaultOptions() {
	opts.Width = 40
	opts.Height = 10
}

func MakeFixture(f etk.Comp) func(*eval.Evaler) {
	return func(ev *eval.Evaler) {
		w := etk.AdaptToWidget(f)
		ev.ExtendGlobal(eval.BuildNs().AddGoFns(map[string]any{
			"setup": func(m vals.Map) {
				w = etk.AdaptToWidget(etk.WithStates(f, must.OK1(convertSetStates(m))...))
			},
			"send-keys": func(args ...any) error {
				keys, err := parseKeys(args)
				if err != nil {
					return err
				}
				for _, key := range keys {
					w.Handle(term.KeyEvent(key))
				}
				return nil
			},
			"render": func(fm *eval.Frame, opts renderOpts) error {
				out := fm.ByteOutput()
				buf := w.Render(opts.Width, opts.Height)
				sd, err := bufferToStyleDown(buf, globalStylesheet)
				if err != nil {
					return err
				}
				_, err = out.WriteString(sd)
				return err
			},
		}).Ns())
	}
}

func parseKeys(args []any) ([]ui.Key, error) {
	var keys []ui.Key
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			for _, r := range arg {
				keys = append(keys, ui.Key{Rune: r})
			}
		case vals.List:
			for it := arg.Iterator(); it.HasElem(); it.Next() {
				elem := it.Elem()
				switch elem := elem.(type) {
				case string:
					key, err := ui.ParseKey(elem)
					if err != nil {
						return nil, err
					}
					keys = append(keys, key)
				default:
					return nil, fmt.Errorf("element of list argument must be string, got %s", vals.ReprPlain(elem))
				}
			}
		default:
			return nil, fmt.Errorf("argument must be string or list, got %s", vals.ReprPlain(arg))
		}
	}
	return keys, nil
}

// TODO: This duplicates part of styledown pkg.
var builtinStyleDownChars = map[ui.Style]rune{
	{}:                 ' ',
	{Bold: true}:       '*',
	{Underlined: true}: '_',
	{Inverse: true}:    '#',
	{Fg: ui.Red}:       'R',
	{Fg: ui.Green}:     'G',
}

// TODO: This duplicates much of (*term.Buffer).TTYString.
func bufferToStyleDown(b *term.Buffer, ss stylesheet) (string, error) {
	var sb strings.Builder
	// Top border
	sb.WriteString("┌" + strings.Repeat("─", b.Width) + "┐\n")
	for i, line := range b.Lines {
		// Write the content line.
		sb.WriteRune('│')
		usedWidth := 0
		for _, cell := range line {
			sb.WriteString(cell.Text)
			usedWidth += wcwidth.Of(cell.Text)
		}
		var rightPadding string
		if usedWidth < b.Width {
			rightPadding = strings.Repeat(" ", b.Width-usedWidth)
			sb.WriteString(rightPadding)
		}
		sb.WriteString("│\n")

		// Write the style line.
		// TODO: I shouldn't have to keep track of the column number manually
		sb.WriteRune('│')
		col := 0
		for _, cell := range line {
			style := ui.StyleFromSGR(cell.Style)
			var styleChar rune
			if char, ok := builtinStyleDownChars[style]; ok {
				styleChar = char
			} else if char, ok := ss.charForStyle[style]; ok {
				styleChar = char
			} else {
				return "", fmt.Errorf("no char for style: %v", style)
			}
			styleStr := string(styleChar)
			if i == b.Dot.Line && col == b.Dot.Col {
				styleStr += "\u0305\u0302" // combining overline + combining circumflex
			}
			sb.WriteString(strings.Repeat(styleStr, wcwidth.Of(cell.Text)))
			col += wcwidth.Of(cell.Text)
		}
		if col <= b.Dot.Col {
			sb.WriteString(strings.Repeat(" ", b.Dot.Col-col+1))
			sb.WriteString("\u0305\u0302")
			sb.WriteString(strings.Repeat(" ", b.Width-b.Dot.Col-1))
		} else {
			sb.WriteString(rightPadding)
		}
		sb.WriteString("│\n")
	}
	// Bottom border
	sb.WriteString("└" + strings.Repeat("─", b.Width) + "┘\n")

	return sb.String(), nil
}

var globalStylesheet = newStylesheet(map[rune]string{
	'r': "red",
})

type stylesheet struct {
	stringStyling map[rune]string
	charForStyle  map[ui.Style]rune
}

func newStylesheet(stringStyling map[rune]string) stylesheet {
	charForStyle := make(map[ui.Style]rune)
	for r, s := range stringStyling {
		var st ui.Style
		ui.ApplyStyling(st, ui.ParseStyling(s))
		charForStyle[st] = r
	}
	return stylesheet{stringStyling, charForStyle}
}

// Same as convertSetStates from pkg/etk, copied to avoid the need to export it.
func convertSetStates(m vals.Map) ([]any, error) {
	var setStates []any
	for it := m.Iterator(); it.HasElem(); it.Next() {
		k, v := it.Elem()
		name, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("key should be string")
		}
		setStates = append(setStates, name, v)
	}
	return setStates, nil
}
