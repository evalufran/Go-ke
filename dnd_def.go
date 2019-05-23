package main

import (
	_"fmt"
	"html/template"
	"log"
	"net/http"
)

type Datas struct {
	Nome   string
	Genere string
	Razza  string
}


func main() {
	tmpl1 := template.Must(template.ParseFiles("dnd.html"))
	tmpl2 := template.Must(template.ParseFiles("answer.html"))
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl1.Execute(w, nil)
			return
		}

		
		tmpl1.Execute(w, struct{ Success bool }{true})
		
		})

		http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				tmpl2.Execute(w, nil)
				return
			}
			details := Datas{
			Nome:   r.FormValue("firstname"),
			Genere: r.FormValue("genere"),
			Razza:  r.FormValue("razza"),
		}

			// do something with details
			log.Println(details.Nome)
			log.Println(details.Genere)
			log.Println(details.Razza)
			tmpl2.Execute(w, struct{ Success bool }{true})
			})
	 http.ListenAndServe(":8080", nil)
		
}
