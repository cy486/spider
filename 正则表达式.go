package main

import (
	"fmt"
	"regexp"
)

/*(?s:(.*?))*/
func main() {
	str := `101.200.101.207 - - [10/Nov/2016:00:01:08 +0800] "GET /static/img/common/logo.png?t=1478707267406 HTTP/1.1" 200 4468 "static.neuedu.com" "http://www.neusoft.com/article/11325" - "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b" "199.30.25.88" 10.100.15.240:80 200 0.001 0.001`
	//解析编译正则表达式
	ret := regexp.MustCompile(`"http://www.neusoft.com/article/(?s:(.*?))"`) //``表示原生字符串
	//提取需要的信息
	alls := ret.FindAllStringSubmatch(str, -1)
	fmt.Println(alls[0][0])
}
