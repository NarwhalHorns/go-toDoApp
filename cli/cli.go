package cli

import (
	"fmt"
	"toDoApp/store"
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

func exitCLI() {
	fmt.Println("Quitting...")
	loop = false
}

func addPrompt() {
	var t string
	for {
		t = scanInputLine("Enter title: ")
		if len(t) <= 0 {
			fmt.Println("Title must exist")
			continue
		}
		break
	}

	p := store.Priority(scanInputLine("Enter priority (High, Medium, Low) or enter for default (Medium): "))

	go store.ToDoList.AddItem(t, p)
	println("Item added to list")
}

func listItems() {
	var res = make(chan string)
	go store.ToDoList.GetAllItems(res)
	listString := <-res
	fmt.Print(listString)
}

var loop = true

func StartCLI(mainChan chan bool) {
	fmt.Println("------- ToDoApp --------")
	for loop {
		input := scanInputLine("Enter command: ")
		switch input {
		case "add":
			addPrompt()
		case "delete":
			// deletePrompt()
		case "edit":
			// editPrompt()
		case "search":
			// searchPrompt()
		case "list":
			listItems()
		case "exit":
			exitCLI()
		default:
			displayHelp()
		}
	}

	close(mainChan)
}
