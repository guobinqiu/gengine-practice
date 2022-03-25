//package必须是main
package main

import (
	"fmt"
)

type cat struct{}

func (b *cat) Action() {
	fmt.Println("喵喵喵")
}

//导出
var Cat = cat{}
