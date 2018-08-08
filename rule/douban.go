package rule

import (
	"strings"
	"net/url"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

type Douban struct {
	url string
	filepath string
	page int
}

func (d *Douban) UrlRule() (url string){
	return "https://book.douban.com/tag/Programming?start=20&type=T"
}
func (d *Douban) PageRule(currentPage int) (page string) {
	return "start="+strconv.Itoa(currentPage+20)+"&type=T"
}
func (d *Douban) ImageRule(doc *goquery.Document, f func(image string)){

	doc.Find(".pic").Each(func(i int, contentSelection *goquery.Selection) {
		imgpath, exist := contentSelection.Find("img").Attr("src")
		if !exist {
			return
		}
		f(imgpath)
	})
}
func (d *Douban) ParseUrl(imageUrl string) string {
	lastIndex:=strings.LastIndex(imageUrl, "/")
	name:=imageUrl[lastIndex+1: len(imageUrl)]
	u, err := url.Parse(imageUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Path)
	return  name;
}