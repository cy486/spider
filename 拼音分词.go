package main

import (
	"fmt"
	"github.com/Lofanmi/pinyin-golang/pinyin"
	orm2 "github.com/astaxie/beego/orm"
	"strings"
)

func shouzimu(s string) (is_a string) {
	dict := pinyin.NewDict()
	pin := dict.Sentence(s).Unicode()
	is_a = strings.ToUpper(string(pin[0]))
	return
}
func Cha(i int) {
	orm := orm2.NewOrm()
	hotword := HotWord{Id: i}
	err := orm.Read(&hotword)
	if err != nil {
		fmt.Println("查找失败")
		return
	}

	word := hotword.Word
	hotword.IsA = shouzimu(string(word))
	_, err = orm.Update(&hotword)
	if err != nil {
		fmt.Println("更新失败")
	}
}
func main() {
	for i := 1; i < 2598; i++ {
		Cha(i)
	}
	dict := pinyin.NewDict()
	str := `Amer`
	pin := dict.Sentence(str).Unicode()
	fmt.Println(strings.ToUpper(string(pin[0])))
}
