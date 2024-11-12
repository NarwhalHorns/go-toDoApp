package store

import (
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkAddItemEmptyStore(b *testing.B) {
	store := CreateAndStartStore(nil)
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			AddItem(&store, "banana", HighPriority)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkAddItemPopulatedStore(b *testing.B) {
	startList := createToDoListWithItems(b.N, uuid.New())
	store := CreateAndStartStore(startList)
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			AddItem(&store, "banana", HighPriority)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkEditItem(b *testing.B) {
	startList := createToDoListWithItems(b.N, uuid.New())
	store := CreateAndStartStore(startList)
	idChan := make(chan uuid.UUID, b.N)
	for id := range store.toDoList {
		idChan <- id
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(3)
		itemID := <-idChan
		go func() {
			EditTitle(&store, itemID, "new Title")
			wg.Done()
		}()
		go func() {
			EditPriority(&store, itemID, MediumPriority)
			wg.Done()
		}()
		go func() {
			ToggleComplete(&store, itemID)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestDeleteItem(t *testing.T) {
	var tests = []struct {
		name string
		list toDoList
		id   uuid.UUID
		want error
	}{
		{name: "happy path", list: createToDoListWithItems(3, knownID), id: knownID, want: nil},
		{name: "empty list", list: createtoDoList(), id: knownID, want: ErrNoItem},
		{name: "item not in list", list: createToDoListWithItems(3, uuid.New()), id: uuid.New(), want: ErrNoItem},
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
		{name: "populated list", list: createToDoListWithItems(3, uuid.New())},
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

func createToDoListWithItems(amount int, firstID uuid.UUID) toDoList {
	toDoList := toDoList{
		firstID: createToDoItemPointer("banana", HighPriority),
	}
	for i := 0; i < amount-1; i++ {
		toDoList[uuid.New()] = createToDoItemPointer("banana"+strconv.Itoa(i), HighPriority)
	}
	return toDoList
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
