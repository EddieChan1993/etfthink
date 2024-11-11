package main

import (
	"etfthink/core"
	"fmt"
	"github.com/EddieChan1993/gcore/utils/cast"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("etfthink path isUp(bool)")
		return
	}
	path := cast.ToString(args[1])
	isUp := cast.ToBool(args[2])
	fmt.Println(path, isUp)
	core.Run(path, isUp)
}
