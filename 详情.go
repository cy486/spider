package main

import (
	"./Util"
	"fmt"
	"regexp"
	"strings"
)

func main() {
	a := "大数据"
	url := "http://app.idcquan.com/?app=search&controller=index&action=search&type=all&wd=" + a
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println("读取网页错误", err)
	}
	urlret := regexp.MustCompile(`<div class="news_nr">(?s:(.*?))" target="_blank" class="d1">`)
	alls := urlret.FindAllStringSubmatch(result, -1)
	for i := 0; i < len(alls); i++ {
		alls[i][1] = strings.Replace(alls[i][1], `<a href="`, "", -1)
		alls[i][1] = strings.TrimSpace(alls[i][1])
		fmt.Println(alls[i][1])
	}
}
