package utils

import (
	"PushMe/constant"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var linkPath string

func init() {
	linkPath = path.Join(constant.UserDir, constant.LinkSuffix)
}

// 启用开机启动（使用 wails3 内置 API）
func EnableAutostart() (bool, error) {
	app := application.Get()
	if app == nil {
		return false, fmt.Errorf("application not initialized")
	}
	err := app.Autostart.Enable()
	if err != nil {
		fmt.Println("启用开机启动失败:", err)
		return false, err
	}
	return true, nil
}

// 禁用开机启动（使用 wails3 内置 API）
func DisableAutostart() bool {
	app := application.Get()
	if app == nil {
		return false
	}
	err := app.Autostart.Disable()
	if err != nil {
		fmt.Println("禁用开机启动失败:", err)
		return false
	}
	return true
}

// 迁移存量用户的开机启动设置
// 检查是否存在旧的快捷方式，如果存在则删除并使用新 API 重新设置
func MigrateAutostart() {
	if !IsWindows() {
		return
	}

	// 处理当前版本的快捷方式
	if _, err := os.Stat(linkPath); err == nil {
		fmt.Println("检测到旧的开机启动快捷方式，开始迁移...")

		// 删除旧的快捷方式
		err := os.Remove(linkPath)
		if err != nil {
			fmt.Println("删除旧快捷方式失败:", err)
			return
		}

		// 使用新 API 重新启用开机启动
		app := application.Get()
		if app == nil {
			fmt.Println("application not initialized, 无法迁移")
			return
		}

		err = app.Autostart.Enable()
		if err != nil {
			fmt.Println("重新设置开机启动失败:", err)
			return
		}

		fmt.Println("开机启动迁移成功")
	}

	// 处理旧版本 push-me-client 的快捷方式
	oldLinkPath := strings.Replace(linkPath, "push-me", "push-me-client", 1)
	if _, err := os.Stat(oldLinkPath); err == nil {
		fmt.Println("检测到旧版本 push-me-client 的快捷方式，开始清理...")

		err := os.Remove(oldLinkPath)
		if err != nil {
			fmt.Println("删除旧版本快捷方式失败:", err)
			return
		}

		fmt.Println("旧版本快捷方式清理成功")
	}
}

// 卸载应用（删除所有快捷方式）
func Uninstall() (bool, error) {
	if !IsWindows() {
		return false, nil
	}

	// 使用 wails3 API 禁用开机启动
	app := application.Get()
	if app != nil {
		err := app.Autostart.Disable()
		if err != nil {
			fmt.Println("禁用开机启动失败:", err)
		}
	}

	// 删除旧的启动快捷方式（兼容存量用户）
	err := os.Remove(linkPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("删除启动快捷方式失败:", err)
		return false, err
	}

	// 删除旧版本 push-me-client 的快捷方式
	oldLinkPath := strings.Replace(linkPath, "push-me", "push-me-client", 1)
	err = os.Remove(oldLinkPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("删除旧版本快捷方式失败:", err)
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
