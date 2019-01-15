package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/**
1.输入要爬取的页码
2.调用爬取函数并进行流程控制 ToWork函数
3.使用spider类，进行页面的爬取，可以进行筛选和保存文件的操作
4.进行网页的数据读取HttpGet2，期间进行对网页的二次读取。
5.保存文件操作savepage
*/
//获取一个网页的所有内容
func HttpGet2(url string) (result string, err error) {
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

func SaveJokePage(idex int, title, content []string) {
	f, err := os.Create("第" + strconv.Itoa(idex) + "页段子.txt")
	if err != nil {
		fmt.Println("打开文件失败")
		return
	}
	defer f.Close()
	for i := 0; i < len(title); i++ {
		_, _ = f.WriteString(title[i] + "\n" + content[i] + "\n")
		_, _ = f.WriteString("--------------------------------" + "\n")
	}
	/*_, _ = f.WriteString("标题" + title + "/n" + "内容：" + content)*/
}

//抓取一个网页，带有10个段子 --- 10url
func SpiderPage1(idex int, page chan int) {
	url := "http://www.pengfu.com/xiaohua_" + strconv.Itoa(idex) + ".html"
	//封装函数获取段子的url
	result, err := HttpGet2(url)
	if err != nil {
		fmt.Println("Httpget err", err)
		return
	}
	//解析编译正则
	ret := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`)
	alls := ret.FindAllStringSubmatch(result, -1)

	//创建用于存储的切片，初始容量为0
	fileTitle := make([]string, 0)
	fileContent := make([]string, 0)
	for _, jokeURL := range alls {
		//fmt.Println(jokeURL[1])
		title, content, err := SpiderJokePage(jokeURL[1])
		if err != nil {
			fmt.Println("jokespidererr:", err)
			continue
		}
		/*fmt.Println(title)
		fmt.Println(content)*/
		fileTitle = append(fileTitle, title)
		fileContent = append(fileContent, content)
	}
	SaveJokePage(idex, fileTitle, fileContent)
	//防止主程序提前完成
	page <- idex
}

//爬取标题和内容
func SpiderJokePage(url string) (title, content string, err error) {
	result, err1 := HttpGet2(url)
	if err != nil {
		err = err1
		return
	}
	//编译解析正则表达式
	ret := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	alls := ret.FindAllStringSubmatch(result, 1) //有两处提取第一个
	for _, tmptitle := range alls {
		title = tmptitle[1]
		title = strings.Replace(title, " ", "", -1)
	}
	ret1 := regexp.MustCompile(`<div class="content-txt pt10">(?s:(.*?))<a id="prev"`)
	alls1 := ret1.FindAllStringSubmatch(result, -1) //有两处提取第一个
	for _, tmpcontent := range alls1 {
		content = tmpcontent[1]
		content = strings.Replace(content, "\n", "", -1)
		content = strings.Replace(content, "\t", "", -1)
		content = strings.Replace(content, "&nbsp;", " ", -1)
	}
	return
}

//爬取的程序
func ToWork1(start, end int) {
	fmt.Printf("要爬取的是从第%d到%d页", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage1(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d页爬取完成", <-page)
	}
}

//爬取段子，可以点击页面的链接
func main() {
	var start, end int
	fmt.Println("请输入要爬取的开始页面：")
	fmt.Scan(&start)
	fmt.Println("请输入要爬取的结束的界面：")
	fmt.Scan(&end)
	ToWork1(start, end)
}
