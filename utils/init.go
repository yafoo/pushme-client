package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"PushMeClient/constant"
)

var (
	AppPath     string
	UserDir     string
	RootDir     string
	LogoPath    string
	SettingPath string
	DbPath      string
)

func init() {
	userInfo, err := user.Current()
	if nil != err {
		fmt.Println(err)
	}
	UserDir = userInfo.HomeDir
	RootDir = path.Join(UserDir, constant.AppDir)
	if _, err := os.Stat(RootDir); os.IsNotExist(err) {
		err = os.Mkdir(RootDir, os.ModeDir)
		if err != nil {
			fmt.Println(err)
		}
	}
	LogoPath = path.Join(RootDir, "logo.png")
	SettingPath = path.Join(RootDir, "setting.json")
	DbPath = path.Join(RootDir, "data.db")

	AppPath, err = os.Executable()
	if err != nil {
		fmt.Println("获取AppPath失败", err)
	}
}
