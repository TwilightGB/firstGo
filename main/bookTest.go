package main

import (
	"firstGo/database"
	"firstGo/entity"
	"firstGo/goUtils"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func main() {
	//proxy.New()
	//createTable()
	doParse()
	//getIps()
}

func getIps() {
	ipRegexp := regexp.MustCompile(`https?://([\w]*:[\w]*@)?[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+`)
	doc := goUtils.GetPageFromUrlOrigin("http://ip.zdaye.com/FreeIPlist.html")
	str := doc.Find("#ipc tr td").Text()
	list := ipRegexp.FindAllString(str, -1)
	for _, value := range list {
		fmt.Println(value)
	}
}

func doParse() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	startTime := time.Now()
	book := entity.NewBook()
	var page = 20
	var url = book.UrlRule()
	for i := 0; i < 3; i++ {
		ch1 := make(chan string, 32)
		doc := goUtils.GetPageFromUrl(url)
		//list := make([]string, 0)
		go func() {
			book.SubjectItems(doc, func(image string) {
				ch1 <- image
				//list = append(list, image)
			})
			close(ch1)
		}()
		for value := range ch1 {
			valueS := strings.Split(value, "/")
			if database.IsExist("SbuId", valueS[len(valueS)-2]) {
				break
			}
			newBook := entity.AppendDetail(value, entity.NewBook())
			chWrite := make(chan *entity.Book, 32)
			go func() {
				chWrite <- newBook
				close(chWrite)
			}()
			go func() {
				bookRead := <-chWrite
				if !database.IsExist("isbn", bookRead.ISBN) {
					if err := database.Create(bookRead); err != nil {
						return
					}
				}
			}()
		}
		url = "https://book.douban.com/tag/Programming?" + book.PageRule(page)
		page += 20
	}
	elapsed := time.Since(startTime)
	fmt.Println(elapsed)
}
