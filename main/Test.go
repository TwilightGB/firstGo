package main

import (
	"bufio"
	"firstGo/goUtils"
	"firstGo/rule"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"os"
	"runtime"
)

type image struct {
	surl      string
	filepath  string
	page      int
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan string, 100)
	go func(ch chan string) {
		uri := "E:\\mygo\\src\\firstGo\\config.properties"
		file, error := os.Open(uri)
		if error != nil {
			panic(error)
		}
		defer file.Close()
		r := bufio.NewReader(file)
		for {
			byte, _, error := r.ReadLine()
			if error != nil {
				if error == io.EOF {
					break
				}
				panic(error)
			}
			ch <- string(byte)
		}
	}(ch)
	f, error := os.Create("E:\\e")
	if error != nil {
		panic(error)
	}
	defer f.Close()
	x, ok := <-ch
	if ok != false {
		fmt.Println(x)
	}
	dMap := goUtils.InitConfig("E:\\mygo\\src\\firstGo\\config.properties")
	var resultList = make([]string, 0)
	//images:=&image{surl:"",filepath:"",page:1,ImageRule:&rule.Douban{}}
	//fmt.Println(images)
	image := new(image)
	image.surl = dMap["purl"]
	image.filepath = dMap["filepath"]
	image.page = 20
	image.ImageRule = &rule.Douban{}
	for i := 0; i < 3; i++ {
		lnewUrl := image.ImageRule.PageRule(image.page)
		result := image.getImage(image.surl + lnewUrl)
		for _, datau := range result {
			resultList = append(resultList, datau)
		}
		image.page += 20
	}
	for _, data := range resultList {
		fmt.Println(data)
		sufixx := image.ImageRule.ParseUrl(data)
		goUtils.DownLoad(image.filepath+sufixx, data)
	}
}
