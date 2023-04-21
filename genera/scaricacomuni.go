// scarica l'elenco dei comuni dall'istat, crea un file comuni.go con l'elenco dei comuni, il codice istat ed altre informazioni

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/encoding/charmap"
)

const (
	comuniURL = "https://www.istat.it/storage/codici-unita-amministrative/Elenco-comuni-italiani.csv"
)

//Comunecodice struttura per memorizzare le informazioni estratte
type Comunecodice struct {
	Codice       string
	Comune       string
	Provincia    string
	Targa        string
	Regione      string
	Incittametro bool
	CoIdx        string
}

//ByCoIdx implementa interface per riordinare l'elenco per comune-indice (sort.Sort...)
type ByCoIdx []Comunecodice

func (a ByCoIdx) Len() int           { return len(a) }
func (a ByCoIdx) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCoIdx) Less(i, j int) bool { return strings.Compare(a[i].CoIdx, a[j].CoIdx) <= 0 }

//Normalizza : esegue alcune operazioni per permettere di confrontare i nomi in maniera agnostica dalle vocali
func Normalizza(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile("è|é").ReplaceAllString(s, "e")
	s = regexp.MustCompile("à").ReplaceAllString(s, "a")
	s = regexp.MustCompile("ù").ReplaceAllString(s, "u")
	s = regexp.MustCompile("ò").ReplaceAllString(s, "o")
	s = regexp.MustCompile("ì").ReplaceAllString(s, "i")
	return regexp.MustCompile("[^a-z]").ReplaceAllString(s, "")
}

func main() {
	var s, c, prv string
	var cm bool
	cc := make([]Comunecodice, 0, 8000)

	response, err := http.Get(comuniURL)
	if err != nil {
		log.Fatal("Errore", err)
	}
	defer response.Body.Close()
	//reader per la decodifica da Windows 1252/ISO8859_1 a UTF-8
	rv := charmap.ISO8859_1.NewDecoder().Reader(response.Body)
	//legge dal csv
	r := csv.NewReader(rv)
	r.Comma = ';'
	//evita la prima linea - intestazioni
	if intestazioni, err := r.Read(); err != nil {
		log.Fatal("Errore:", err)
	} else {
		fmt.Println("Intestazioni:", intestazioni)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		s = strings.TrimSpace(record[19])
		if s != "" {
			c = strings.TrimSpace(record[5])
			//sceglie fra città metropolitana e provincia
			if prv = strings.TrimSpace(record[10]); prv == "-" {
				prv = strings.TrimSpace(record[11])
				cm = false
			} else {
				cm = true
			}
			cc = append(cc, Comunecodice{
				Comune: c, Codice: s, Provincia: prv,
				Targa:        strings.TrimSpace(record[13]),
				Regione:      strings.TrimSpace(record[9]),
				Incittametro: cm, CoIdx: Normalizza(c),
			})
		}
	}
	fmt.Println("comuni letti:", len(cc))

	sort.Sort(ByCoIdx(cc))

	f, err := os.Create("comuni.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	comuniTemplate.Execute(f, struct {
		Timestamp    time.Time
		URL          string
		Comunecodice []Comunecodice
	}{
		Timestamp:    time.Now(),
		URL:          comuniURL,
		Comunecodice: cc,
	})

}

var comuniTemplate = template.Must(template.New("").Parse(`// go generate
// FILE GENERATO AUTOMATICAMENTE; NON MODIFICARE
// Generato il:
// {{ .Timestamp }}
// usando dati scaricati da:
// {{ .URL }}

package generacodicefiscale
// Comunecodice : array con il codice istat del comune,il nome
// Provincia, SiglaTarga (se esiste, '-' altrimenti), Regione,
// Incittametro:Se è in città metropolitana, CoIdx:Nome comune normalizzato per indice
type Comunecodice struct {
	Codice, Comune, Provincia, Targa, Regione, CoIdx string
	Incittametro bool
}

// Comunecod : codici dei comuni
var Comunecod = []Comunecodice{
{{- range .Comunecodice}}
	{Codice:"{{ .Codice }}",Comune:"{{ .Comune }}", Provincia:"{{ .Provincia }}", Targa:"{{ .Targa }}",
	Regione:"{{ .Regione }}", Incittametro: {{.Incittametro}}, CoIdx:"{{ .CoIdx }}"},
{{- end}}
}
`))
