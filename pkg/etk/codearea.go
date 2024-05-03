package etk

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/parse"
	"src.elv.sh/pkg/ui"
)

func CodeArea(c Context) (View, React) {
	quotePasteVar := State(c, "quote-paste", false)

	pastingVar := State(c, "pasting", false)
	pasteBufferVar := State(c, "paste-buffer", &strings.Builder{})
	innerView, innerReact := codeAreaWithAbbr(c)
	bufferVar := BindState[CodeBuffer](c, "buffer")

	return innerView, c.WithBinding("codearea", func(event term.Event) Action {
		switch event := event.(type) {
		case term.PasteSetting:
			startPaste := bool(event)
			// TODO:
			// resetInserts()
			if startPaste {
				pastingVar.Set(true)
			} else {
				text := pasteBufferVar.Get().String()
				pasteBufferVar.Set(new(strings.Builder))
				pastingVar.Set(false)

				if quotePasteVar.Get() {
					text = parse.Quote(text)
				}
				bufferVar.Swap(insertAtDot(text))
			}
			return Consumed
		case term.KeyEvent:
			key := ui.Key(event)
			if pastingVar.Get() {
				if isFuncKey(key) {
					// TODO: Notify the user of the error, or insert the
					// original character as is.
				} else {
					pasteBufferVar.Get().WriteRune(key.Rune)
				}
			}
		}
		return innerReact(event)
	})
}

func codeAreaWithAbbr(c Context) (View, React) {
	abbrVar := State(c, "abbr", func(func(a, f string)) {})
	cmdAbbrVar := State(c, "cmd-abbr", func(func(a, f string)) {})
	smallWordAbbr := State(c, "small-word-abbr", func(func(a, f string)) {})

	streakVar := State(c, "streak", "")
	innerView, innerReact := codeAreaCore(c)
	bufferVar := BindState[CodeBuffer](c, "buffer")
	return innerView, func(event term.Event) Action {
		if keyEvent, ok := event.(term.KeyEvent); ok {
			bufferBefore := bufferVar.Get()
			action := innerReact(event)
			if action != Consumed {
				return action
			}
			buffer := bufferVar.Get()
			if inserted, ok := isLiteralInsert(keyEvent, bufferBefore, buffer); ok {
				streak := streakVar.Get() + inserted
				if newBuffer, ok := expandSimpleAbbr(abbrVar.Get(), buffer, streak); ok {
					bufferVar.Set(newBuffer)
					streakVar.Set("")
					return Consumed
				}
				if newBuffer, ok := expandCmdAbbr(cmdAbbrVar.Get(), buffer, streak); ok {
					bufferVar.Set(newBuffer)
					streakVar.Set("")
					return Consumed
				}
				if newBuffer, ok := expandSmallWordAbbr(smallWordAbbr.Get(), buffer, streak, keyEvent.Rune, categorizeSmallWord); ok {
					bufferVar.Set(newBuffer)
					streakVar.Set("")
					return Consumed
				}
				streakVar.Set(streak)
			} else {
				streakVar.Set("")
			}
			return Consumed
		} else {
			return innerReact(event)
		}
	}
}

func isLiteralInsert(event term.KeyEvent, before, after CodeBuffer) (string, bool) {
	key := ui.Key(event)
	if isFuncKey(key) {
		return "", false
	} else {
		text := string(key.Rune)
		if after == insertAtDot(text)(before) {
			return text, true
		} else {
			return "", false
		}
	}
}

func codeAreaCore(c Context) (View, React) {
	promptVar := State(c, "prompt", ui.T(""))
	rpromptVar := State(c, "rprompt", ui.T(""))
	bufferVar := State(c, "buffer", CodeBuffer{})
	pendingVar := State(c, "pending", PendingCode{})
	highlighterVar := State(c, "highlighter",
		func(code string) (ui.Text, []ui.Text) { return ui.T(code), nil })

	buffer := bufferVar.Get()
	code, pFrom, pTo := patchPending(buffer, pendingVar.Get())
	styledCode, tips := highlighterVar.Get()(code.Content)
	if pFrom < pTo {
		// Apply stylingForPending to [pFrom, pTo)
		parts := styledCode.Partition(pFrom, pTo)
		pending := ui.StyleText(parts[1], stylingForPending)
		styledCode = ui.Concat(parts[0], pending, parts[2])
	}

	view := &codeAreaView{
		promptVar.Get(), rpromptVar.Get(),
		styledCode, bufferVar.Get().Dot, tips,
	}
	return view, func(event term.Event) Action {
		if event, ok := event.(term.KeyEvent); ok {
			key := ui.Key(event)
			// Implement the absolute essential functionalities here. Others
			// can be added via keybindings.
			switch key {
			case ui.K(ui.Backspace), ui.K('H', ui.Ctrl):
				bufferVar.Swap(backspace)
				return Consumed
			default:
				if !isFuncKey(key) && unicode.IsGraphic(key.Rune) {
					bufferVar.Swap(insertAtDot(string(key.Rune)))
					return Consumed
				}
			}
		}
		return Unused
	}
}

// CodeBuffer represents the buffer of the CodeArea widget.
type CodeBuffer struct {
	// Content of the buffer.
	Content string
	// Position of the dot (more commonly known as the cursor), as a byte index
	// into Content.
	Dot int
}

func insertAtDot(text string) func(CodeBuffer) CodeBuffer {
	return func(buf CodeBuffer) CodeBuffer {
		return CodeBuffer{
			Content: buf.Content[:buf.Dot] + text + buf.Content[buf.Dot:],
			Dot:     buf.Dot + len(text),
		}
	}
}

func backspace(buf CodeBuffer) CodeBuffer {
	_, chop := utf8.DecodeLastRuneInString(buf.Content[:buf.Dot])
	return CodeBuffer{
		Content: buf.Content[:buf.Dot-chop] + buf.Content[buf.Dot:],
		Dot:     buf.Dot - chop,
	}
}

// PendingCode represents pending code, such as during completion.
type PendingCode struct {
	// Beginning index of the text area that the pending code replaces, as a
	// byte index into RawState.Code.
	From int
	// End index of the text area that the pending code replaces, as a byte
	// index into RawState.Code.
	To int
	// The content of the pending code.
	Content string
}

func isFuncKey(key ui.Key) bool {
	return key.Mod != 0 || key.Rune < 0
}
