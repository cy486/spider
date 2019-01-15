package main

import (
	"fmt"
	"github.com/astaxie/beego"
	orm2 "github.com/astaxie/beego/orm"
	_ "github.com/go-mysql/mysql"
	"os"
)

func save3(a []HotWord) {
	f, err := os.Create("b.txt")
	if err != nil {
		fmt.Println("打开文件失败")
		return
	}
	defer f.Close()
	for i := 0; i < len(a); i++ {
		_, _ = f.WriteString(a[i].Word + "\n" + string(a[i].Content) + "\n")
		_, _ = f.WriteString("--------------------------------" + "\n")
	}
}
func main() {
	orm := orm2.NewOrm()
	var hot []HotWord
	_, err := orm.QueryTable("HotWord").All(&hot)
	if err != nil {
		beego.Info("查询失败")
		return
	}
	save3(hot)
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
	orm2.RegisterModel(new(HotWord))
	err = orm2.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}
}
