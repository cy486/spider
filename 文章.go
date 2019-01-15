package main

import (
	"./Util"
	"fmt"
	orm2 "github.com/astaxie/beego/orm"
	_ "github.com/go-mysql/mysql"
	"regexp"
	"strings"
)

func save1(title, concent, urls string, id int) {
	orm := orm2.NewOrm()
	WC := WordContent{}
	WC.Title = string(title)
	WC.Concent = orm2.TextField(concent)
	WC.Word = int(id)
	WC.Url = string(urls)
	_, err := orm.Insert(&WC)
	if err != nil {
		fmt.Println(err)
	}
}
func small(a, b, c int) (result int) {
	result = a
	if a > b {
		result = b
	} else if a > c {
		result = c
	}
	return
}

//查询文章的方法
func findcontent(word string /*page chan int, */, id int) {
	url := "http://app.idcquan.com/?app=search&controller=index&action=search&type=all&wd=" + word
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println("读取网页错误", err)
	}
	titlerets := regexp.MustCompile(`<span class="title">(?s:(.*?))</a>`)
	title := titlerets.FindAllStringSubmatch(result, -1)

	concentret := regexp.MustCompile(`<span class="nei_rong">...(?s:(.*?))...</span>`)
	concent := concentret.FindAllStringSubmatch(result, -1)

	urlret := regexp.MustCompile(`<div class="news_nr">(?s:(.*?))" target="_blank" class="d1">`)
	alls := urlret.FindAllStringSubmatch(result, -1)
	a := small(len(title), len(concent), len(alls))
	for i := 0; i < a; i++ {
		title[i][1] = strings.Replace(title[i][1], `<span class="keyword">`, "", -1)
		title[i][1] = strings.Replace(title[i][1], `</span>`, "", -1)
		title[i][1] = strings.TrimSpace(title[i][1])
		alls[i][1] = strings.TrimSpace(alls[i][1])
		alls[i][1] = strings.Replace(alls[i][1], `<a href="`, "", -1)
		fmt.Println(alls[i][1])
		save1(title[i][1], concent[i][1], alls[i][1], id)
	}
	/*page <- id*/
}
func main() {
	//在Word表中查询数据，调用方法查询url和标题
	orm := orm2.NewOrm()
	var words HotWord
	/*page := make(chan int)*/
	for i := 1; i < 2598; i++ {
		words = HotWord{Id: i}
		_ = orm.Read(&words)
		findcontent(words.Word /*page, */, words.Id)
	}
	/*for i := 1; i < 2598; i++ {
		fmt.Printf("第%d个界面的爬取完成", <-page)
	}*/
}

/*type WordContent struct {
	Id      int `orm:"pk;auto;"`
	Title   string
	Word    int
	Concent orm2.TextField
	Url     string
}
type HotWord struct {
	Id      int `orm:"pk;auto;"`
	Word    string
	UrL     orm2.TextField
	Content orm2.TextField
	IsA     string
}

func init() {
	err := orm2.RegisterDataBase("default", "mysql", "root:123456@tcp(localhost:3306)/text?charset=UTF8")
	orm2.RegisterModel(new(WordContent), new(HotWord))
	err = orm2.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}
}*/
