package main

import (
	"firstGo/entity"
	"firstGo/goUtils"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
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
func main() {
	//createTable()
	startTime := time.Now()
	book := entity.NewBook()
	var page = 20
	var url = book.UrlRule()
	for i := 0; i < 3; i++ {
		list := getPageList(book, url)
		parseBook(list)
		url = "https://book.douban.com/tag/Programming?" + book.PageRule(page)
		page += 20
	}
	elapsed := time.Since(startTime)
	fmt.Println(elapsed)
}

func parseBook(list []string) {
	for _, value := range list {
		newBook := entity.AppendDetail(value, entity.NewBook())
		newBook.CreatedAt = time.Now()
		newBook.UpdatedAt = time.Now()
		if err := db.Create(newBook).Error; err != nil {
			return
		}
		//fmt.Println(mapInfo)
	}
}

func getPageList(book *entity.Book, url string) []string {
	doc := goUtils.GetPageFromUrl(url)
	list := make([]string, 0)
	book.SubjectItems(doc, func(image string) {
		list = append(list, image)
	})
	return list
}