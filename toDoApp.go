package main

import (
	"toDoApp/cli"
)

func main() {
	mainChan := make(chan bool)

	go cli.StartCLI(mainChan)

	<-mainChan
}
