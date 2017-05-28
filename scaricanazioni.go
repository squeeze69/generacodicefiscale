// scarica l'elenco delle nazioni - sfrutta un file dal ministero della salute
// Esiste anche un file ISTAT ufficiale

// +build ignore


package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
	"sort"
)

const (
	nazioneURL = "http://www.salute.gov.it/imgs/C_17_pubblicazioni_1055_ulterioriallegati_ulterioreallegato_0_alleg.txt"
)

type Nazionecodice struct {
	Codice     string
	CodiceISO  string
	CodiceISO3 string
	Nazione    string
}
type ByNazione []Nazionecodice

func (a ByNazione) Len() int           { return len(a) }
func (a ByNazione) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNazione) Less(i, j int) bool { return strings.Compare(a[i].Nazione, a[j].Nazione)<= 0 }

func main() {
	var s string
	cc := make([]Nazionecodice, 0, 300)

	response, err := http.Get(nazioneURL)
	if err != nil {
		log.Fatal("Errore", err)
	}
	defer response.Body.Close()
	rv := charmap.ISO8859_1.NewDecoder().Reader(response.Body)
	r := csv.NewReader(rv)
	r.Comma = '\t'
	if _, err := r.Read(); err != nil {
		log.Fatal("Errore:", err)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		//praticamente ignoro tutti gli altri errori
		//(altrimenti è instabile distinguere fra altri tipi di errori dei csv
		if err != nil {
			continue
		}
		if len(record) != 4 {
			log.Fatal("Qualche problema", record)
		}
		s = strings.TrimSpace(record[3])
		if s != "" && s != "ND" {
			cc = append(cc,
				Nazionecodice{
					Nazione:    strings.TrimSpace(record[2]),
					Codice:     s,
					CodiceISO:  strings.TrimSpace(record[0]),
					CodiceISO3: strings.TrimSpace(record[1])})
		}
	}
	fmt.Println("Nazioni importate:", len(cc))
	sort.Sort(ByNazione(cc))
	
	f, err := os.Create("nazioni.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	nazioniTemplate.Execute(f, struct {
		Timestamp     time.Time
		URL           string
		Nazionecodice []Nazionecodice
	}{
		Timestamp:     time.Now(),
		URL:           nazioneURL,
		Nazionecodice: cc,
	})

}

var nazioniTemplate = template.Must(template.New("").Parse(`// go generate
// FILE GENERATO AUTOMATICAMENTE; NON MODIFICARE
// Generato il:
// {{ .Timestamp }}
// usando dati scaricati da:
// {{ .URL }}

package generacodicefiscale
// Nazionecodice : array con il codice istat della nazione ed i relativi codici ISO
type Nazionecodice struct {
	Codice, Nazione, CodiceISO,CodiceISO3 string
}
var Nazionecod = []Nazionecodice{
{{- range .Nazionecodice}}
	{Codice:"Z{{ .Codice }}",Nazione:"{{ .Nazione }}",CodiceISO:"{{ .CodiceISO }}",CodiceISO3:"{{ .CodiceISO3 }}"},
{{- end}}
}
`))
