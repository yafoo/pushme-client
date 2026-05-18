package services

import (
	"PushMe/constant"
	db "PushMe/internal/models"
	"PushMe/internal/utils"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

type AppService struct{}

var windowsMap map[int]*application.WebviewWindow = map[int]*application.WebviewWindow{}

var Notifier = notifications.New()

func init() {
	Notifier.OnNotificationResponse(func(result notifications.NotificationResult) {
		if result.Error != nil {
			log.Println(fmt.Errorf("parsing notification result failed: %s", result.Error))
		} else {
			log.Println("Response: %+v\n", result.Response)
			utils.EventEmit("notification:action", result.Response)
		}
	})

	time.AfterFunc(time.Second, func() {
		utils.EventOn("close-message", func(event *application.CustomEvent) {
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

func (a *AppService) GoNotification(msg db.Msg) {
	authorized, err := Notifier.CheckNotificationAuthorization()
	if err != nil {
		utils.Toast("检查通知权限失败: " + err.Error())
		return
	}
	if authorized {
		err := Notifier.SendNotification(notifications.NotificationOptions{
			ID:    strconv.Itoa(int(msg.ID)),
			Title: msg.Title,
			Body:  msg.Content,
			Data: map[string]interface{}{
				"id": msg.ID,
			},
		})
		if err != nil {
			utils.Toast("发送通知失败: " + err.Error())
			return
		}
		log.Println("发送通知成功: ", msg)
	} else {
		authorized, err = Notifier.RequestNotificationAuthorization()
		if err != nil {
			utils.Toast("请求通知权限失败: " + err.Error())
			return
		}
		if !authorized {
			utils.Toast("用户拒绝了通知权限")
		}
	}
}
