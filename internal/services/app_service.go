package services

import (
	"PushMe/constant"
	"PushMe/internal/utils"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type AppService struct{}

var pageType = struct {
	Message string
	Plugin  string
	Note    string
}{
	Message: "message",
	Plugin:  "plugin",
	Note:    "note",
}
var pageStore map[string]*application.WebviewWindow = map[string]*application.WebviewWindow{}

func closePage(pageType string, event *application.CustomEvent) {
	var key string
	switch v := event.Data.(type) {
	case int:
		key = strconv.Itoa(v)
	case int64:
		key = strconv.FormatInt(v, 10)
	case float64: // JSON 数字默认是 float64
		key = fmt.Sprint(v)
	case string:
		key = v
	default:
		return
	}
	win, ok := pageStore[pageType+key]
	if ok {
		win.Close()
	}
}

func init() {
	time.AfterFunc(time.Second, func() {
		utils.EventOn("message:close", func(event *application.CustomEvent) {
			closePage(pageType.Message, event)
		})
	})
	time.AfterFunc(time.Second, func() {
		utils.EventOn("plugin:close", func(event *application.CustomEvent) {
			closePage(pageType.Plugin, event)
		})
	})
	time.AfterFunc(time.Second, func() {
		utils.EventOn("note:close", func(event *application.CustomEvent) {
			closePage(pageType.Note, event)
		})
	})
}

func (a *AppService) OpenPage(url string, title string) *application.WebviewWindow {
	app := application.Get()
	if title == "" {
		title = constant.AppName
	}
	return app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: title,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              url,
		Width:            constant.WindowWidth,
		Height:           constant.WindowHeight,
	})
}

func (a *AppService) GoOpenMessage(id int) {
	var page *application.WebviewWindow = a.OpenPage("/index.html?page=Message&id="+strconv.Itoa(id), "消息内容")
	var pageKey = pageType.Message + strconv.Itoa(id)
	page.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		delete(pageStore, pageKey)
	})
	pageStore[pageKey] = page
}

func (a *AppService) GoOpenDashboard(id int) {
	a.OpenPage("/index.html?page=Dashboard", "数据小屏")
}

func (a *AppService) GoOpenPlugin() {
	a.OpenPage("/index.html?page=Plugin", "消息插件")
}

func (a *AppService) GoOpenPluginEdit(id int) {
	var page *application.WebviewWindow = a.OpenPage("/index.html?page=PluginEdit&id="+strconv.Itoa(id), "插件编辑")
	var pageKey = pageType.Plugin + strconv.Itoa(id)
	page.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		log.Println("events Close Message")
		delete(pageStore, pageKey)
	})
	pageStore[pageKey] = page
}

func (a *AppService) GoOpenSetting(id int) {
	a.OpenPage("/index.html?page=Setting", "系统设置")
}

func (a *AppService) GoOpenUser(user string) {
	var pageKey = pageType.Message + "user_" + user
	// 如果已打开同用户页面，聚焦到已有窗口
	if win, ok := pageStore[pageKey]; ok {
		win.SetAlwaysOnTop(true)
		win.SetAlwaysOnTop(false)
		return
	}
	var page = a.OpenPage("/index.html?page=User&user="+user, user)
	pageStore[pageKey] = page
	page.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		delete(pageStore, pageKey)
	})
}

func (a *AppService) GoOpenNoteEdit(id int) {
	var pageKey string
	var url string
	var title string
	if id > 0 {
		url = "/index.html?page=NoteEdit&id=" + strconv.Itoa(id)
		pageKey = pageType.Note + strconv.Itoa(id)
		title = "编辑便签"
	} else {
		url = "/index.html?page=NoteEdit"
		pageKey = pageType.Note + "0"
		title = "新建便签"
	}
	var page *application.WebviewWindow = a.OpenPage(url, title)
	page.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		delete(pageStore, pageKey)
	})
	pageStore[pageKey] = page
}
