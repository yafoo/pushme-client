package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MsgLabel struct {
	widget.BaseWidget
	text    string
	texts   []string
	width   float32
	lines   int
	measure fyne.Size
}

func NewMsgLabel(text string) *MsgLabel {
	var l = &MsgLabel{text: text, lines: 1, width: 0}
	l.ExtendBaseWidget(l)
	return l
}

func NewMsgLabelWithLines(text string, lines int) *MsgLabel {
	label := NewMsgLabel(text)
	label.lines = lines
	return label
}

type MsgLabelRenderer struct {
	label *MsgLabel
	objs  []fyne.CanvasObject
}

func (l *MsgLabel) CreateRenderer() fyne.WidgetRenderer {
	return &MsgLabelRenderer{label: l}
}

func (r *MsgLabelRenderer) Destroy() {
}

func (r *MsgLabelRenderer) Layout(size fyne.Size) {
	r.label.width = size.Width
	if r.label.measure.Height == 0 {
		r.label.measure = fyne.MeasureText("A", theme.TextSize(), fyne.TextStyle{})
	}

	setTextIntro(r.label)

	pos := fyne.NewPos(0, 0)
	for _, obj := range r.Objects() {
		t := obj.(*canvas.Text)
		t.Move(pos)
		pos = pos.AddXY(0, t.MinSize().Height)
	}
}

func (r *MsgLabelRenderer) MinSize() fyne.Size {
	w, h := float32(0), float32(0)
	for i := 0; i < r.label.lines; i++ {
		if r.label.measure.Width > w {
			w = r.label.measure.Width
		}
		h += r.label.measure.Height
	}
	return fyne.NewSize(w, h)
}

func (r *MsgLabelRenderer) Refresh() {
	if r.label.width != 0 {
		setTextIntro(r.label)
	}

	for i, obj := range r.Objects() {
		t, _ := obj.(*canvas.Text)
		if t == nil {
			continue
		}
		if len(r.label.texts) > i {
			t.Text = r.label.texts[i]
		} else {
			t.Text = ""
		}
		t.Refresh()
	}
}

func (r *MsgLabelRenderer) Objects() []fyne.CanvasObject {
	if len(r.objs) == 0 {
		// setTextIntro(r.label)
		for i := 0; i < r.label.lines; i++ {
			t := canvas.NewText("", theme.ForegroundColor())
			if len(r.label.texts) > i {
				t.Text = r.label.texts[i]
			}
			r.objs = append(r.objs, t)
		}
	}
	return r.objs
}

func (l *MsgLabel) SetText(text string) {
	l.text = text
}

func setTextIntro(label *MsgLabel) {
	t := strings.ReplaceAll(label.text, "\n", " ")
	texts := strings.Split(t, "")

	width := float32(0)
	chars := []string{}
	lines := 0
	for i := 0; i < len(texts); i++ {
		measure := fyne.MeasureText(texts[i], theme.TextSize(), fyne.TextStyle{})
		width += measure.Width

		if width > label.width-2*theme.Padding() {
			if lines >= 1 {
				break
			}
			chars = append(chars, "\n")
			width = 0
			lines++
		}
		if texts[i] == " " && width == 0 {
			continue
		}
		chars = append(chars, texts[i])
	}

	label.texts = strings.Split(strings.Join(chars, ""), "\n")
}
