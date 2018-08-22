package main

import (
	"firstGo/goUtils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func main() {
	var url string
	url = "http://www.zanghaihua.org/645.html"
	filepath := "D:\\temp\\cq.txt"
	var context string
	for i := 0; i < 225; i++ {
		doc := goUtils.GetPageFromUrl(url)
		context += " 【 " + doc.Find("h1").Text() + " 】 "
		doc.Find("#BookText").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			strings.TrimSpace(text)
			strings.Replace(text, "\n", "", -1)
			context += text
		})
		nextUrl, succ := doc.Find(".linkbtn a:last-child").Attr("href")
		if succ {
			context += "【" + strconv.Itoa(i) + "### " + nextUrl + "】"
			fmt.Println(i)
			url = nextUrl
		} else {
			fmt.Println(i)
			break
		}
		goUtils.WriteFile(filepath, context)
	}

}
