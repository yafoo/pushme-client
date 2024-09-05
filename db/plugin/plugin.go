package plugin

import (
	"PushMeClient/db"

	"gorm.io/gorm"
)

func Add(data *db.Plugin) (id int) {
	if data.ID == 0 {
		plugin := &db.Plugin{
			Title:   data.Title,
			Content: data.Content,
			Sort:    Count(),
		}
		db.Db.Save(plugin)
		id = plugin.ID
	} else {
		db.Db.Save(data)
		id = data.ID
	}
	return id
}

func Get(id int) db.Plugin {
	var data db.Plugin
	db.Db.First(&data, id)
	return data
}

func Del(id int) bool {
	var data db.Plugin
	db.Db.First(&data, id)
	if data.ID == 0 {
		return true
	}
	db.Db.Delete(&db.Plugin{}, id)
	db.Db.Model(&db.Plugin{}).Where("sort > ?", data.Sort).Update("sort", gorm.Expr("sort - ?", 1))
	return true
}

func List() []db.Plugin {
	var list []db.Plugin
	db.Db.Order("sort").Order("id").Find(&list)
	return list
}

func ListWithState(state int) []db.Plugin {
	var list []db.Plugin
	db.Db.Where("state == ?", state).Order("sort").Order("id").Find(&list)
	return list
}

func Count() int {
	var count int64
	err := db.Db.Model(&db.Plugin{}).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

func CountWithState(state int) int {
	var count int64
	err := db.Db.Model(&db.Plugin{}).Where("state == ?", state).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

func ToggleState(data *db.Plugin) bool {
	if data.State == 1 {
		data.State = 0
	} else {
		data.State = 1
	}
	db.Db.Save(data)
	return true
}

func SwapSort(data1 *db.Plugin, data2 *db.Plugin) bool {
	sort1 := data1.Sort
	data1.Sort = data2.Sort
	data2.Sort = sort1
	db.Db.Save(data1)
	db.Db.Save(data2)
	return true
}
