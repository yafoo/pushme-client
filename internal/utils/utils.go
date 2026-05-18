package utils

import (
	"PushMe/constant"
	"log"
	"net"
	"os"
	"os/exec"

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

func Toast(message string) {
	EventEmit("toast", message)
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
