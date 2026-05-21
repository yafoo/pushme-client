package utils

import (
	"PushMe/constant"
	db "PushMe/internal/models"
	"PushMe/internal/request"
	"PushMe/internal/setting"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func EventEmit(event string, data ...any) {
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

	repost := request.Request{
		Url: setting.Setting.Repost.Url,
		Data: map[string]interface{}{
			"title":   msg.Title,
			"content": msg.Content,
			"date":    msg.Date,
			"type":    msg.Type,
		},
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
	log.Println("RepostMessage")
}

func CheckVersion() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			Toast(r.(string))
		}
	}()

	req := request.Request{
		Url: constant.UrlVersion,
		Data: map[string]interface{}{
			"version": constant.AppVersion,
		},
	}
	result := req.Post()
	EventEmit("version", result)
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

func Autostart(status bool) {
	// wails3暂未实现
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
