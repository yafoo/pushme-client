package services

import (
	"PushMe/internal/setting"
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
