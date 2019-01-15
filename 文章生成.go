package main

import (
	"./Util"
	"fmt"
	_ "github.com/go-mysql/mysql"
	"os"
	"regexp"
)

func save2(data, detail []string) {
	f, err := os.Create("a.txt")
	if err != nil {
		fmt.Println("打开文件失败")
		return
	}
	defer f.Close()
	for i := 0; i < len(data); i++ {
		_, _ = f.WriteString(data[i] + "\n" + detail[i] + "\n")
		_, _ = f.WriteString("--------------------------------" + "\n")
	}
}
func con1(word string) string {
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
func main() {
	url := "http://app.idcquan.com/tags.php"
	result, err := Util.HttpGet(url)
	if err != nil {
		fmt.Println("读取网页错误", err)
	}
	urlret1 := regexp.MustCompile(`容">(?s:(.*?))</a></li>`)
	alls2 := urlret1.FindAllStringSubmatch(result, -1)
	title := make([]string, 0)
	content := make([]string, 0)
	for i := 0; i < 100; i++ {
		title = append(title, alls2[i][1])
		content = append(content, con1(alls2[i][1]))
	}
	save2(title, content)

}
