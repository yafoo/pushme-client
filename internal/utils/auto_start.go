package utils

import (
	"PushMe/constant"
	"fmt"
	"os"
	"path"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

var linkPath string

func init() {
	linkPath = path.Join(constant.UserDir, constant.LinkSuffix)
}

// 开机启动
func MakeShortcut() (bool, error) {
	if !IsWindows() {
		return false, nil
	}
	err := createShortcut(constant.AppPath, linkPath)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

// 去掉开机启动
func RemoveShortcut() bool {
	if !IsWindows() {
		return false
	}
	err := os.Remove(linkPath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// 卸载应用（删除所有快捷方式）
func Uninstall() (bool, error) {
	if !IsWindows() {
		return false, nil
	}

	// 删除启动目录快捷方式
	err := os.Remove(linkPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("删除启动快捷方式失败:", err)
		return false, err
	}

	// 删除开始菜单快捷方式
	startMenuPath := path.Join(constant.UserDir, constant.StartMenuLnk)
	err = os.Remove(startMenuPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("删除开始菜单快捷方式失败:", err)
		return false, err
	}

	return true, nil
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
