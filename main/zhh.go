package main

import (
	"firstGo/goUtils"
	"firstGo/rule"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	startTime := time.Now() // get current time
	fileName := rand.Intn(1000)
	filepath := "D:\\temp\\" + strconv.Itoa(fileName) + ".txt"
	zhh := rule.InitZhh()
	zhh.FilePath = filepath
	var context string
	for i := 0; i < 225; i++ {
		var doc *goquery.Document
		doc = goUtils.GetPageFromUrl(zhh.Url)
		context += " 【 " + doc.Find("h1").Text() + " 】 "
		context += zhh.GetText(doc)

		nextUrl, succ := zhh.NextUrl(doc)
		if succ {
			context += "【" + strconv.Itoa(i) + "### " + nextUrl + "】"
			fmt.Println(i)
			zhh.Url = nextUrl
		} else {
			fmt.Println(i)
			break
		}
		goUtils.WriteFile(filepath, context)
	}
	elapsed := time.Since(startTime)
	fmt.Println(elapsed)
}
