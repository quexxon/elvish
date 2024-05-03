package etk

import (
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/ui"
)

var Empty = EmptyView{}

type EmptyView struct{}

func (e EmptyView) Render(width, height int) *term.Buffer { return term.NewBuffer(width) }

type TextView struct {
	Spans     []ui.Text
	DotBefore int
}

func Text(spans ...ui.Text) TextView { return TextView{spans, 0} }

func (t TextView) WithDotBefore(i int) TextView { return TextView{t.Spans, i} }

func (t TextView) Render(width, height int) *term.Buffer {
	bb := term.NewBufferBuilder(width)
	for i, span := range t.Spans {
		bb.WriteStyled(span)
		if i+1 == t.DotBefore {
			bb.SetDotHere()
		}
	}
	buf := bb.Buffer()
	// TODO: dot line
	buf.TrimToLines(0, height)
	return buf
}

type VBoxView struct {
	Rows  []View
	Focus int
}

func VBox(rows ...View) VBoxView { return VBoxView{rows, 0} }

func (v VBoxView) WithFocus(i int) VBoxView { return VBoxView{v.Rows, i} }

func (v VBoxView) Render(width, height int) *term.Buffer {
	if len(v.Rows) == 0 {
		return term.NewBuffer(width)
	}
	buf := v.Rows[0].Render(width, height-len(v.Rows)-1)
	for i := 1; i < len(v.Rows); i++ {
		rowHeight := height - len(buf.Lines) - len(v.Rows) - i - 1
		if rowHeight <= 0 {
			break
		}
		buf.Extend(v.Rows[i].Render(width, rowHeight), i == v.Focus)
	}
	return buf
}

type HBoxView struct {
	Cols  []View
	Focus int
}

func HBox(cols ...View) HBoxView { return HBoxView{cols, 0} }

func (h HBoxView) WithFocus(i int) HBoxView { return HBoxView{h.Cols, i} }

func (h HBoxView) Render(width, height int) *term.Buffer {
	if len(h.Cols) == 0 {
		return term.NewBuffer(width)
	}
	colWidth := width / len(h.Cols)

	buf := h.Cols[0].Render(colWidth, height)
	for i := 1; i < len(h.Cols); i++ {
		// TODO: Focus
		buf.ExtendRight(h.Cols[i].Render(colWidth, height))
	}
	return buf
}

type HBoxFlexView struct {
	Cols  []View
	Focus int
	Gap   int
}

func HBoxFlex(cols ...View) HBoxFlexView { return HBoxFlexView{cols, 0, 0} }

func (h HBoxFlexView) WithFocus(i int) HBoxFlexView { return HBoxFlexView{h.Cols, i, h.Gap} }
func (h HBoxFlexView) WithGap(g int) HBoxFlexView   { return HBoxFlexView{h.Cols, h.Focus, g} }

func (h HBoxFlexView) Render(width, height int) *term.Buffer {
	buf := term.NewBuffer(0)
	if len(h.Cols) == 0 {
		return buf
	}
	// TODO: Handle very narrow width

	for i, col := range h.Cols {
		bufCol := col.Render(width-(h.Gap+1)*(len(h.Cols)-i-1), height)
		actualWidth := term.CellsWidth(bufCol.Lines[0])
		for _, line := range bufCol.Lines[1:] {
			actualWidth = max(actualWidth, term.CellsWidth(line))
		}
		bufCol.Width = actualWidth
		// TODO: Focus
		if i > 0 {
			buf.Width += h.Gap
		}
		buf.ExtendRight(bufCol)
	}
	return buf
}
