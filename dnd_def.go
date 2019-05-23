package main
//importazione pacchetti
import (
	"html/template"
	"net/http"
	"log"
	_"fmt"
)
//creazione di una struttura
type Datas struct {
	Name   string
	Class	[]string
	Razza	[]string
}

func main() {
	//lettura del file HTML, permette di riconoscere il template che deve essere abilitato nella porta :8080
	tmpl := template.Must(template.ParseFiles("dnd.html"))
 // funzione che gestisce i dati, quando attivata invia i dati inseriti nel form alla struttura 
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
//se la funzione ha inviato corrrettamente i dati ritorna true dovrebbe fare il redirect su altra pagina html
		tmpl.Execute(w, struct{ Success bool }{true})
	})
//abilitazione della porta 8080
	http.ListenAndServe(":8080", nil)
}