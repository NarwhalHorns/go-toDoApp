package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
)

const filepathPrefix = "store/perm/"

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
	Title    string
	Priority Priority
	Complete complete
}

func createToDoItem(t string, p Priority) toDoItem {
	return toDoItem{
		Title:    t,
		Priority: p,
		Complete: false,
	}
}

func (i *toDoItem) UpdateTitle(t string) response {
	i.Title = t
	return response{}
}

func (i *toDoItem) UpdatePriority(p Priority) response {
	i.Priority = p
	return response{}
}

func (i *toDoItem) ToggleComplete() response {
	i.Complete = !i.Complete
	return response{}
}

func ToDoItemToString(id uuid.UUID, i toDoItem) string {
	var fullString string
	fullString = "id: " + id.String() + "\n"
	fullString += "\ttitle: " + i.Title + "\n"
	fullString += "\tpriority: " + string(i.Priority) + "\n"
	fullString += "\tcomplete: " + i.Complete.String() + "\n"
	return fullString
}

func (i toDoItem) GetValues() (string, string, bool) {
	return i.Title, string(i.Priority), bool(i.Complete)
}

type toDoList map[uuid.UUID]*toDoItem

func createtoDoList() toDoList {
	return make(toDoList)
}

func (l toDoList) AddItem(i toDoItem) response {
	l[uuid.New()] = &i
	return response{}
}

func (l toDoList) RemoveItem(id uuid.UUID) response {
	delete(l, id)
	return response{}
}

func (l toDoList) GetAllItems() response {
	return response{
		list: l,
	}
}

func ToDoListToString(l toDoList) string {
	var fullString string
	for id, item := range l {
		fullString += ToDoItemToString(id, *item)
	}
	return fullString
}

func readListFromFile(filepath string) (toDoList, error) {
	jsonData, err := os.ReadFile(filepathPrefix + filepath)
	if err != nil {
		return createtoDoList(), err
	}
	var list toDoList
	err = json.Unmarshal(jsonData, &list)
	if err != nil {
		return createtoDoList(), err
	}
	return list, nil
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
	item      toDoItem
	itemPtr   *toDoItem
}

type Store struct {
	toDoList  toDoList
	inputChan chan input
	filepath  string
}

func CreateStore(l toDoList, filepath string) Store {
	list, err := readListFromFile(filepath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
	for id, item := range l {
		list[id] = item
	}
	return Store{
		toDoList:  list,
		inputChan: make(chan input),
		filepath:  filepath,
	}
}

func (s *Store) Start() {
	go func() {
		for input := range s.inputChan {
			switch input.operation {
			case addItem:
				input.response <- s.toDoList.AddItem(input.item)
			case deleteItem:
				input.response <- s.toDoList.RemoveItem(input.id)
			case updatePriority:
				input.response <- input.itemPtr.UpdatePriority(input.priority)
			case updateTitle:
				input.response <- input.itemPtr.UpdateTitle(input.title)
			case toggleComplete:
				input.response <- input.itemPtr.ToggleComplete()
			case listItems:
				input.response <- s.toDoList.GetAllItems()
			case search:
				// myToDoList.Search()
			}
		}
	}()
}

func (s *Store) WriteToJson() error {
	jsonData, err := json.Marshal(s.toDoList)
	if err != nil {
		return err
	}
	file, err := os.Create(filepathPrefix + s.filepath)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func CreateAndStartStore(l toDoList, filepath string) Store {
	store := CreateStore(l, filepath)
	store.Start()
	return store
}

func AddItem(s *Store, title string, priority Priority) {
	item := createToDoItem(title, priority)
	res := make(chan response)
	var input = input{
		operation: addItem,
		item:      item,
		response:  res,
	}
	s.inputChan <- input
	<-res
}

func GetAllItems(s *Store) toDoList {
	res := make(chan response)
	var input = input{
		operation: listItems,
		response:  res,
	}
	s.inputChan <- input
	response := <-res
	return response.list
}

var ErrNoItem = errors.New("no item found")

func DeleteItem(s *Store, id uuid.UUID) error {
	_, ok := s.toDoList[id]
	if !ok {
		return ErrNoItem
	}
	res := make(chan response)
	var input = input{
		operation: deleteItem,
		id:        id,
		response:  res,
	}
	s.inputChan <- input
	<-res
	return nil
}

func EditPriority(s *Store, id uuid.UUID, priority Priority) error {
	item, ok := s.toDoList[id]
	if !ok {
		return ErrNoItem
	}
	res := make(chan response)
	var input = input{
		operation: updatePriority,
		priority:  priority,
		itemPtr:   item,
		response:  res,
	}
	s.inputChan <- input
	<-res
	return nil
}

func EditTitle(s *Store, id uuid.UUID, title string) error {
	item, ok := s.toDoList[id]
	if !ok {
		return ErrNoItem
	}
	res := make(chan response)
	var input = input{
		operation: updateTitle,
		itemPtr:   item,
		title:     title,
		response:  res,
	}
	s.inputChan <- input
	<-res
	return nil
}

func ToggleComplete(s *Store, id uuid.UUID) error {
	item, ok := s.toDoList[id]
	if !ok {
		return ErrNoItem
	}
	res := make(chan response)
	var input = input{
		operation: toggleComplete,
		itemPtr:   item,
		response:  res,
	}
	s.inputChan <- input
	<-res
	return nil
}
