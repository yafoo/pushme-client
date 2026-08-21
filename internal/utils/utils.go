package utils

import (
	"PushMe/constant"
	db "PushMe/internal/models"
	"PushMe/internal/request"
	"PushMe/internal/setting"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func EventEmit(event string, data any) {
	app := application.Get()
	app.Event.Emit(event, data)
}

func EventOn(event string, callback func(event *application.CustomEvent)) {
	app := application.Get()
	app.Event.On(event, callback)
}

func EventOff(event string) {
	app := application.Get()
	app.Event.Off(event)
}

func ContainsAny(title string, keywords []string) bool {
	lowerTitle := strings.ToLower(title)
	for _, keyword := range keywords {
		if strings.Contains(lowerTitle, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func FilterEmpty(keywords []string) []string {
	var result []string
	for _, keyword := range keywords {
		if strings.TrimSpace(keyword) != "" {
			result = append(result, keyword)
		}
	}
	return result
}

func RepostMessage(msg db.Msg) {
	if !setting.Setting.Repost.Enable || setting.Setting.Repost.Url == "" {
		return
	}

	if setting.Setting.Repost.Limit != "" {
		limitKeywords := FilterEmpty(strings.Split(setting.Setting.Repost.Limit, "|"))
		ContainsAnyLimitKeyword := len(limitKeywords) == 0 || ContainsAny(msg.Title, limitKeywords)
		if !ContainsAnyLimitKeyword {
			return
		}
	}

	if setting.Setting.Repost.Omit != "" {
		omitKeywords := FilterEmpty(strings.Split(setting.Setting.Repost.Omit, "|"))
		containsNoOmitKeywords := len(omitKeywords) == 0 || !ContainsAny(msg.Title, omitKeywords)
		if !containsNoOmitKeywords {
			return
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			Toast(r.(string))
		}
	}()

	repostData := map[string]interface{}{
		"title":   msg.Title,
		"content": msg.Content,
		"date":    msg.Date,
		"type":    msg.Type,
	}

	// 如果开启了转发push_key，则带上push_key
	if setting.Setting.Repost.PushKey {
		repostData["push_key"] = setting.Setting.Host.PushKey
	}

	repost := request.Request{
		Url:  setting.Setting.Repost.Url,
		Data: repostData,
	}
	switch setting.Setting.Repost.Method {
	case "GET":
		repost.Get()
	case "POST/JSON":
		repost.Post()
	case "POST/FORM":
		repost.PostForm()
	default:
		repost.Post()
	}
	log.Println("RepostMessage", repostData)
}

func CheckVersion() (result string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			switch v := r.(type) {
			case string:
				result = v
			case error:
				result = v.Error()
			default:
				result = fmt.Sprintf("unknown panic: %v", v)
			}
		}
	}()

	req := request.Request{
		Url: constant.UrlVersion,
		Data: map[string]interface{}{
			"version": constant.AppVersion,
		},
	}
	return req.Post()
}

func Toast(message string) {
	EventEmit("event:toast", message)
}

func GetIps() (ips []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				ips = append(ips, ip)
			}
		}
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ipnet.IP.To4() == nil {
				ip := ipnet.IP.String()
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

// 辅助函数：检查字符串是否在切片中
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		cmd = "open" // 默认尝试 open 命令
		args = []string{url}
	}

	return exec.Command(cmd, args...).Run()
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// wails3暂未实现，本程序仅支持windows
func Autostart(status bool) {
	if status {
		MakeShortcut()
	} else {
		RemoveShortcut()
	}
}

func Restart() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	cmd := exec.Command(constant.AppPath)
	cmd.Start()
	os.Exit(0)
}
