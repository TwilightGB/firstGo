package database

import (
	"firstGo/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		panic(err)
	}
}
func createTable() {
	if !db.HasTable(entity.NewBook()) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB AUTO_INCREMENT=1 CHARSET=utf8").CreateTable(entity.NewBook()).Error; err != nil {
			panic(err)
		}
	}
}

func IsExist(column string, value string) bool {
	d := db.Where(column+"=?", value).Find(&entity.Book{})
	if d.RecordNotFound() {
		return false
	}
	return true
}

func Create(var1 interface{}) error {
	return db.Create(var1).Error
}
