package services

import (
	"PushMe/constant"
	"PushMe/internal/request"
	"PushMe/internal/setting"
	"PushMe/internal/toast"
	"PushMe/internal/utils"
	"log"
)

type SettingService struct{}

func (s *SettingService) GoGetSetting() setting.SettingType {
	return setting.Setting
}

func (s *SettingService) GoGetSettingDefault() setting.SettingType {
	return setting.GetSettingDefault()
}

func (s *SettingService) GoSaveSetting(settingData setting.SettingType) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			utils.Toast(r.(string))
		}
	}()

	setting.SaveSetting(settingData)
	utils.Autostart(settingData.System.Enable)

	return true
}

func (s *SettingService) GoGetHost() setting.HostType {
	return setting.Setting.Host
}

func (s *SettingService) GoGetHtmlJs() bool {
	return setting.Setting.Other.HtmlJs
}

func (s *SettingService) GoGetIps() []string {
	return utils.GetIps()
}

func (s *SettingService) GoRestart() {
	utils.Restart()
}

func (s *SettingService) GoGetVersion() string {
	return constant.AppVersion
}

func (s *SettingService) GoCheckVersion() {
	go CheckVersion()
}

func CheckVersion() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			toast.Toast(r.(string))
		}
	}()

	req := request.Request{
		Url: constant.UrlVersion,
		Data: map[string]interface{}{
			"version": constant.AppVersion,
		},
	}
	result := req.Post()
	utils.EventEmit("version", result)
}
