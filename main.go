package main

import (
	"fmt"
	goserver "statkiGo/go-server"
)

func main() {
	// str := goserver.InitGame()
	// fmt.Printf(str)
	f, err := goserver.InitGame()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f)
	}
}
