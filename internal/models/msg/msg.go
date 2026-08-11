package msg

import (
	db "PushMe/internal/models"
	"encoding/json"
	"strconv"
	"strings"
)

var typeKey = "type in (?)"

type ChartData struct {
	Type  string    `json:"type"`
	Size  uint      `json:"size"`
	List  []float64 `json:"list"`
	Label []string  `json:"label"`
}

func parseChartData(str string) ChartData {
	data := ChartData{
		Type:  db.ChartTypeDefault, // 默认类型
		Size:  10,                  // 默认大小
		List:  []float64{},
		Label: []string{},
	}

	parts := strings.Split(str, "::")
	if len(parts) > 1 {
		typeSize := strings.Split(parts[0], "|")
		if len(typeSize) > 1 {
			data.Type = typeSize[0]
			size, err := strconv.ParseUint(typeSize[1], 10, 64)
			if err == nil {
				data.Size = uint(size)
			}
		} else if db.ChartMap[typeSize[0]] {
			data.Type = typeSize[0]
		} else {
			size, err := strconv.ParseUint(typeSize[0], 10, 64)
			if err == nil {
				data.Size = uint(size)
			}
		}
		str = parts[1]
	}

	// 解析数据点和标签
	for _, s := range strings.Split(str, ",") {
		valueLabel := strings.Split(s, "|")
		if len(valueLabel) > 0 {
			value, err := strconv.ParseFloat(valueLabel[0], 64)
			if err == nil {
				data.List = append(data.List, value)
			} else {
				data.List = append(data.List, 0)
			}
			if len(valueLabel) > 1 {
				data.Label = append(data.Label, valueLabel[1])
			} else {
				data.Label = append(data.Label, "")
			}
		}
	}

	return data
}

func mergeChartData(first, second ChartData) ChartData {
	merged := ChartData{
		Type:  second.Type,
		Size:  second.Size,
		List:  make([]float64, 0),
		Label: make([]string, 0),
	}

	labelMap := make(map[string]int)
	for i, label := range first.Label {
		if label != "" {
			labelMap[label] = i
		}
	}

	merged.List = append(merged.List, first.List...)
	merged.Label = append(merged.Label, first.Label...)

	for i, label := range second.Label {
		if label != "" {
			if index, exists := labelMap[label]; exists {
				merged.List[index] = second.List[i]
			} else {
				merged.List = append(merged.List, second.List[i])
				merged.Label = append(merged.Label, label)
			}
		} else {
			merged.List = append(merged.List, second.List[i])
			merged.Label = append(merged.Label, label)
		}
	}

	if uint(len(merged.List)) > merged.Size {
		start := len(merged.List) - int(merged.Size)
		merged.List = merged.List[start:]
		merged.Label = merged.Label[start:]
	}

	return merged
}

func Add(data *db.Msg) db.Msg {
	if data.IsDataMsg() {
		res := db.Msg{}
		db.Db.Where("title like ?", data.Title).Order("id desc").First(&res)
		if res.ID > 0 {
			// 便签新增不允许标题重复
			if db.DataMap[data.Type] {
				return db.Msg{}
			}
			res.Title = data.Title
			res.Type = data.Type
			res.Date = data.Date
			if data.Type != "chart" {
				res.Content = data.Content
			} else {
				newData := parseChartData(data.Content)
				oldData := ChartData{}
				json.Unmarshal([]byte(res.Content), &oldData)
				mergeData := mergeChartData(oldData, newData)
				jsonBytes, err := json.Marshal(mergeData)
				if err == nil {
					res.Content = string(jsonBytes)
				}
			}
			db.Db.Save(res)
			return res
		} else if data.Type == "chart" {
			newData := parseChartData(data.Content)
			jsonBytes, err := json.Marshal(newData)
			if err == nil {
				data.Content = string(jsonBytes)
			}
		}
	} else if !data.IsTextMsg() {
		data.Type = db.TextTypeDefault
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

func Update(data *db.Msg) bool {
	result := db.Db.Save(data)
	return result.Error == nil
}

func ClearText() bool {
	db.Db.Where("type IN ?", *db.TextTypes).Delete(&db.Msg{}, "type IN ?", *db.TextTypes)
	return true
}

func MsgList(page int, pageSize int) []db.Msg {
	offset := (page - 1) * pageSize
	var list []db.Msg
	db.Db.Where(typeKey, *db.TextTypes).Offset(offset).Limit(pageSize).Order("id desc").Find(&list)
	return list
}

func DataList() []db.Msg {
	var list []db.Msg
	db.Db.Where(typeKey, *db.DataTypes).Find(&list)
	return list
}

func CountText() int {
	var count int64
	err := db.Db.Model(&db.Msg{}).Where("type IN ?", *db.TextTypes).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}
