package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
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
	Utente          string
	Razza           string
	Genere          string
	NomePersonaggio string
	Allineamento    string
	Taglia          string
	Dio             string
	Classe          string
	Selection       []int
}

//genera scheda del personaggio randomizzata
func (p *Personaggio) Genera() error {

	selezioneNome := Conf.NomePersonaggio[p.Selection[1]][p.Selection[0]] //crea il nome del personaggio basandosi su razza e genere
	rand.Seed(time.Now().UnixNano())
	p.NomePersonaggio = selezioneNome[rand.Intn(len(selezioneNome))]
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl1 := template.Must(template.ParseFiles("dnd.html"))
	homeMap := make(map[string]interface{})
	homeMap["Razza"] = Conf.Razza
	homeMap["Genere"] = Conf.Genere
	tmpl1.Execute(w, homeMap)

}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl2 := template.Must(template.ParseFiles("answer.html"))

	//array di stringhe in cui vengono salvate le scelte dell'utente
	selezioni := []string{r.FormValue("firstname"), r.FormValue("genere"), r.FormValue("razza")}

	processMap := make(map[string]interface{}) //mappa per salvare i parametri
	processMap["Utente"] = selezioni[0]
	convertiGenere, _ := strconv.Atoi(selezioni[1]) //funzione che converte in interi i valori di selezioni
	processMap["Genere"] = Conf.Genere[convertiGenere]
	convertiRazza, _ := strconv.Atoi(selezioni[2]) //funzione che converte in interi i valori di selezioni
	processMap["Razza"] = Conf.Razza[convertiRazza]
	personaggio := new(Personaggio) //crea l'oggetto personaggio
	//inserisce le key dei valori scelti dall'utente dentro p.selection, per poter generare il nome adatto
	personaggio.Selection = []int{convertiGenere, convertiRazza}

	personaggio.Genera()
	processMap["Nome"] = personaggio.NomePersonaggio
	processMap["Allineamento"] = personaggio.Allineamento
	processMap["Taglia"] = personaggio.Taglia
	processMap["Classe"] = personaggio.Classe
	processMap["Divinita"] = personaggio.Dio

	tmpl2.Execute(w, processMap)

	pdf := gofpdf.New("P", "mm", "A4", "") //crea il pdf
	pdf.AddPage()                          //crea la pagina
	pdf.SetFont("Arial", "B", 12)          //imposta il font
	pdf.Cell(40, 10, "Giocatore")          //crea la Nome Giocatore
	pdf.Cell(40, 10, "Personaggio")        //crea la Nome Personaggio
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, r.FormValue("firstname"))    //la seconda scritta è consecutiva alla prima
	pdf.Cell(40, 10, personaggio.NomePersonaggio) //la seconda scritta è consecutiva alla prima
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 12) //imposta il font
	pdf.Cell(40, 10, "Razza")     //crea la scritta
	pdf.Cell(40, 10, "Genere")    //crea la scritta
	pdf.Ln(8)                     //a capo (spaziatura normale)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, Conf.Razza[convertiRazza])   //la seconda scritta è consecutiva alla prima
	pdf.Cell(40, 10, Conf.Genere[convertiGenere]) //la seconda scritta è consecutiva alla prima
	pdf.Ln(8)                                     //a capo (spaziatura normale)
	pdf.SetFont("Arial", "B", 12)                 //imposta il font
	pdf.Cell(40, 10, "Allineamento")              //crea la scritta
	pdf.Cell(40, 10, "Taglia")                    //crea la scritta
	pdf.Cell(40, 10, "Classe")                    //crea la scritta
	pdf.Cell(40, 10, "Divinita'")                 //crea la scritta
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, personaggio.Allineamento)           //la seconda scritta è consecutiva alla prima
	pdf.Cell(40, 10, personaggio.Taglia)                 //la seconda scritta è consecutiva alla prima
	pdf.Cell(40, 10, personaggio.Classe)                 //la seconda scritta è consecutiva alla prima
	pdf.Cell(40, 10, personaggio.Dio)                    //la seconda scritta è consecutiva alla prima
	pdf.OutputFileAndClose("LaMiaSchedaPersonaggio.pdf") //salva il pdf

}

func main() {

	http.HandleFunc("/home", homeHandler)      //handler della pagina home
	http.HandleFunc("/process", answerHandler) //handler della pagina di risposta /process

	log.Fatal(http.ListenAndServe(":8080", nil)) //hosting pagina
}
