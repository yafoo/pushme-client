package db

import (
	"log"
	"regexp"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var TextTypes *[]string = &[]string{"text", "markdown", "html", "url", ""}
var DataTypes *[]string = &[]string{"data", "markdata", "chart", "echarts", "svg", "note"}
var ChartTypes *[]string = &[]string{"bar", "line", "pie"}

const TextTypeDefault string = "text"
const ChartTypeDefault string = "bar"

var TextMap = make(map[string]bool)
var DataMap = make(map[string]bool)
var ChartMap = make(map[string]bool)

func init() {
	for _, t := range *TextTypes {
		TextMap[t] = true
	}
	for _, t := range *DataTypes {
		DataMap[t] = true
	}
	for _, t := range *ChartTypes {
		ChartMap[t] = true
	}
}

type Msg struct {
	ID      uint   `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Title   string `json:"title" gorm:"type:varchar(255)"`
	Content string `json:"content"`
	Date    string `json:"date" gorm:"type:varchar(255);default:"`
	Type    string `json:"type" gorm:"type:varchar(255);default:text"`
	Theme   string `json:"theme" gorm:"type:varchar(10);default:"`
	User    string `json:"user" gorm:"type:varchar(255);default:;index"`
	Face    string `json:"face" gorm:"type:varchar(255);default:"`
	PushKey string `json:"push_key" gorm:"-"`
}

func (Msg) TableName() string {
	return "msgs"
}

func (msg *Msg) IsTextMsg() bool {
	return TextMap[msg.Type]
}

func (msg *Msg) IsDataMsg() bool {
	return DataMap[msg.Type]
}

type Plugin struct {
	ID      int    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Title   string `json:"title" gorm:"type:varchar(255)"`
	Content string `json:"content"`
	Sort    int    `json:"sort" gorm:"default:0"`
	State   int    `json:"state" gorm:"default:1"`
}

var Db *gorm.DB
var err error

func InitDb(dbPath string) {
	Db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database " + err.Error())
	}
	Db.AutoMigrate(&Msg{})
	Db.AutoMigrate(&Plugin{})

	// 迁移存量数据：回填 user/face 字段
	MigrateUserField()
}

// MigrateUserField 为存量消息解析 theme/user/face，清理标题中的标记
func MigrateUserField() {
	// 只查需要迁移的：有分组或频道标记的
	var msgs []Msg
	Db.Where("title LIKE '%[#%' OR title LIKE '%[~%'").Find(&msgs)
	if len(msgs) == 0 {
		return
	}

	log.Printf("迁移 %d 条消息的分组数据...", len(msgs))

	// 用事务批量更新
	Db.Transaction(func(tx *gorm.DB) error {
		for i := range msgs {
			info := parseTitleForMigration(msgs[i].Title)
			if info.Theme != "" || info.User != "" || info.Title != msgs[i].Title {
				tx.Model(&msgs[i]).Updates(map[string]interface{}{
					"theme": info.Theme,
					"user":  info.User,
					"face":  info.Face,
					"title": info.Title,
				})
			}
		}
		return nil
	})
}

type titleInfoSimple struct {
	Theme string
	User  string
	Face  string
	Title string
}

// parseTitleForMigration 迁移用的标题解析（避免循环引用 msg 包）
func parseTitleForMigration(input string) titleInfoSimple {
	result := titleInfoSimple{Title: strings.TrimSpace(input)}
	title := result.Title

	// 提取 theme: [s] [i] [w] [f]
	themeRe := regexp.MustCompile(`^\[([iswfISWF])\]\s*(.*)`)
	if match := themeRe.FindStringSubmatch(title); len(match) > 0 {
		result.Theme = strings.ToLower(match[1])
		title = match[2]
	}

	// 移除频道标记: [~channel]
	channelRe := regexp.MustCompile(`\[~([^\]]+)\]`)
	title = channelRe.ReplaceAllString(title, "")

	// 提取分组: [#user!face]
	groupRe := regexp.MustCompile(`\[#([^!\]\n]+)(?:!([^\]\n]+))?\]`)
	if match := groupRe.FindStringSubmatch(title); len(match) > 1 {
		result.User = strings.TrimSpace(match[1])
		if len(match) > 2 {
			result.Face = strings.TrimSpace(match[2])
		}
		title = groupRe.ReplaceAllString(title, "")
	}

	result.Title = strings.TrimSpace(title)
	return result
}
