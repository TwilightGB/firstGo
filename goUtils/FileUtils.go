package goUtils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//get response body form url
func GetPageFromUrl(url string) *goquery.Document {
	res, error := http.Get(url)
	if error != nil {
		fmt.Print(error)
	}
	doc, error2 := goquery.NewDocumentFromReader(res.Body)
	if error2 != nil {
		fmt.Println(error2)
	}
	return doc
}

func IsExist(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

//download from url
func DownLoad(filepath string, url string) error {
	trueOrFalse := IsExist(filepath)
	if trueOrFalse {
		fmt.Println("have exists")
		return errors.New("have exists")
	}
	file, errorCreate := os.Create(filepath)
	if errorCreate != nil {
		return errorCreate
	}
	defer file.Close()

	resp, errorGet := http.Get(url)
	if errorGet != nil {
		return errorGet
	}
	defer resp.Body.Close()

	pix, _ := ioutil.ReadAll(resp.Body)
	_, errorCopy := io.Copy(file, bytes.NewReader(pix))
	if errorCopy != nil {
		return errorCopy
	}
	return nil
}

//int config file
func InitConfig(path string) map[string]string {
	var resultMap = make(map[string]string)
	file, error := os.Open(path)
	if error != nil {
		panic(error)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		byte, _, error := reader.ReadLine()
		if error != nil {
			if error == io.EOF {
				break
			}
			panic(error)
		}
		sLine := strings.TrimSpace((string(byte)))
		//
		index := strings.Index(sLine, "=")
		if index < 0 {
			continue
		}

		key := strings.TrimSpace(sLine[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(sLine[index+1 : len(sLine)])
		if len(value) == 0 {
			continue
		}
		resultMap[key] = value
	}
	return resultMap
}

//write file by cache
func WriteFile(path string, context string) {
	trueOrfalse := IsExist(path)
	if trueOrfalse {
		file := doWrite(path, context)
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
		doWrite(path, context)
		defer file.Close()
	}
}

func doWrite(path string, context string) *os.File {
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
