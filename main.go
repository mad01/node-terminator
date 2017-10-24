package main

func main() {
	LogInit(false)
	err := runCmd()
	if err != nil {
		panic(err)
	}
}
