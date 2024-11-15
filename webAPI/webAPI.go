package webAPI

import (
	"net/http"
	"text/template"
	"toDoApp/store"

	"github.com/google/uuid"
)

func mainPage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "nope", 404)
		return
	}
	items := store.GetAllItems(memStore)
	err := tmpl.Execute(w, items)
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
	http.Redirect(w, req, "/", http.StatusSeeOther)
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
	err = store.DeleteItem(memStore, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func updateTitle(w http.ResponseWriter, req *http.Request) {
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
	title := req.FormValue("title")

	err = store.EditTitle(memStore, id, title)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func updatePriority(w http.ResponseWriter, req *http.Request) {
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
	priority := req.FormValue("priority")

	err = store.EditPriority(memStore, id, store.Priority(priority))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func updateComplete(w http.ResponseWriter, req *http.Request) {
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

	err = store.ToggleComplete(memStore, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

var memStore *store.Store
var tmpl *template.Template

func Start(s *store.Store) {
	memStore = s

	var err error
	tmpl, err = template.New("toDoList.go.tmpl").ParseFiles("webAPI/toDoList.go.tmpl")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/create", create)
	http.HandleFunc("/update/title", updateTitle)
	http.HandleFunc("/update/priority", updatePriority)
	http.HandleFunc("/update/complete", updateComplete)
	http.HandleFunc("/delete", delete)
	go http.ListenAndServe(":8080", nil)
}
