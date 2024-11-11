package store

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkStore(b *testing.B) {
	Start()
	for i := 0; i < b.N; i++ {
		go AddItem("banana ", "Medium")
		go EditPriority(uuid.New(), "High")
		go EditTitle(uuid.New(), "ubunga")
		go ToggleComplete(uuid.New())
		go GetAllItems()
		go DeleteItem(uuid.New())
	}
}

func TestRemoveItem(t *testing.T) {
	var tests = []struct {
		name string
		list toDoList
		id   uuid.UUID
		want response
	}{
		{name: "happy path", list: createToDoListWithItems(), id: knownID, want: response{}},
		{name: "empty list", list: make(toDoList), id: knownID, want: response{err: ErrNoItem}},
		{name: "item not found", list: createToDoListWithItems(), id: uuid.New(), want: response{err: ErrNoItem}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.list.RemoveItem(tt.id)
			if !compareResonses(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
			if tt.list[tt.id] != nil {
				t.Errorf("item still exists")
			}
		})
	}
}

var knownID = uuid.New()

func createToDoListWithItems() toDoList {
	return toDoList{
		knownID:    createToDoItem("banana", HighPriority),
		uuid.New(): createToDoItem("triple", MediumPriority),
		uuid.New(): createToDoItem("frog", LowPriority),
	}
}

func compareResonses(res1, res2 response) bool {
	if res1.err != res2.err {
		fmt.Println("error fail")
		return false
	}
	if res1.item != res2.item {
		fmt.Println("item fail")
		return false
	}
	if !reflect.DeepEqual(res1.list, res2.list) {
		fmt.Println("list fail")
		return false
	}
	return true
}

// func compareToDoLists(list1, list2 toDoList) bool {
// 	for id, value := range list1 {
// 		if list2[id] != value {
// 			return false
// 		}
// 	}
// 	return true
// }
