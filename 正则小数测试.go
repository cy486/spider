package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := `116.231.194.12 - - [10/Nov/2016:00:01:13 +0800] "POST /course/collectvideo HTTP/1.1" 200 72 "www.neusoft.com" "http://www.neusoft.com/video/2469" renderingMode=0&bufferTime=0&videoFileName=http%3A%2F%2Fv2.neuedu.com%2F86c47cd6-ba30-43ee-a6f4-5b0533294a42%2FM.mp4%3Fauth_key%3D1478707756-0-0-1b94c64fbd08fd2ef78ea5cc3bd88249&videoId=2469&errorMsg=&currentHd=1&oldHd=NaN&type=4&fullscreen=1&cdn=v2.neuedu.com&source=1&currentSpeed=1.0+X&oldSpeed=&winWidth=1366&winHeight=768 "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36" "-" 10.100.136.65:80 200 0.163 0.163`
	ipret := regexp.MustCompile(`([0-9]{1,3}\.){3}[0-9]{1,3}`)
	timeret := regexp.MustCompile(`\[(.+)\]`)
	articleret := regexp.MustCompile(`http://www.neusoft.com/article/[0-9]{1,4}`)
	videoret := regexp.MustCompile(`http://www.neusoft.com/video/[0-9]{1,4}`)
	ip := ipret.FindAllStringSubmatch(str, -1)
	time := timeret.FindAllStringSubmatch(str, -1)
	article := articleret.FindAllStringSubmatch(str, -1)
	video := videoret.FindAllStringSubmatch(str, -1)
	/*fmt.Println("all:", alls)*/ /*fmt.Println("all:", alls)*/
	for _, one := range ip {
		fmt.Printf("ip:%s\n", one[0])
	}
	for _, one := range time {
		fmt.Printf("时间：%s\n", one[0])
	}
	for _, one := range article {
		fmt.Printf("纹章地址：%s\n", one[0])
	}
	for _, one := range video {
		fmt.Printf("视频地址：%s\n", one[0])
	}
}
