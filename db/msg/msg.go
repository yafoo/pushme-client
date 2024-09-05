package msg

import "PushMeClient/db"

var typeKey = "type in (?)"
var textWhere = []string{"text", "markdown", ""}
var dataWhere = []string{"data", "markdata", "chart"}

func Add(data *db.Msg) db.Msg {
	if data.Type == "data" || data.Type == "markdata" || data.Type == "chart" {
		res := db.Msg{}
		db.Db.Where("title like ?", data.Title).Order("id desc").First(&res)
		if res.ID > 0 {
			res.Title = data.Title
			res.Content = data.Content
			res.Type = data.Type
			res.Date = data.Date
			db.Db.Save(res)
			return res
		}
	}
	db.Db.Save(data)
	return *data
}

func Get(id int) db.Msg {
	var data db.Msg
	db.Db.First(&data, id) //id不存在时报错
	return data
}

func Del(id int) bool {
	db.Db.Delete(&db.Msg{}, id)
	return true
}

func ClearText() bool {
	db.Db.Where("type IN ?", []string{"", "text", "markdown"}).Delete(&db.Msg{}, "type IN ?", []string{"", "text", "markdown"})
	return true
}

func MsgList(page int, pageSize int) []db.Msg {
	offset := (page - 1) * pageSize
	var list []db.Msg
	db.Db.Where(typeKey, textWhere).Offset(offset).Limit(pageSize).Order("id desc").Find(&list)
	return list
}

func DataList() []db.Msg {
	var list []db.Msg
	db.Db.Where(typeKey, dataWhere).Find(&list)
	return list
}

func CountText() int {
	var count int64
	err := db.Db.Model(&db.Msg{}).Where("type IN ?", []string{"", "text", "markdown"}).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}
