package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var TextTypes *[]string = &[]string{"text", "markdown", "html", "url", ""}
var DataTypes *[]string = &[]string{"data", "markdata", "chart", "echarts", "svg"}
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
}
