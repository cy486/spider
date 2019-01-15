package main

import (
	"./Util"
	"fmt"
	orm2 "github.com/astaxie/beego/orm"
	_ "github.com/go-mysql/mysql"
	"regexp"
	"strings"
)

func save(data, detail, urls string) {
	orm := orm2.NewOrm()
	crawlpage := Hot{}
	crawlpage.Word = string(data)
	crawlpage.Content = orm2.TextField((detail))
	crawlpage.UrL = orm2.TextField(urls)
	_, err := orm.Insert(&crawlpage)
	if err != nil {
		fmt.Println("插入错误", err)
		return
	}
}
func con(word string) string {
	url := "https://baike.baidu.com/item/" + word
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println("读取网页错误", err)
	}
	urlret := regexp.MustCompile(`<meta name="description" content="(?s:(.*?))...">`)
	alls := urlret.FindAllStringSubmatch(result, 1)
	if alls != nil {
		return alls[0][1]
	}
	return "a"
}
func findurl(word string) (urlsm string) {
	url := "http://app.idcquan.com/?app=search&controller=index&action=search&type=all&wd=" + word
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println("读取网页错误", err)
	}
	urlret := regexp.MustCompile(`<div class="news_nr">(?s:(.*?))" target="_blank" class="d1">`)
	alls := urlret.FindAllStringSubmatch(result, 2)
	if alls == nil {
		return
	} else {

		alls[0][1] = strings.Replace(alls[0][1], `<a href="`, "", -1)
		urlsm = alls[0][1]

	}
	return
}
func dowork(i string, id int, page chan int) {
	detail := con(i)
	urls := findurl(i)
	save(i, detail, urls)
	page <- id
}
func main() {
	url := "http://app.idcquan.com/tags.php"
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println("读取网页错误", err)
	}
	urlret1 := regexp.MustCompile(`容">(?s:(.*?))</a></li>`)
	alls2 := urlret1.FindAllStringSubmatch(result, -1)
	page := make(chan int)
	for i := 0; i < len(alls2); i++ {
		go dowork(alls2[i][1], i, page)
	}
	for i := 0; i < len(alls2); i++ {
		fmt.Printf("第%d个界面的爬取完成", <-page)
	}
}

func init() {
	err := orm2.RegisterDataBase("default", "mysql", "root:123456@tcp(localhost:3306)/text?charset=UTF8")
	orm2.RegisterModel(new(HotWord))
	err = orm2.RunSyncdb("default", true, true)
	if err != nil {
		fmt.Println(err)
	}
}
