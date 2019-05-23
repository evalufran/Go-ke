package main

import (
	"html/template"
	"net/http"
	"log"
	_"fmt"
)

type Datas struct {
	Name   string
	Class	[]string
	Razza	[]string
}

func main() {
	tmpl := template.Must(template.ParseFiles("dnd.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := Datas{
			Name:   r.FormValue("firstname"),
			Class:	r.Form["classe"],
			Razza:	r.Form["razza"],
		}

		
		// do something with details
		log.Println(details.Name)

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":8080", nil)
}