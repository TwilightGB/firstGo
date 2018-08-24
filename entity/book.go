package entity

import (
	"firstGo/goUtils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type Book struct {
	gorm.Model
	id           int     `gorm:"primary_key;AUTO_INCREMENT"`
	Title        string  `gorm:"type:varchar(20)"`
	SubTitle     string  `gorm:"type:varchar(20)"`
	Author       string  `gorm:"type:varchar(20)"`
	OriginAuth   string  `gorm:"type:varchar(20)"`
	TransAuth    string  `gorm:"type:varchar(20)"`
	Bind         string  `gorm:"type:varchar(20)"`
	Year         string  `gorm:"type:varchar(20)"`
	Page         int     `gorm:"type:int(10)"`
	Price        string  `gorm:"type:varchar(20)"`
	Rating       float64 `gorm:"type:float(20)"`
	ISBN         string  `gorm:"type:varchar(20)"`
	Ratingpeople string  `gorm:"type:varchar(20)"`
	Publishcomp  string  `gorm:"type:varchar(55)"`
	Content      string  `gorm:"type:varchar(1024)"`
	Image        string  `gorm:"type:varchar(256)"`
	Classify     string  `gorm:"type:varchar(32)"`
	SbuId        string  `gorm:"type:varchar(32)"`
}

func NewBook() *Book {
	return &Book{}
}

//func InitTable() {
//	database.DB()
//	table := database.New()
//	table.SetTableName("testDouban")
//	table.AddColumn("Title,varchar(20)", "SUBTITLE,varchar(20)", "AUTHOR,varchar(20)", "ORIGINAUTH,varchar(20)", "TRANSAUTH,varchar(20)", "BIND,varchar(20)", "YEAR,varchar(20)", "PAGE,int(20)", "PRICE,varchar(20)", "RATING,int(20)", "ISBN,varchar(20)", "RATINGPEOPLE,varchar(20)", "PUBLISHCOMP,varchar(36)", "CONTENT,varchar(256)", "IMAGE,varchar(256)", )
//	table.Create()
//	database.DbClose()
//}

func (d *Book) UrlRule() (url string) {
	return "https://book.douban.com/tag/Programming?start=20&type=T"
}
func (d *Book) PageRule(currentPage int) (page string) {
	return "start=" + strconv.Itoa(currentPage+20) + "&type=T"
}
func (d *Book) SubjectItems(doc *goquery.Document, f func(image string)) {

	doc.Find(".subject-item .info h2").Each(func(i int, contentSelection *goquery.Selection) {
		imgpath, exist := contentSelection.Find("a").Attr("href")
		if !exist {
			return
		}
		f(imgpath)
	})
}

func AppendDetail(value string, book *Book) *Book {
	valueSplic := strings.Split(value, "/")
	book.SbuId = valueSplic[len(valueSplic)-2]
	docDetail := goUtils.GetPageFromUrl(value)
	bookTitle := strings.TrimSpace(docDetail.Find("h1").Text())
	book.Title = bookTitle
	//r.ReplaceAllString(bookName, "")
	bookInfo := strings.TrimSpace(docDetail.Find("#info").Text())
	bookInfo = strings.Replace(bookInfo, "\n ", "", -1)
	infoArray := strings.Split(bookInfo, "\n")
	err := appendInfo(infoArray, book)
	bookRating := strings.TrimSpace(docDetail.Find(".rating_num").Text())
	rating, err := strconv.ParseFloat(bookRating, 64)
	if err != nil {
		log.Print(err)
	} else {
		book.Rating = rating

	}
	bookRateNum := strings.TrimSpace(docDetail.Find(".rating_people").Text())
	book.Ratingpeople = bookRateNum
	bookImage, exist := docDetail.Find("#mainpic img").Attr("src")
	if exist {
		book.Image = bookImage
	}
	bookContext := docDetail.Find(".intro p").Text()
	book.Content = bookContext
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	fmt.Println(book)
	return book
}

func appendInfo(infoArray []string, book *Book) error {
	mapInfo := goUtils.ParseDate(infoArray)
	book.SubTitle = mapInfo["副标题"]
	book.Author = mapInfo["作者"]
	book.OriginAuth = mapInfo["原作名"]
	book.TransAuth = mapInfo["译者"]
	book.Bind = mapInfo["装帧"]
	book.Year = mapInfo["出版年"]
	book.Price = mapInfo["定价"]
	book.ISBN = mapInfo["ISBN"]
	book.Publishcomp = mapInfo["出版社"]
	book.Classify = mapInfo["丛书"]
	page, err := strconv.Atoi(mapInfo["页数"])
	if err != nil {
		log.Println(err)
	} else {
		book.Page = page
	}
	return err
}
