package main

import "fmt"

func main() {
	LogInit(false)
	err := runCmd()
	if err != nil {
		fmt.Println(err.Error())
	}
}
