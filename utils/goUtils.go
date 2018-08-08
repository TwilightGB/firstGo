package goUtils

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"bytes"
	"errors"
	"bufio"
	"strings"
)

func IsExist(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}



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

func InitConfig(path string) map[string]string  {
	var resultMap = make(map[string]string)
	file,error:=os.Open(path)
	if error!=nil{
		panic(error)
	}
	defer file.Close()
	reader :=bufio.NewReader(file)
	for{
		byte,_,error:=reader.ReadLine()
		if error!=nil{
			if error==io.EOF{
				break
			}
			panic(error)
		}
		sLine :=strings.TrimSpace((string(byte)))
		//
		index:=strings.Index(sLine,":")
		if index<0{
			continue
		}

		key:=strings.TrimSpace(sLine[:index])
		if len(key)==0{
			continue
		}
		value:=strings.TrimSpace(sLine[index:])
		if len(value)==0{
			continue
		}
		resultMap[key]=value
	}
	return  resultMap
}