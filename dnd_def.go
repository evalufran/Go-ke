package main

import (
	"html/template"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"strconv"
)
//funzione di errore
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
//Crea una struttura dal json
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
//crea la struttura per la scheda personaggio
type Personaggio struct {
	Utente       string
	Razza        string
	Genere       string
	NomePersonaggio         string
	Allineamento string
	Taglia       string
	Dio          string
	Classe       string
	Selection    []int
}
//genera scheda del personaggio randomizzata
func (p *Personaggio) Genera() error {

	selezioneNome := Conf.NomePersonaggio[p.Selection[0]][p.Selection[1]] //crea il nome del personaggio basandosi su razza e genere 
	rand.Seed(time.Now().UnixNano())
	p.NomePersonaggio =  selezioneNome[rand.Intn(len(selezioneNome))] 
	p.Allineamento = Conf.Allineamento[rand.Intn(len(Conf.Allineamento))]
	p.Taglia = Conf.Taglia[rand.Intn(len(Conf.Taglia))]
	p.Classe = Conf.Classe[rand.Intn(len(Conf.Classe))]
	p.Dio = Conf.Dio[rand.Intn(len(Conf.Dio))]

	return nil
}
//legge il json
func init() {
	checkErrors(ReadFromJSON(&Conf, "conf.json"))
}

func main() {
	//parsa i templates
	tmpl1 := template.Must(template.ParseFiles("dnd.html"))
	tmpl2 := template.Must(template.ParseFiles("answer.html"))
//funzione che gestisce le preferenze dell'utente nella pagina home
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {

		homeMap := make(map[string]interface{})
		homeMap["Razza"] = Conf.Razza
		homeMap["Genere"] = Conf.Genere
		tmpl1.Execute(w, homeMap)

	})
	//funzione che gestisce la pagina di riposta /process
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		//array di stringhe in cui vengono salvate le scelte dell'utente
		selezioni:=[]string{ r.FormValue("firstname"), r.FormValue("genere"), r.FormValue("razza")}
		
		processMap := make(map[string]interface{})
		processMap["Utente"] =selezioni[0]
		//funzione che converte in interi i valori di selezioni
		convertiGenere,_:=strconv.Atoi(selezioni[1])
		processMap["Genere"]=Conf.Genere[convertiGenere]
		//funzione che converte in interi i valori di selezioni
		convertiRazza,_:=strconv.Atoi(selezioni[2])
		processMap["Razza"]=Conf.Razza[convertiRazza]
		//crea l'oggetto personaggio
		personaggio := new(Personaggio)
		//inserisce le key dei valori scelti dall'utente dentro p.selection, per poter generare il nome adatto
		personaggio.Selection =[]int{convertiGenere, convertiRazza}
		
		personaggio.Genera()
		processMap["Nome"] = personaggio.NomePersonaggio
		processMap["Allineamento"] = personaggio.Allineamento
		processMap["Taglia"] = personaggio.Taglia
		processMap["Classe"] = personaggio.Classe
		processMap["Divinita"] = personaggio.Dio
		
		
		tmpl2.Execute(w, processMap)

	})

	
	http.ListenAndServe(":8080", nil)
}
