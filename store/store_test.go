package store

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkStore(b *testing.B) {
	store := CreateAndStartStore(nil)
	for i := 0; i < b.N; i++ {
		AddItem(&store, "banana ", "Medium")
		EditPriority(&store, uuid.New(), "High")
		EditTitle(&store, uuid.New(), "ubunga")
		ToggleComplete(&store, uuid.New())
		GetAllItems(&store)
		DeleteItem(&store, uuid.New())
	}
}

func TestDeleteItem(t *testing.T) {
	var tests = []struct {
		name string
		list toDoList
		id   uuid.UUID
		want error
	}{
		{name: "happy path", list: createToDoListWithItems(), id: knownID, want: nil},
		{name: "empty list", list: createtoDoList(), id: knownID, want: ErrNoItem},
		{name: "item not in list", list: createToDoListWithItems(), id: uuid.New(), want: ErrNoItem},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := CreateAndStartStore(tt.list)
			got := DeleteItem(&store, tt.id)
			if tt.want != got {
				t.Errorf("got %v, want %v", got, tt.want)
			}
			_, ok := store.toDoList[tt.id]
			if ok {
				t.Errorf("item not deleted")
			}
		})
	}
}

func TestCreateStorePopulatesList(t *testing.T) {
	var tests = []struct {
		name string
		list toDoList
	}{
		{name: "populated list", list: createToDoListWithItems()},
		{name: "nil list", list: createtoDoList()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := CreateAndStartStore(tt.list)
			got := store.toDoList
			if !reflect.DeepEqual(tt.list, got) {
				t.Errorf("got %v, want %v", got, tt.list)
			}
		})
	}
}

var knownID = uuid.New()

func createToDoListWithItems() toDoList {
	return toDoList{
		knownID:    createToDoItemPointer("banana", HighPriority),
		uuid.New(): createToDoItemPointer("triple", MediumPriority),
		uuid.New(): createToDoItemPointer("frog", LowPriority),
	}
}

func createToDoItemPointer(t string, p Priority) *toDoItem {
	item := createToDoItem(t, p)
	return &item
}

// func getFirstToDoItem(l toDoList) toDoItem {
// 	var returnItem toDoItem
// 	for _, item := range l {
// 		returnItem = *item
// 		break
// 	}
// 	return returnItem
// }

// func compareResonses(res1, res2 response) bool {
// 	if res1.err != res2.err {
// 		fmt.Println("error fail")
// 		return false
// 	}
// 	if res1.item != res2.item {
// 		fmt.Println("item fail")
// 		return false
// 	}
// 	if !reflect.DeepEqual(res1.list, res2.list) {
// 		fmt.Println("list fail")
// 		return false
// 	}
// 	return true
// }

// func compareToDoLists(list1, list2 toDoList) bool {
// 	for id, value := range list1 {
// 		if list2[id] != value {
// 			return false
// 		}
// 	}
// 	return true
// }
