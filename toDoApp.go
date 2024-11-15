package main

import (
	"fmt"
	"os"
	"os/signal"
	"toDoApp/cli"
	"toDoApp/store"
	"toDoApp/webAPI"
)

func main() {
	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Interrupt)

	store := store.CreateAndStartStore(nil, "toDoList.json")
	cli.Start(killChan, &store)
	webAPI.Start(&store)

	<-killChan
	err := store.WriteToJson()
	if err != nil {
		fmt.Println("failed to write list to json. error: ", err)
	}
}
