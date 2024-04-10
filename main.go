package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

type Msg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
	Type    string `json:"type"`
}

var (
	msgs = []Msg{}
	list = &widget.List{}
)

func main() {
	initFont()

	app := app.NewWithID("me.i-i.push")
	window := app.NewWindow("PushMe Client")

	go initServer(func(msg Msg) {
		msgs = append([]Msg{msg}, msgs...)
		list.Refresh()
		list.ScrollToTop()

		ntfy := fyne.NewNotification(msg.Title, msg.Content)
		app.SendNotification(ntfy)

		window.RequestFocus()
	})

	window.SetIcon(resourceIconPng)
	window.Resize(fyne.NewSize(300, 500))

	list = widget.NewList(
		func() int {
			return len(msgs)
		},
		func() fyne.CanvasObject {
			return NewMsgCard(&Msg{})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := o.(*MsgCard)
			item.Msg = &msgs[i]
			item.Refresh()
		},
	)
	list.OnSelected = func(i widget.ListItemID) {
		list.Unselect(i)
		del := func() {
			msgs = append(msgs[:i], msgs[i+1:]...)
			list.Refresh()
		}
		clr := func() {
			msgs = []Msg{}
			list.Refresh()
		}
		NewMessage(app, &msgs[i], del, clr)
	}

	window.SetContent(container.NewStack(list))
	window.ShowAndRun()

	tidyUp()
}

func initFont() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func tidyUp() {
	os.Unsetenv("FYNE_FONT")
	fmt.Println("PushMe Client已退出")
}
