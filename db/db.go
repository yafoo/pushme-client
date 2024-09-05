package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Msg struct {
	ID      uint   `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Title   string `json:"title" gorm:"type:varchar(255)"`
	Content string `json:"content"`
	Date    string `json:"date" gorm:"type:varchar(255);default:"`
	Type    string `json:"type" gorm:"type:varchar(255);default:text"`
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
		panic("failed to connect database")
	}
	Db.AutoMigrate(&Msg{})
	Db.AutoMigrate(&Plugin{})
}
