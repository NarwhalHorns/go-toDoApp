package main

import (
	"os"
	"os/signal"
	"toDoApp/cli"
	"toDoApp/store"
)

func main() {
	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Interrupt)

	store.Start()
	cli.Start(killChan)

	<-killChan
}
