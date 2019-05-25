package main

import (
	_ "fmt"
	_ "github.com/balacode/one-file-pdf"
	"html/template"

	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func checkErrors(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

//ReadFromJSON function load a json file into a struct or return error
func ReadFromJSON(t interface{}, filename string) error {

	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonFile), t)
	if err != nil {
		return err
	}

	return nil
}

type Datas struct {
	Classe          []string     `json:"classe"`
	Genere          []string     `json:"genere"`
	Razza           []string     `json:"razza"`
	Allineamento    []string     `json:"allineamento"`
	Taglia          []string     `json:"taglia"`
	Dio             []string     `json:"dio"`
	NomePersonaggio [][][]string `json:"nomepersonaggio"`
}

var Conf Datas

type Personaggio struct {
	Utente       string
	Razza        string
	Genere       string
	Nome         string
	Allineamento string
	Taglia       string
	Dio          string
	Classe       string
	Selection    []int
}

func (p *Personaggio) Genera() error {

	selezioneNome := Conf.NomePersonaggio[p.Selection[0]][p.Selection[1]]
	rand.Seed(time.Now().UnixNano())
	p.Nome = selezioneNome[rand.Intn(len(selezioneNome))]
	p.Allineamento = Conf.Allineamento[rand.Intn(len(Conf.Allineamento))]
	p.Taglia = Conf.Taglia[rand.Intn(len(Conf.Taglia))]
	p.Classe = Conf.Classe[rand.Intn(len(Conf.Classe))]
	p.Dio = Conf.Dio[rand.Intn(len(Conf.Dio))]

	return nil
}

func init() {
	checkErrors(ReadFromJSON(&Conf, "conf.json"))
}

func main() {
	tmpl1 := template.Must(template.ParseFiles("dnd.html"))
	tmpl2 := template.Must(template.ParseFiles("answer.html"))

	/*template.FuncMap{
		"Iterate": func(count *uint) []uint {
			var i uint
			var Items []uint
			for i=0; i<(*count); i++{
				Items = append(Items,i)
			}
			return Items
		},
	}*/

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {

		homeMap := make(map[string]interface{})
		homeMap["Razza"] = Conf.Razza
		homeMap["Genere"] = Conf.Genere
		tmpl1.Execute(w, homeMap)

		log.Println(homeMap["Razza"])
		log.Println(Conf.Razza)

	})
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {

		// 	//selezioni := []string{r.FormValue("firstname"), r.FormValue("genere"), r.FormValue("razza")}
		// 	personaggioUtente := new(Personaggio)

		// 	checkErrors(personaggioUtente.Genera())
		// 	personaggioUtente.Nome = r.FormValue("Nome")
		// 	personaggioUtente.Genere = r.FormValue("Genere")
		// 	personaggioUtente.Selection
		// 	personaggioUtente.Razza = r.FormValue("Razza")

		// 	//l'indice dell'array delle razze all'elemento con valore = personaggioUtente.Razza
		// 	//personaggioUtente.Allineamento = personaggioUtente.Allineamento
		// 	scelte := make(map[string]interface{})
		// 	scelte["Nome"] = personaggioUtente.Nome
		// 	scelte["Genere"]= personaggioUtente.Genere
		// 	scelte["Razza"] = personaggioUtente.Razza
		// 	//scelte["all"] = personaggioUtente.Allineamento
		//s	scelte["all"]=all
		tmpl2.Execute(w, struct{ Success bool }{true})

		//log.Println(scelte)

	})

	http.ListenAndServe(":8080", nil)
}
