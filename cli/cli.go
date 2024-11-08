package cli

import (
	"fmt"
	"os"
	"toDoApp/store"

	"github.com/google/uuid"
)

func scanInputLine(prompt string) string {
	fmt.Print(prompt)
	var input string
	var char rune
	var err error
	for err = nil; err == nil; {
		_, err = fmt.Scanf("%c", &char)
		if char == '\n' {
			break
		}
		input += string(char)
	}
	return input
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

	store.AddItem(t, p)
	println("Item added to list")
}

func deletePrompt() {
	id, err := uuid.Parse(scanInputLine("Enter id: "))
	if err != nil {
		fmt.Println(err)
		return
	}
	deleted := store.DeleteItem(id)
	if deleted {
		fmt.Println("Item deleted")
	} else {
		fmt.Println("Item not found")
	}
}

func listItems() {
	items := store.GetAllItems()
	allItemsString := store.ToDoListToString(items)
	fmt.Println(allItemsString)
}

func editPrompt() {
	id, err := uuid.Parse(scanInputLine("Enter id of item to edit: "))
	if err != nil {
		fmt.Println(err)
		return
	}
	input := scanInputLine("Enter 't' to edit title, 'p' to edit priority or 'c' to toggle the compelte state.")
	switch input {
	case "c":
		store.ToggleComplete(id)
	case "t":
		t := inputTitle()
		store.EditTitle(id, t)
	case "p":
		p := scanInputLine("Enter priority: ")
		store.EditPriority(id, store.Priority(p))
	}
	fmt.Println("Item edited")
}

// func searchPrompt() {
// 	t := scanInputLine("Enter title to search by: ")
// 	item := store.Search(t)
// 	fmt.Println(item)
// }

func Start(killChan chan os.Signal) {
	fmt.Println("------- ToDoApp --------")
	go func() {
		var loop = true
		for loop {
			input := scanInputLine("Enter command: ")
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
				loop = false
			default:
				displayHelp()
			}
		}
	}()
}
