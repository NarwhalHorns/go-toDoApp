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

func createToDoItem(t string, p Priority) toDoItem {
	return toDoItem{
		title:    t,
		priority: p,
		complete: false,
	}
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

func (i *toDoItem) String() string {
	var fullString string
	fullString = "\ttitle: " + i.title + "\n"
	fullString += "\tpriority: " + string(i.priority) + "\n"
	fullString += "\tcomplete: " + i.complete.String() + "\n"
	return fullString
}

type toDoList map[uuid.UUID]toDoItem

func createtoDoList() toDoList {
	return make(toDoList)
}

func (l toDoList) AddItem(t string, p Priority) {
	l[uuid.New()] = createToDoItem(t, p)
}

func (l toDoList) RemoveItem(id uuid.UUID) {
	delete(l, id)
}

func (l toDoList) GetAllItems(res chan string) {
	var allItemsString string
	for id, item := range l {
		allItemsString += "id: " + id.String() + "\n"
		allItemsString += item.String()
	}

	res <- allItemsString
}

var ToDoList = createtoDoList()
