package services

import (
	db "PushMe/internal/models"
	dbPlugin "PushMe/internal/models/plugin"
)

type PluginService struct{}

func (p *PluginService) GoGetPluginList() []db.Plugin {
	return dbPlugin.List()
}

func (p *PluginService) GoGetPlugin(id int) db.Plugin {
	return dbPlugin.Get(id)
}

func (p *PluginService) GoEditPlugin(plugin db.Plugin) int {
	return dbPlugin.Add(&plugin)
}

func (p *PluginService) GoDelPlugin(id int) bool {
	return dbPlugin.Del(id)
}

func (p *PluginService) GoEditPluginState(plugin db.Plugin) int {
	return dbPlugin.Add(&plugin)
}

func (p *PluginService) GoGetPluginCount() int {
	return dbPlugin.CountWithState(1)
}

func (p *PluginService) GoSwapPlugin(plugin1 db.Plugin, plugin2 db.Plugin) bool {
	return dbPlugin.SwapSort(&plugin1, &plugin2)
}

func (p *PluginService) GoGetPluginListWithState(state int) []db.Plugin {
	return dbPlugin.ListWithState(state)
}
