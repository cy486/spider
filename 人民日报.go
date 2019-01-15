package main

import (
	"./Util"
	"fmt"
	orm2 "github.com/astaxie/beego/orm"
	_ "github.com/go-mysql/mysql"
	"regexp"
	"strconv"
	"strings"
)

//爬取网页内容的函数
func SpiderContentPage(url string) (title, content string) {
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println(err)
	}
	titleret := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	all := titleret.FindAllStringSubmatch(result, 1)
	title = all[0][1]
	contentret := regexp.MustCompile(`<div class="d2txt_con clearfix">(?s:(.*?))</div>`)
	all1 := contentret.FindAllStringSubmatch(result, 1)
	all1[0][1] = strings.Replace(all1[0][1], `<p style="text-indent: 2em;">`, "", -1)
	all1[0][1] = strings.Replace(all1[0][1], `<p>`, "", -1)
	all1[0][1] = strings.Replace(all1[0][1], `</p>`, "", -1)
	content = all1[0][1]
	return
}

//保存到数据库
func SaveData(fileTitle, fileContent, fileUrl []string) {
	for i := 0; i < len(fileUrl); i++ {
		orm := orm2.NewOrm()
		crawlpage := Crawlpage{}
		crawlpage.Title = string(fileTitle[i])
		crawlpage.Content = orm2.TextField(fileContent[i])
		crawlpage.Url = string(fileUrl[i])
		_, err := orm.Insert(&crawlpage)
		if err != nil {
			fmt.Println("插入错误", err)
			continue
		}
	}
}

//爬取网页的函数
func SpiderPage(index int, page chan int) {
	url := "http://jhsjk.people.cn/result/" + strconv.Itoa(index) + "?keywords=%E4%BA%92%E8%81%94%E7%BD%91&button=%E6%90%9C%E7%B4%A2"
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println("读取网页错误", err)
	}
	urlret := regexp.MustCompile(`<li><a href="(?s:(.*?))target="_blank">`)
	alls := urlret.FindAllStringSubmatch(result, -1)
	fileTitle := make([]string, 0)
	fileContent := make([]string, 0)
	fileUrl := make([]string, 0)
	for i := 0; i < len(alls)-1; i++ {
		alls[i][1] = strings.Replace(alls[i][1], `"`, "", -1)
		alls[i][1] = alls[i][1][:16]
		/*fmt.Println(alls[i][1])*/
		url1 := "http://jhsjk.people.cn/" + alls[i][1]
		title, content := SpiderContentPage(url1)
		fileTitle = append(fileTitle, title)
		fileContent = append(fileContent, content)
		fileUrl = append(fileUrl, url1)
		/*fmt.Println(content)*/
		/*fmt.Println(alls[i][1])*/
	}
	//存入数据库
	SaveData(fileTitle, fileContent, fileUrl)
	page <- index
}

//控制流程的函数
func Towork(start, end int) {
	page := make(chan int)
	fmt.Printf("要爬取的页数是第%d页到%d页", start, end)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d个界面的爬取完成", <-page)
	}
}

//主函数
func main() {
	var start, end int
	fmt.Println("请输入要爬取的开始页")
	fmt.Scan(&start)
	fmt.Println("请输入要爬取的结束页")
	fmt.Scan(&end)
	Towork(start, end)
}

type Crawlpage struct {
	Id      int `orm:"pk;auto;"`
	Title   string
	Content orm2.TextField
	Segment string
	IsTec   string
	IsSoup  string
	IsMR    string
	IsMath  string
	IsNews  string
	Url     string
}

func init() {
	err := orm2.RegisterDataBase("default", "mysql", "root:123456@tcp(localhost:3306)/text?charset=UTF8")
	orm2.RegisterModel(new(Crawlpage))
	err = orm2.RunSyncdb("default", true, true)
	if err != nil {
		fmt.Println(err)
	}
}
