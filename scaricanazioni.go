// scarica l'elenco delle nazioni - sfrutta un file dal ministero della salute
// Esiste anche un file ISTAT ufficiale

// +build ignore

package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"go.dev/x/text/encoding/charmap"
)

const (
	nazioneURL     = "https://www.istat.it/it/files//2011/01/Elenco-codici-e-denominazioni-unita-territoriali-estere.zip"
	fileDaGenerare = "nazioni.go"
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
func (a ByNazione) Less(i, j int) bool { return strings.Compare(a[i].Nazione, a[j].Nazione) <= 0 }

//leggiCSVinZIP : legge il primo csv contenuto in uno zip file scaricato al volo in memoria
func leggiCSVinZIP(url string) (data []byte, err error) {

	r, er := http.Get(url)
	if er != nil {
		log.Fatal(er)
	}
	defer r.Body.Close()
	tutto, er := io.ReadAll(r.Body)
	if er != nil {
		log.Fatal(er)
	}
	zipR, er := zip.NewReader(bytes.NewReader(tutto), int64(len(tutto)))
	if er != nil {
		log.Fatal(er)
	}

	// cerca il CSV
	for _, zipf := range zipR.File {
		if strings.Contains(strings.ToLower(zipf.Name), ".csv") {
			zf, er := zipf.Open()
			if er != nil {
				log.Fatal(er)
			}
			defer zf.Close()
			return io.ReadAll(zf)
		}
	}
	return nil, errors.New("non è stato trovato niente")
}

func main() {
	var s string
	var lette int
	cc := make([]Nazionecodice, 0, 300)

	response, err := leggiCSVinZIP(nazioneURL)
	if err != nil {
		log.Fatal(err)
	}
	rfile := bytes.NewReader(response)
	//reader per la decodifica da Windows 1252/ISO8859_1 a UTF-8
	//rv := charmap.ISO8859_1.NewDecoder().Reader(response.Body)
	rv := charmap.ISO8859_1.NewDecoder().Reader(rfile)
	r := csv.NewReader(rv)
	r.Comma = ';'
	if intestazioni, err := r.Read(); err != nil {
		log.Fatal("Errore:", err)
	} else {
		if intestazioni[9] != "Codice AT" {
			log.Fatal("Non ho trovato la colonna con il 'Codice AT'")
		}
		if intestazioni[6] != "Denominazione IT" {
			log.Fatal("Non ho trovato la colonna con il 'Denominazione IT'")
		}
		if intestazioni[11] != "Codice ISO 3166 alpha2" {
			log.Fatal("Non ho trovato la colonna con il 'Codice ISO 3166 alpha2'")
		}
		if intestazioni[12] != "Codice ISO 3166 alpha3" {
			log.Fatal("Non ho trovato la colonna con il 'Codice ISO 3166 alpha3'")
		}
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
		lette += 1
		if len(record) < 13 {
			continue
		}
		s = strings.TrimSpace(record[9])
		if s != "" && s != "n.d." {
			cc = append(cc,
				Nazionecodice{
					Nazione:    strings.TrimSpace(record[6]),
					Codice:     s,
					CodiceISO:  strings.TrimSpace(record[11]),
					CodiceISO3: strings.TrimSpace(record[12])})
		}
	}
	fmt.Printf("Nazioni importate: %d, lette: %d\n", len(cc), lette)
	sort.Sort(ByNazione(cc))

	f, err := os.Create(fileDaGenerare)
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
// Nazionecod : array con i dati, nazione per nazione
var Nazionecod = []Nazionecodice{
{{- range .Nazionecodice}}
	{Codice:"{{ .Codice }}",Nazione:"{{ .Nazione }}",CodiceISO:"{{ .CodiceISO }}",CodiceISO3:"{{ .CodiceISO3 }}"},
{{- end}}
}
`))
