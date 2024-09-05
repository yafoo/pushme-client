package setting

import (
	"PushMeClient/utils"
	"encoding/json"
	"fmt"
	"os"
)

type SettingType struct {
	Api    ApiType    `json:"api"`
	Host   HostType   `json:"host"`
	Repost RepostType `json:"repost"`
	Notice NoticeType `json:"notice"`
	System SystemType `json:"system"`
}

type ApiType struct {
	Enable bool   `json:"enable"`
	Ip     string `json:"ip"`
	Port   string `json:"port"`
}

type HostType struct {
	Enable     bool   `json:"enable"`
	Protocol   string `json:"protocol"`
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	OfflineMsg bool   `json:"offline_msg"`
	PushKey    string `json:"push_key"`
}

type RepostType struct {
	Enable bool   `json:"enable"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

type NoticeType struct {
	Enable   bool   `json:"enable"`
	Duration string `json:"duration"`
	Audio    string `json:"audio"`
}

type SystemType struct {
	Enable bool `json:"enable"`
}

var Setting SettingType
var BaseApi string

func init() {
	if _, err := os.Stat(utils.SettingPath); os.IsNotExist(err) {
		initFile()
	}
	InitSetting()
}

func InitSetting() {
	data, err := os.ReadFile(utils.SettingPath)
	if err != nil {
		panic("读取配置文件出错：" + err.Error())
	}
	err = json.Unmarshal(data, &Setting)
	if err != nil {
		panic("解析配置文件json出错：" + err.Error())
	}

	BaseApi = "http://" + Setting.Api.Ip + ":" + Setting.Api.Port
}

func SaveSetting(setting SettingType) {
	saveFile(setting)
	InitSetting()
}

func GetSettingDefault() SettingType {
	api := ApiType{
		Enable: true,
		Ip:     "127.0.0.1",
		Port:   "3200",
	}
	host := HostType{
		Enable:     false,
		Protocol:   "ws",
		Ip:         "",
		Port:       "3010",
		OfflineMsg: false,
		PushKey:    "",
	}
	repost := RepostType{
		Enable: false,
		Url:    "",
		Method: "POST/JSON",
	}
	notice := NoticeType{
		Enable:   true,
		Duration: "short",
		Audio:    "default",
	}
	system := SystemType{
		Enable: false,
	}
	setting := SettingType{
		Api:    api,
		Host:   host,
		Repost: repost,
		Notice: notice,
		System: system,
	}

	return setting
}

func initFile() {
	setting := GetSettingDefault()
	saveFile(setting)
}

func saveFile(setting SettingType) {
	data, err := json.MarshalIndent(setting, "", "  ")
	if err != nil {
		panic("生成配置文件json数据出错：" + err.Error())
	}

	err = os.WriteFile(utils.SettingPath, data, 0644)
	if err != nil {
		panic("创建配置文件出错：" + err.Error())
	}
	fmt.Println("创建配置文件：", utils.SettingPath)
}
