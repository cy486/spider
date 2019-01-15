package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var start, end int
	fmt.Println("请输入要爬取的起始页")
	fmt.Scan(&start)
	fmt.Println("其输入爬取的终止页")
	fmt.Scan(&end)

	working(start, end)
}

//爬取一个页面的程序
func SpiderPage(i int, page chan int) {
	url := "https://tieba.baidu.com/f?kw=%E7%BB%BF%E5%AE%9D%E7%9F%B3493&ie=utf-8&pn=" + strconv.Itoa((i-1)*50)
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*fmt.Println(result)*/
	f, err := os.Create("第" + strconv.Itoa(i) + "页" + ".html")
	if err != nil {
		fmt.Println("HttpGet err", err)
		return
	}
	_, _ = f.WriteString(result)
	_ = f.Close() //保存好文件后关闭
	page <- i
}

//爬取代码：并发调用
func working(start, end int) {
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d个界面的爬取完成", <-page)
	}
}

//获取页面信息并封装送回
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1 //峰会钻杆数内部传递给调用者
	}
	defer resp.Body.Close()
	//循环读取网页数据
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
