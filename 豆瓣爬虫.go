package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func main() {
	var start, end int
	fmt.Println("请输入要爬取的开始页面：")
	fmt.Scan(&start)
	fmt.Println("请输入要爬取的结束的界面：")
	fmt.Scan(&end)
	ToWork(start, end)
}

//爬取的程序
func ToWork(start, end int) {
	fmt.Printf("要爬取的是从第%d到%d页", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go DoMain(i, page)
	}
	for i := start; i <= end; i++ {
		<-page
	}
}
func DoMain(index int, page chan int) {
	//获取url
	url := "https://movie.douban.com/top250?start=" + strconv.Itoa((index-1)*25) + "&filter="
	result, err := HttpGet1(url)
	if err != nil {
		fmt.Printf("url加载错误%s", err)
		return
	}
	namereg := regexp.MustCompile(`<img width="100" alt="(?s:(.*?))"`)
	editorreg := regexp.MustCompile(`导演: (?s:(.*?))&nbsp;&nbsp;`)
	scorereg := regexp.MustCompile(` <span class="rating_num" property="v:average">(?s:(.*?))</span>`)
	name := namereg.FindAllStringSubmatch(result, -1)
	editor := editorreg.FindAllStringSubmatch(result, -1)
	score := scorereg.FindAllStringSubmatch(result, -1)
	SaveFile(index, name, editor, score)
	//写入chan防止程序提前结束
	page <- index
}
func HttpGet1(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1 //峰会钻杆数内部传递给调用者
	}
	defer resp.Body.Close()
	//循环读取网页数据1
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("读取网页完成")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[:n])
	}
	return
}
func SaveFile(index int, name, editor, score [][]string) {
	f, err := os.Create("第" + strconv.Itoa(index) + "页.txt")
	if err != nil {
		fmt.Println("文件打开错误", err)
		return
	}
	defer f.Close()
	n := len(editor)
	_, _ = f.WriteString("电影名称" + "\t\t\t\t\t" + "导演" + "\t\t\t\t\t\t" + "电影评分" + "\t\t\t\t\t\t\t" + "\n")
	for i := 0; i < n; i++ {
		_, _ = f.WriteString(name[i][1] + "\t\t\t" + editor[i][1] + "\t\t\t" + score[i][1] /*+ "\t\t\t" + year[i][1]*/ + "\n")
	}
}
