package webAPI

import (
	"net/http"
	"text/template"
	"toDoApp/store"

	"github.com/google/uuid"
)

type displayItem struct {
	Id       string
	Title    string
	Priority string
	Complete bool
}

func mainPage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "nope", 404)
	} else {
		displayPage(w)
	}
}

func displayPage(w http.ResponseWriter) {
	items := store.GetAllItems(memStore)
	var displayItems []displayItem
	for id, item := range items {
		t, p, c := item.GetValues()
		displayItems = append(displayItems, displayItem{Id: id.String(), Title: t, Priority: p, Complete: c})
	}

	tmpl, err := template.New("toDoList.go.tmpl").ParseFiles("webAPI/toDoList.go.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, displayItems)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func create(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "nope", 404)
		return
	}
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	store.AddItem(memStore, req.PostFormValue("title"), store.Priority(req.PostFormValue("priority")))
	displayPage(w)
}

func delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "nope", 404)
		return
	}
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var id uuid.UUID
	id, err = uuid.Parse(req.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	store.DeleteItem(memStore, id)
	displayPage(w)
}

var memStore *store.Store
var tmplFile = "webAPI/toDoList.go.tmpl"
var tmpl *template.Template

func Start(s *store.Store) {
	memStore = s

	// var err error
	// tmpl, err = template.New(tmplFile).ParseFiles(tmplFile)
	// if err != nil {
	// 	panic(err)
	// }

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/create", create)
	// http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	go http.ListenAndServe(":8080", nil)
}