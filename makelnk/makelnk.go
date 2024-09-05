package makelnk

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"PushMeClient/constant"
	"PushMeClient/server/html"
	"PushMeClient/utils"
)

//go:embed makelnk/*
var Makelnk embed.FS

func initMakelnk() error {
	files, err := Makelnk.ReadDir("makelnk")
	if err != nil {
		return err
	}
	makeLnkDir := path.Join(utils.RootDir, "makelnk")
	for _, file := range files {
		content, err := Makelnk.ReadFile("makelnk/" + file.Name())
		if err != nil {
			return err
		}
		err = os.WriteFile(makeLnkDir+"/"+file.Name(), content, 0755)
		if err != nil {
			return err
		}
	}

	content, err := html.Static.ReadFile("static/logo.png")
	if err != nil {
		return err
	}
	err = os.WriteFile(utils.LogoPath, content, 0755)
	if err != nil {
		return err
	}

	return nil
}

func RegisterAppUserModeId() {
	go func() {
		makeLnkDir := path.Join(utils.RootDir, "makelnk")
		if _, err := os.Stat(makeLnkDir); os.IsNotExist(err) {
			err = os.Mkdir(makeLnkDir, os.ModeDir)
			if err != nil {
				fmt.Println(err)
			}
			err = initMakelnk()
			if err != nil {
				fmt.Println("初始化makelnk失败", err)
			}
		}

		lnkPath := path.Join(utils.UserDir, constant.StartMenuLnk)
		if _, err := os.Stat(lnkPath); os.IsNotExist(err) {
			makeLnkExe := path.Join(utils.RootDir, "makelnk", "makelnk.exe")
			cmd := exec.Command(makeLnkExe, constant.AppID, utils.AppPath, constant.LnkName)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			err = cmd.Start()
			if err != nil {
				fmt.Println("注册AppID失败", err)
			}
			fmt.Println("注册AppID成功")
		}
	}()

}
