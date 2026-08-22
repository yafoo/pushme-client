package services

import (
	"PushMe/constant"
	db "PushMe/internal/models"
	"PushMe/internal/utils"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (u *UtilsService) GoProxyImage(imageUrl string) string {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", imageUrl, nil)
	if err != nil {
		log.Printf("创建图片请求失败: %v", err)
		return ""
	}

	// 设置常见的 User-Agent 和 Referer 绕过防盗链
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", imageUrl)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("下载图片失败: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("图片下载状态码: %d", resp.StatusCode)
		return ""
	}

	// 限制图片大小 10MB
	if resp.ContentLength > 10*1024*1024 {
		log.Printf("图片过大: %d bytes", resp.ContentLength)
		return ""
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取图片数据失败: %v", err)
		return ""
	}

	// 检测图片类型
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		contentType = "image/jpeg"
	}

	// 转换为 base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
}
