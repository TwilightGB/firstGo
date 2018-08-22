package main

import (
	"bufio"
	"firstGo/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func writeFile(path string, context string) {
	trueOrfalse := goUtils.IsExist(path)
	if trueOrfalse {
		file := funcName(path, context)
		defer file.Close()
	} else {
		file, errorCreate := os.Create(path)
		if errorCreate != nil {
			return
		}
		// 查找文件末尾的偏移量
		// 从末尾的偏移量开始写入内容
		//n, _ := file.Seek(0, os.SEEK_END)
		//_, _ = file.WriteAt([]byte(context), n)
		funcName(path, context)
		defer file.Close()
	}
}

func funcName(path string, context string) *os.File {
	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	bufferedWriter := bufio.NewWriter(file)
	//bytesWritten, err := bufferedWriter.Write([]byte{})
	bytesWritten, err := bufferedWriter.WriteString(context)
	fmt.Printf("Bytes written: %d\n", bytesWritten)
	//unflushedBufferSize := bufferedWriter.Buffered()
	//log.Printf("Bytes buffered: %d\n", unflushedBufferSize)
	bufferedWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}
	return file
}
func main() {
	var url string
	url = "http://www.zanghaihua.org/645.html"
	filepath := "D:\\temp\\cq.txt"
	var context string
	for i := 0; i < 225; i++ {
		res, error := http.Get(url)
		if error != nil {
			fmt.Print(error)
		}
		doc, error2 := goquery.NewDocumentFromReader(res.Body)
		if error2 != nil {
			fmt.Println(error2)
		}
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
		writeFile(filepath, context)
	}

}
