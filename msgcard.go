package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MsgCard struct {
	widget.BaseWidget
	Msg *Msg
}

func NewMsgCard(msg *Msg) *MsgCard {
	var card = &MsgCard{Msg: msg}
	card.ExtendBaseWidget(card)
	return card
}

type MsgCardRenderer struct {
	card    *MsgCard
	objs    []fyne.CanvasObject
	title   *canvas.Text
	content *MsgLabel
	date    *canvas.Text
}

func (c *MsgCard) CreateRenderer() fyne.WidgetRenderer {
	var r = &MsgCardRenderer{card: c}
	return r
}

func (r *MsgCardRenderer) Destroy() {
}

func (r *MsgCardRenderer) Layout(size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	objs := r.Objects()
	for i, obj := range objs {
		if i == 1 {
			obj.Resize(size)
		}
		obj.Move(pos)
		pos = pos.AddXY(0, obj.MinSize().Height)
	}
}

func (r *MsgCardRenderer) MinSize() fyne.Size {
	var size fyne.Size
	objs := r.Objects()
	for _, obj := range objs {
		size = size.Add(obj.MinSize())
	}
	return size
}

func (r *MsgCardRenderer) Refresh() {
	r.title.Text = r.card.Msg.Title
	r.content.SetText(r.card.Msg.Content)
	r.date.Text = r.card.Msg.Date

	r.title.Refresh()
	r.content.Refresh()
	r.date.Refresh()
}

func (r *MsgCardRenderer) Objects() []fyne.CanvasObject {
	if len(r.objs) == 0 {
		var objs []fyne.CanvasObject
		r.title = canvas.NewText(r.card.Msg.Title, theme.ForegroundColor())
		r.title.TextSize = 16
		r.content = NewMsgLabelWithLines(r.card.Msg.Content, 2)
		r.content.Resize(fyne.NewSize(1, 1))
		r.date = canvas.NewText(r.card.Msg.Date, theme.ForegroundColor())
		r.objs = append(objs, container.NewPadded(r.title), container.NewPadded(r.content), container.NewPadded(r.date))
	}
	return r.objs
}
