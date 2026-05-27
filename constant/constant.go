package constant

import (
	db "PushMe/internal/models"
	"log"
	"os"
	"os/user"
	"path"
)

const AppDir string = "push-me"
const AppID string = "PushMe"
const AppName string = "PushMe"
const AppVersion string = "v4.0.0"
const UrlVersion string = "https://push.i-i.me/api/version"
const WindowWidth int = 320
const WindowHeight int = 520

const (
	LnkName      string = AppID + ".lnk"
	LinkSuffix   string = "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup/" + LnkName
	StartMenuLnk string = "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/" + LnkName
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
		log.Println(err)
	}
	UserDir = userInfo.HomeDir
	RootDir = path.Join(UserDir, AppDir)
	log.Println("RootDir", RootDir)
	if _, err := os.Stat(RootDir); os.IsNotExist(err) {
		err = os.Mkdir(RootDir, os.ModeDir)
		if err != nil {
			log.Println(err)
		}
	}
	LogoPath = path.Join(RootDir, "logo.png")
	SettingPath = path.Join(RootDir, "setting.json")
	DbPath = path.Join(RootDir, "data.db")

	AppPath, err = os.Executable()
	if err != nil {
		log.Println("获取AppPath失败", err)
	}
	// DbPath = path.Join(filepath.Dir(AppPath), "data.db")

	log.Println("DbPath", DbPath)
	db.InitDb(DbPath)
}
