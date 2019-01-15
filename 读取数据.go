package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

/*(?s:(.*?))*/

func findLine(page chan string) {
	str := <-page
	ipret := regexp.MustCompile(`(?s:(.*?)) - - `)
	timeret := regexp.MustCompile(`\[(?s:(.*?))\]`)
	articleret := regexp.MustCompile(`"http://www.neusoft.com/article/(?s:(.*?))"`)
	videoret := regexp.MustCompile(`"http://www.neusoft.com/video/(?s:(.*?))"`)
	ip := ipret.FindAllStringSubmatch(str, -1)
	time := timeret.FindAllStringSubmatch(str, -1)
	article := articleret.FindAllStringSubmatch(str, -1)
	video := videoret.FindAllStringSubmatch(str, -1)
	if video == nil {
		a := [][]string{{" "}}
		video = a
	}
	if article == nil {
		a := [][]string{{" "}}
		article = a
	}
	f, err := os.OpenFile("筛选的数据.txt", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		/*fmt.Println("all:", alls)*/ /*fmt.Println("all:", alls)*/
		if ip != nil {
			n, _ := f.Seek(0, os.SEEK_END)
			// 从末尾的偏移量开始写入内容
			_, err = f.WriteAt([]byte(ip[0][1]+"\t"+time[0][1]+"\t"+video[0][0]+"\t"+article[0][0]+"\n"), n)
		}
	}
}
func main() {
	fileName := "log.log"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()
	fmt.Println("file size=", size)
	buf := bufio.NewReader(file)
	f, err := os.Create("筛选的数据.txt")
	defer f.Close()
	for {
		page := make(chan string, 1)
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		page <- line
		go findLine(page)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
	}
}
