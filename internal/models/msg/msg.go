package msg

import (
	db "PushMe/internal/models"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

var typeKey = "type in (?)"

// 标题解析正则
var (
	themePattern  = regexp.MustCompile(`^\[([iswfISWF])\]\s*(.*)`)
	channelPattern = regexp.MustCompile(`\[~([^\]]+)\]`)
	groupPattern  = regexp.MustCompile(`\[#([^!\]\n]+)(?:!([^\]\n]+))?\]`)
)

// TitleInfo 标题解析结果
type TitleInfo struct {
	Theme   string
	Title   string
	User    string
	Face    string
	Channel string
}

// CalcTitleInfo 从标题中提取主题、分组等信息
func CalcTitleInfo(input string) TitleInfo {
	info := TitleInfo{}
	if strings.TrimSpace(input) == "" {
		return info
	}

	title := strings.TrimSpace(input)

	// 提取主题
	if match := themePattern.FindStringSubmatch(title); len(match) > 0 {
		info.Theme = strings.ToLower(match[1])
		title = match[2]
	}

	// 提取频道
	if match := channelPattern.FindStringSubmatch(title); len(match) > 0 {
		info.Channel = strings.TrimSpace(match[1])
		title = strings.Replace(title, match[0], "", 1)
	}

	// 提取分组
	if match := groupPattern.FindStringSubmatch(title); len(match) > 0 {
		info.User = strings.TrimSpace(match[1])
		if len(match) > 2 {
			info.Face = strings.TrimSpace(match[2])
		}
		title = strings.Replace(title, match[0], "", 1)
	}

	info.Title = strings.TrimSpace(title)
	return info
}

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
	// 从标题中解析分组和主题信息
	titleInfo := CalcTitleInfo(data.Title)
	data.Theme = titleInfo.Theme
	data.User = titleInfo.User
	data.Face = titleInfo.Face
	data.Title = titleInfo.Title // 保存干净的标题

	if data.IsDataMsg() {
		res := db.Msg{}
		db.Db.Where("title like ?", data.Title).Order("id desc").First(&res)
		if res.ID > 0 {
			// 便签新增不允许标题重复
			// if data.Type == "note" {
			// 	return db.Msg{}
			// }
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

// MsgListByUser 按用户/分组查询消息列表
func MsgListByUser(user string, page int, pageSize int) []db.Msg {
	offset := (page - 1) * pageSize
	var list []db.Msg
	db.Db.Where(typeKey, *db.TextTypes).Where("user = ?", user).Offset(offset).Limit(pageSize).Order("id desc").Find(&list)
	return list
}

// CountByUser 统计某用户的消息数量
func CountByUser(user string) int {
	var count int64
	db.Db.Model(&db.Msg{}).Where(typeKey, *db.TextTypes).Where("user = ?", user).Count(&count)
	return int(count)
}

// DelByUser 删除某用户的所有消息
func DelByUser(user string) bool {
	result := db.Db.Where("user = ?", user).Delete(&db.Msg{})
	return result.Error == nil
}

// UserList 获取所有有分组的消息的用户列表（含最近消息信息和未读数量）
func UserList() []map[string]interface{} {
	var results []map[string]interface{}
	db.Db.Model(&db.Msg{}).
		Select(`user, face, count(*) as count, max(id) as last_id`).
		Where(typeKey, *db.TextTypes).
		Where("user != ''").
		Group("user").
		Order("last_id desc").
		Find(&results)
	return results
}

// MsgListGrouped 首页分组列表：每个 user 只显示最新一条 + 所有无 user 的消息，按 id 降序
func MsgListGrouped(page int, pageSize int) []db.Msg {
	offset := (page - 1) * pageSize
	var list []db.Msg
	tableName := db.Msg{}.TableName()

	// 用子查询取每个 user 的最新消息 id，再 UNION 无 user 的消息
	db.Db.Raw(`
		SELECT * FROM `+"`"+tableName+"`"+`
		WHERE id IN (
			SELECT MAX(id) FROM `+"`"+tableName+"`"+`
			WHERE user != '' AND user IS NOT NULL AND type IN (?)
			GROUP BY user
		)
		UNION ALL
		SELECT * FROM `+"`"+tableName+"`"+`
		WHERE (user = '' OR user IS NULL) AND type IN (?)
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, *db.TextTypes, *db.TextTypes, pageSize, offset).Scan(&list)

	return list
}

// CountGrouped 分组后的总条数（用于判断是否有更多）
func CountGrouped() int {
	var count int64

	// 每个 user 算 1 条 + 无 user 的消息数
	var userCount int64
	db.Db.Model(&db.Msg{}).
		Where(typeKey, *db.TextTypes).
		Where("user != '' AND user IS NOT NULL").
		Distinct("user").
		Count(&userCount)

	var noUserCount int64
	db.Db.Model(&db.Msg{}).
		Where(typeKey, *db.TextTypes).
		Where("user = '' OR user IS NULL").
		Count(&noUserCount)

	count = userCount + noUserCount
	return int(count)
}
