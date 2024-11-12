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

	store := store.CreateAndStartStore(nil)
	cli.Start(killChan, &store)

	<-killChan
}
