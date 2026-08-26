package services

import (
	"PushMe/internal/api"
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

	// 保存旧设置，用于对比
	oldSetting := setting.Setting

	setting.SaveSetting(settingData)
	utils.Autostart(settingData.System.Enable)

	// 检测 API 相关配置是否变化，变化则热重启
	applyApiChanges(oldSetting, settingData)

	utils.EventEmit("setting:change", settingData)

	return true
}

// applyApiChanges 检测 API 配置变化并执行相应操作
func applyApiChanges(old, new setting.SettingType) {
	wasRunning := old.Api.Enable
	shouldBeRunning := new.Api.Enable

	if wasRunning && !shouldBeRunning {
		// API 从启用变为关闭 → 停止服务
		log.Println("API 已关闭，停止 HTTP 服务...")
		go api.Stop()
	} else if !wasRunning && shouldBeRunning {
		// API 从关闭变为启用 → 启动服务
		log.Println("API 已启用，启动 HTTP 服务...")
		go api.Start()
	} else if wasRunning && shouldBeRunning {
		// 都在运行，检查 ip/port 是否变化 → 重启服务
		if old.Api.Ip != new.Api.Ip || old.Api.Port != new.Api.Port {
			log.Println("API 地址/端口变更，重启 HTTP 服务...")
			go api.Restart()
		}
	}
}

func (s *SettingService) GoGetHost() setting.HostType {
	return setting.Setting.Host
}

func (s *SettingService) GoGetHtmlJs() bool {
	return setting.Setting.Other.HtmlJs
}

func (s *SettingService) GoGetSettingNotice() setting.NoticeType {
	return setting.Setting.Notice
}
