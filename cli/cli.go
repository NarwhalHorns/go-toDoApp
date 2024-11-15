package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"toDoApp/store"

	"github.com/google/uuid"
)

func scanInputLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.Split(input, "\n")[0]
}

func displayHelp() {
	fmt.Println("Enter one of the following commands: add, delete, edit, search, list. Or enter 'exit' to quit.")
}

func exitCLI(killChan chan os.Signal) {
	fmt.Println("Quitting...")
	killChan <- os.Interrupt
}

func inputTitle() string {
	for {
		t := scanInputLine("Enter title: ")
		if len(t) > 0 {
			return t
		}
		fmt.Println("Title must exist")
	}
}

func addPrompt() {
	t := inputTitle()

	p := store.Priority(scanInputLine("Enter priority (High, Medium, Low) or enter for default (Medium): "))

	store.AddItem(memStore, t, p)
	println("Item added to list")
}

func deletePrompt() {
	id, err := uuid.Parse(scanInputLine("Enter id: "))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = store.DeleteItem(memStore, id)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Item deleted")
	}
}

func listItems() {
	items := store.GetAllItems(memStore)
	allItemsString := store.ToDoListToString(items)
	fmt.Println(allItemsString)
}

func editPrompt() {
	id, err := uuid.Parse(scanInputLine("Enter id of item to edit: "))
	if err != nil {
		fmt.Println(err)
		return
	}
	input := scanInputLine("Enter 't' to edit title, 'p' to edit priority or 'c' to toggle the complete state: ")
	switch input {
	case "c":
		err = store.ToggleComplete(memStore, id)
		if err != nil {
			fmt.Println(err)
		}
	case "t":
		t := inputTitle()
		err = store.EditTitle(memStore, id, t)
		if err != nil {
			fmt.Println(err)
		}
	case "p":
		p := scanInputLine("Enter priority: ")
		err = store.EditPriority(memStore, id, store.Priority(p))
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Item edited")
}

// func searchPrompt() {
// 	t := scanInputLine("Enter title to search by: ")
// 	item := store.Search(t)
// 	fmt.Println(item)
// }

func commandSwitch(input string, killChan chan os.Signal, loop *bool) {
	switch input {
	case "add":
		addPrompt()
	case "delete":
		deletePrompt()
	case "edit":
		editPrompt()
	case "search":
		// searchPrompt()
	case "list":
		listItems()
	case "exit":
		exitCLI(killChan)
		*loop = false
	default:
		displayHelp()
	}
}

var memStore *store.Store

func Start(killChan chan os.Signal, store *store.Store) {
	memStore = store
	fmt.Println("------- ToDoApp --------")
	go func() {
		var loop = true
		for loop {
			input := scanInputLine("Enter command: ")
			commandSwitch(input, killChan, &loop)
		}
	}()
}
