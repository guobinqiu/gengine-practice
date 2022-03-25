//package必须是main
package main

import (
	"fmt"
)

type dog struct{}

func (b *dog) Action() {
	fmt.Println("汪汪汪")
}

//导出
var Dog = dog{}
