package sys

import (
	"PushMeClient/constant"
	"PushMeClient/db"
	"PushMeClient/utils"
	"PushMeClient/utils/setting"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/go-toast/toast"
)

var linkPath string

func init() {
	linkPath = path.Join(utils.UserDir, constant.LinkSuffix)
}

func Notification(msg db.Msg) {
	if !setting.Setting.Notice.Enable {
		return
	}

	messageRune := []rune(msg.Content)
	var message string
	if len(messageRune) > 255 {
		message = string(messageRune[:255]) + "..."
	} else {
		message = msg.Content
	}
	notice := toast.Notification{
		AppID:   "PushMeClient",
		Title:   msg.Title,
		Message: message,
		Icon:    utils.LogoPath,
		// Actions: []toast.Action{
		// 	{"protocol", "查看", "pushme://message/"+fmt.Sprint(msg.ID)},
		// },
	}

	switch setting.Setting.Notice.Duration {
	case "short":
		notice.Duration = toast.Short
	case "long":
		notice.Duration = toast.Long
	default:
		notice.Duration = toast.Short
	}

	switch setting.Setting.Notice.Audio {
	case "default":
		notice.Audio = toast.Default
	case "im":
		notice.Audio = toast.IM
	case "mail":
		notice.Audio = toast.Mail
	case "reminder":
		notice.Audio = toast.Reminder
	case "sms":
		notice.Audio = toast.SMS
	case "loopingalarm":
		notice.Audio = toast.LoopingAlarm
	case "loopingalarm2":
		notice.Audio = toast.LoopingAlarm2
	case "loopingalarm3":
		notice.Audio = toast.LoopingAlarm3
	case "loopingalarm4":
		notice.Audio = toast.LoopingAlarm4
	case "loopingalarm5":
		notice.Audio = toast.LoopingAlarm5
	case "loopingalarm6":
		notice.Audio = toast.LoopingAlarm6
	case "loopingalarm7":
		notice.Audio = toast.LoopingAlarm7
	case "loopingalarm8":
		notice.Audio = toast.LoopingAlarm8
	case "loopingalarm9":
		notice.Audio = toast.LoopingAlarm9
	case "loopingalarm10":
		notice.Audio = toast.LoopingAlarm10
	case "loopingcall":
		notice.Audio = toast.LoopingCall
	case "loopingcall2":
		notice.Audio = toast.LoopingCall2
	case "loopingcall3":
		notice.Audio = toast.LoopingCall3
	case "loopingcall4":
		notice.Audio = toast.LoopingCall4
	case "loopingcall5":
		notice.Audio = toast.LoopingCall5
	case "loopingcall6":
		notice.Audio = toast.LoopingCall6
	case "loopingcall7":
		notice.Audio = toast.LoopingCall7
	case "loopingcall8":
		notice.Audio = toast.LoopingCall8
	case "loopingcall9":
		notice.Audio = toast.LoopingCall9
	case "loopingcall10":
		notice.Audio = toast.LoopingCall10
	case "silent":
		notice.Audio = toast.Silent
	default:
		notice.Audio = toast.Default
	}

	err := notice.Push()
	if err != nil {
		panic(err)
	}
}

func GetIps() (ips []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

// 开机启动
func MakeShortcut() (bool, error) {
	err := createShortcut(utils.AppPath, linkPath)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

// 去掉开机启动
func RemoveShortcut() bool {
	err := os.Remove(linkPath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func createShortcut(source string, target string) error {
	var err error
	err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer ole.CoUninitialize()
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wShell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wShell.Release()
	cs, err := oleutil.CallMethod(wShell, "CreateShortcut", target)
	if err != nil {
		return err
	}
	iDispatch := cs.ToIDispatch()
	_, err = oleutil.PutProperty(iDispatch, "TargetPath", source)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(iDispatch, "Save")
	if err != nil {
		return err
	}
	return nil
}

func Restart() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	cmd := exec.Command(utils.AppPath)
	cmd.Start()
	os.Exit(0)
}
