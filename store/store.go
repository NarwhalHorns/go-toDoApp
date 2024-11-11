package store

import (
	"errors"

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

var ErrNoItem = errors.New("no item found")

func (l toDoList) RemoveItem(id uuid.UUID) response {
	_, ok := l[id]
	if !ok {
		return response{
			err: ErrNoItem,
		}
	}
	delete(l, id)
	return response{}
}

func (l toDoList) GetAllItems() response {
	return response{
		list: l,
	}
}

func (l toDoList) DoesItemExist(id uuid.UUID) bool {
	return l[id] != nil
}

func ToDoListToString(l toDoList) string {
	var fullString string
	for id, item := range l {
		fullString += ToDoItemToString(id, *item)
	}
	return fullString
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

type response struct {
	list toDoList
	item toDoItem
	err  error
}

type input struct {
	id        uuid.UUID
	operation operation
	title     string
	priority  Priority
	response  chan response
}

var myToDoList = createtoDoList()
var inputChan = make(chan input)

func Start() {
	go func() {
		for input := range inputChan {
			switch input.operation {
			case addItem:
				myToDoList.AddItem(input.title, input.priority)
				input.response <- response{}
			case deleteItem:
				input.response <- myToDoList.RemoveItem(input.id)
			case updatePriority:
				if !myToDoList.DoesItemExist(input.id) {
					input.response <- response{
						err: errors.New("item does not exist to update"),
					}
					continue
				}
				myToDoList[input.id].UpdatePriority(input.priority)
				input.response <- response{}
			case updateTitle:
				if !myToDoList.DoesItemExist(input.id) {
					input.response <- response{
						err: errors.New("item does not exist to update"),
					}
					continue
				}
				myToDoList[input.id].UpdateTitle(input.title)
				input.response <- response{}
			case toggleComplete:
				if !myToDoList.DoesItemExist(input.id) {
					input.response <- response{
						err: errors.New("item does not exist to update"),
					}
					continue
				}
				myToDoList[input.id].ToggleComplete()
				input.response <- response{}
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
	res := make(chan response)
	var input = input{
		operation: listItems,
		response:  res,
	}
	inputChan <- input
	response := <-res
	return response.list
}

func DeleteItem(id uuid.UUID) error {
	res := make(chan response)
	var input = input{
		operation: deleteItem,
		id:        id,
		response:  res,
	}
	inputChan <- input
	response := <-res
	if response.err != nil {
		return response.err
	}
	return nil
}

func EditPriority(id uuid.UUID, priority Priority) error {
	res := make(chan response)
	var input = input{
		operation: updatePriority,
		priority:  priority,
		id:        id,
		response:  res,
	}
	inputChan <- input
	response := <-res
	if response.err != nil {
		return response.err
	}
	return nil
}

func EditTitle(id uuid.UUID, title string) error {
	res := make(chan response)
	var input = input{
		operation: updateTitle,
		id:        id,
		title:     title,
		response:  res,
	}
	inputChan <- input
	response := <-res
	if response.err != nil {
		return response.err
	}
	return nil
}

func ToggleComplete(id uuid.UUID) error {
	res := make(chan response)
	var input = input{
		operation: toggleComplete,
		id:        id,
		response:  res,
	}
	inputChan <- input
	response := <-res
	if response.err != nil {
		return response.err
	}
	return nil
}
