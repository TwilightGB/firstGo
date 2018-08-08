package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"firstGo/rule"
	"firstGo/utils"
)

type image struct {
	surl string
	filepath string
	page int
	ImageRule rule.Rule
}

func (im *image) getImage(surl string) []string {
	doc, err := goquery.NewDocument(surl)
	if err != nil {
		fmt.Println(err)
	}
	imageList := make([]string, 0)
	im.ImageRule.ImageRule(doc, func(image string) {
		imageList = append(imageList, image)
	})
	return imageList
}

//, strconv.FormatInt(time.Now().Unix(), 10)
//strconv.Itoa(rand.Intn(10)) + ".jpg"
func main() {
	var resultList = make([]string,0)
	//images:=&image{surl:"",filepath:"",page:1,ImageRule:&rule.Douban{}}
	//fmt.Println(images)
	image :=new(image)
	image.surl="https://book.douban.com/tag/Programming?"
	image.filepath="D:\\temp\\"
	image.page=20
	image.ImageRule=&rule.Douban{}
	for i:=0;i<3;i++{
		lnewUrl:=image.ImageRule.PageRule(image.page)
		result:=image.getImage(image.surl+lnewUrl)
		for _,datau:= range result{
			resultList=append(resultList,datau)
		}
		image.page+=20
	}
	for _,data:= range resultList{
		fmt.Println(data)
		sufixx:=image.ImageRule.ParseUrl(data)
		goUtils.DownLoad(image.filepath+sufixx,data)
	}
}
