package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Message struct {
	Window fyne.Window
}

func (m *Message) Close() {
	m.Window.Close()
}

func NewMessage(app fyne.App, msg *Msg, del func(), clr func()) *Message {
	window := app.NewWindow("Message")
	window.SetIcon(resourceIconPng)
	window.Resize(fyne.NewSize(300, 500))

	title := canvas.NewText(msg.Title, theme.ForegroundColor())
	title.TextSize = 16
	title.Alignment = fyne.TextAlignCenter

	date := canvas.NewText(msg.Date, theme.ForegroundColor())
	date.TextSize = 12
	date.Alignment = fyne.TextAlignCenter

	var content *widget.RichText
	if msg.Type == "markdown" {
		content = widget.NewRichTextFromMarkdown(msg.Content)
	} else {
		content = widget.NewRichTextWithText(msg.Content)
	}
	content.Wrapping = fyne.TextWrapWord

	message := &Message{Window: window}
	buttons := container.NewBorder(
		nil,
		nil,
		widget.NewButton("清空", func() {
			window.Close()
			clr()
		}),
		widget.NewButton("删除", func() {
			window.Close()
			del()
		}),
	)

	window.SetContent(container.NewBorder(container.NewVBox(container.NewPadded(title), container.NewPadded(date), content), buttons, nil, nil))
	window.Show()
	return message
}
