package main

import "fmt"

// TODO: implement to handle number of concurrent terminations
// TODO: implement to take wait time.Duration before one node is considered done
//		 and the worker proceds to next node

func main() {
	LogInit(false)
	err := runCmd()
	if err != nil {
		fmt.Println(err.Error())
	}
}
