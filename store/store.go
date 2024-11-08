package store

import (
	"github.com/google/uuid"
)

type Priority string

const (
	MediumPriority Priority = "Medium"
	HighPriority   Priority = "High"
	LowPriority    Priority = "Low"
)

type complete bool

func (c complete) String() string {
	if c {
		return "true"
	}
	return "false"
}

type toDoItem struct {
	title    string
	priority Priority
	complete complete
}

func createToDoItem(t string, p Priority) *toDoItem {
	itemPtr := &toDoItem{
		title:    t,
		priority: p,
		complete: false,
	}
	return itemPtr
}

func (i *toDoItem) UpdateTitle(t string) {
	i.title = t
}

func (i *toDoItem) UpdatePriority(p Priority) {
	i.priority = p
}

func (i *toDoItem) ToggleComplete() {
	i.complete = !i.complete
}

func ToDoItemToString(id uuid.UUID, i toDoItem) string {
	var fullString string
	fullString = "id: " + id.String() + "\n"
	fullString += "\ttitle: " + i.title + "\n"
	fullString += "\tpriority: " + string(i.priority) + "\n"
	fullString += "\tcomplete: " + i.complete.String() + "\n"
	return fullString
}

type toDoList map[uuid.UUID]*toDoItem

func createtoDoList() toDoList {
	return make(toDoList)
}

func (l toDoList) AddItem(t string, p Priority) {
	l[uuid.New()] = createToDoItem(t, p)
}

func (l toDoList) RemoveItem(id uuid.UUID) bool {
	_, ok := l[id]
	if !ok {
		return false
	}
	delete(l, id)
	return true
}

func (l toDoList) GetAllItems() toDoList {
	return l
}

func ToDoListToString(l toDoList) string {
	res := make(chan string)
	go func() {
		var fullString string
		for id, item := range l {
			fullString += ToDoItemToString(id, *item)
		}
		res <- fullString
	}()
	return <-res
}

type operation int

const (
	addItem operation = iota
	deleteItem
	updatePriority
	updateTitle
	toggleComplete
	listItems
	search
)

type input struct {
	id        uuid.UUID
	operation operation
	title     string
	priority  Priority
	response  chan any
}

var myToDoList = createtoDoList()
var inputChan = make(chan input)

func Start() {
	go func() {
		for input := range inputChan {
			switch input.operation {
			case addItem:
				myToDoList.AddItem(input.title, input.priority)
			case deleteItem:
				myToDoList.RemoveItem(input.id)
			case updatePriority:
				myToDoList[input.id].UpdatePriority(input.priority)
			case updateTitle:
				myToDoList[input.id].UpdateTitle(input.title)
			case toggleComplete:
				myToDoList[input.id].ToggleComplete()
			case listItems:
				input.response <- myToDoList.GetAllItems()
			case search:
				// myToDoList.Search()
			}
		}
	}()
}

func AddItem(title string, priority Priority) {
	var input = input{
		operation: addItem,
		title:     title,
		priority:  priority,
	}
	inputChan <- input
}

func GetAllItems() toDoList {
	res := make(chan toDoList)
	var input = input{
		operation: listItems,
		response:  res,
	}
	inputChan <- input
	return <-res
}

func DeleteItem(id uuid.UUID) bool {
	res := make(chan bool)
	var input = input{
		operation: deleteItem,
		id:        id,
		response:  res,
	}
	inputChan <- input
	return <-res
}

func EditPriority(id uuid.UUID, priority Priority) bool {
	res := make(chan bool)
	var input = input{
		operation: updatePriority,
		priority:  priority,
		id:        id,
		response:  res,
	}
	inputChan <- input
	return <-res
}

func EditTitle(id uuid.UUID, title string) bool {
	res := make(chan bool)
	var input = input{
		operation: updateTitle,
		id:        id,
		title:     title,
		response:  res,
	}
	inputChan <- input
	return <-res
}

func ToggleComplete(id uuid.UUID) bool {
	res := make(chan bool)
	var input = input{
		operation: toggleComplete,
		id:        id,
		response:  res,
	}
	inputChan <- input
	return <-res
}
