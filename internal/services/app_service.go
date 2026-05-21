package services

import (
	"PushMe/constant"
	"PushMe/internal/utils"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppService struct{}

var windowsMap map[int]*application.WebviewWindow = map[int]*application.WebviewWindow{}

func init() {
	time.AfterFunc(time.Second, func() {
		utils.EventOn("message:close", func(event *application.CustomEvent) {
			var id int
			switch v := event.Data.(type) {
			case int:
				id = v
			case int64:
				id = int(v)
			case float64: // JSON 数字默认是 float64
				id = int(v)
			case string:
				if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
					id = int(parsed)
				} else {
					return
				}
			default:
				return
			}
			win, ok := windowsMap[id]
			if ok {
				win.Close()
			}
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
	windowsMap[id] = a.OpenPage("/index.html?page=Message&id="+strconv.Itoa(id), "消息内容")
}

func (a *AppService) GoOpenDashboard(id int) {
	a.OpenPage("/index.html?page=Dashboard", "数据小屏")
}

func (a *AppService) GoOpenPlugin() {
	a.OpenPage("/index.html?page=Plugin", "消息插件")
}

func (a *AppService) GoOpenPluginEdit(id int) {
	a.OpenPage("/index.html?page=PluginEdit&id="+strconv.Itoa(id), "插件编辑")
}

func (a *AppService) GoOpenSetting(id int) {
	a.OpenPage("/index.html?page=Setting", "系统设置")
}
