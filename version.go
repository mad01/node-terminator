package main

import (
	"fmt"
)

var (
	Version = "not set"
)

func getVersion() {
	fmt.Printf("Version: %v\n", Version)
}
