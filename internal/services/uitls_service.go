package services

import (
	"PushMe/constant"
	db "PushMe/internal/models"
	"PushMe/internal/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

type UtilsService struct{}

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
}

func (u *UtilsService) GoNotification(msg db.Msg) {
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
				"id":   msg.ID,
				"type": msg.Type,
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

func (u *UtilsService) GoLog(str string) {
	go log.Println(str)
}

func (u *UtilsService) GoGetIps() []string {
	return utils.GetIps()
}

func (u *UtilsService) GoGetVersion() string {
	return constant.AppVersion
}

func (u *UtilsService) GoCheckVersion() string {
	return utils.CheckVersion()
}

func (u *UtilsService) GoOpenBrowser(url string) bool {
	err := utils.OpenBrowser(url)
	return err == nil
}

func (u *UtilsService) GoRestart() {
	utils.Restart()
}
