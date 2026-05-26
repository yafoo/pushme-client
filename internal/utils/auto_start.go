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
