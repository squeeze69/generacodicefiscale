// scarica l'elenco dei comuni dall'istat, crea un file comuni.go con l'elenco dei comuni, il codice istat ed altre informazioni

package main

import (
	"bytes"
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
	comuniURL           = "https://www.istat.it/storage/codici-unita-amministrative/Elenco-comuni-italiani.csv"
	comuniVariazioniURL = "https://www.anagrafenazionale.interno.it/wp-content/uploads/ANPR_archivio_comuni.csv"
	formatoData         = "2006-01-02"
	dataNonValida       = "1000-01-01"
)

// Comunecodice struttura per memorizzare le informazioni estratte
type Comunecodice struct {
	Codice         string
	Comune         string
	Provincia      string
	Targa          string
	Regione        string
	DataCessazione string
	CoIdx          string
}

// ByCoIdx implementa interface per riordinare l'elenco per comune-indice (sort.Sort...)
type ByCoIdx []Comunecodice

func (a ByCoIdx) Len() int      { return len(a) }
func (a ByCoIdx) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCoIdx) Less(i, j int) bool {
	return strings.Compare(a[i].CoIdx+a[i].Regione, a[j].CoIdx+a[j].Regione) <= 0
}

// per rendere più elegante la riconciliazione targa ->regione e provincia
type regioneprovincia struct {
	Regione, Provincia string
}

// inizializza per "Normalizza"
var ns1 = regexp.MustCompile("è|é")
var ns2 = regexp.MustCompile("à")
var ns3 = regexp.MustCompile("ù")
var ns4 = regexp.MustCompile("ò")
var ns5 = regexp.MustCompile("ì")
var ns6 = regexp.MustCompile("[^a-z]")

// Normalizza : esegue alcune operazioni per permettere di confrontare i nomi in maniera agnostica dalle vocali
func Normalizza(s string) string {
	s = strings.ToLower(s)
	s = ns1.ReplaceAllString(s, "e")
	s = ns2.ReplaceAllString(s, "a")
	s = ns3.ReplaceAllString(s, "u")
	s = ns4.ReplaceAllString(s, "o")
	s = ns5.ReplaceAllString(s, "i")
	return ns6.ReplaceAllString(s, "")
}

func main() {
	var s, c, prv string
	bom3utf8 := []byte{0xef, 0xbb, 0xbf}

	cc := make([]Comunecodice, 0, 11520)

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
	if _, err := r.Read(); err != nil {
		log.Fatal("Errore:", err)
	}
	// tiene traccia dei codici comunali attivi, serve per il passo successivo (codici disattivati)
	codiceattivo := make(map[string]bool, 10000)
	targaRegioneProvincia := make(map[string]regioneprovincia, 80)
	defaultCessazione := "9999-12-31"

	// vecchie exclave italiane
	targaRegioneProvincia["ZA"] = regioneprovincia{Provincia: "Zara", Regione: "Zara"}
	targaRegioneProvincia["PL"] = regioneprovincia{Provincia: "Pola", Regione: "Pola"}
	targaRegioneProvincia["FU"] = regioneprovincia{Provincia: "Fiume", Regione: "Friuli-Venezia Giulia"}

	// targhe mutate nel tempo (diventate rispettivamente: PU e FC)
	targaRegioneProvincia["PS"] = regioneprovincia{Provincia: "Pesaro-Urbino", Regione: "Marche"}
	targaRegioneProvincia["FO"] = regioneprovincia{Provincia: "Forlì-Cesena", Regione: "Emilia-Romagna"}

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
			prv = strings.TrimSpace(record[11])
			codiceattivo[s] = true
			tg := strings.TrimSpace(record[14])
			rg := strings.TrimSpace(record[10])
			targaRegioneProvincia[tg] = regioneprovincia{Regione: rg, Provincia: prv}
			cc = append(cc, Comunecodice{
				Comune: c, Codice: s, Provincia: prv,
				Targa:          tg,
				Regione:        rg,
				DataCessazione: defaultCessazione,
				CoIdx:          Normalizza(c),
			})
		}
	}
	fmt.Println("Comuni attivi letti:", len(cc))

	// caccia ai comunni cessati (di solito trasformati in frazioni di altri comuni)
	response1, err := http.Get(comuniVariazioniURL)
	if err != nil {
		log.Fatal("Errore", err)
	}
	defer response1.Body.Close()
	// legge dal csv
	// va fatto lo skip del BOM a 3 bytes se presente
	ra, err := io.ReadAll(response1.Body)
	if err != nil {
		log.Fatal("Errore:", err)
	}

	rb := bytes.NewReader(ra)
	if bytes.Equal(bom3utf8[:], ra[:3]) {
		rb.Seek(3, io.SeekStart)
	}
	r = csv.NewReader(rb)
	r.Comma = ','
	// evita la prima linea - intestazioni
	if _, err := r.Read(); err != nil {
		log.Fatal("Errore:", err)
	}

	var dc string
	var datenonbuone int
	//mappa codice ISTAT, le variazini possono essere più di una
	inattivi := make(map[string]Comunecodice, 4000)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// codice del comune
		s = strings.TrimSpace(record[4])
		if s != "" {
			c = strings.TrimSpace(record[5])
			// salta il codice se già fra le città attive
			if _, ok := codiceattivo[s]; ok {
				continue
			}
			_, err := time.Parse(formatoData, record[2])
			if err != nil {
				datenonbuone++
				dc = dataNonValida
			} else {
				dc = record[2]
			}
			tg := strings.TrimSpace(record[14])
			rp, ok := targaRegioneProvincia[tg]
			if !ok {
				rp = regioneprovincia{Provincia: "?", Regione: "?"}
			}
			comuneattuale := Comunecodice{
				Comune: c, Codice: s, Provincia: rp.Provincia,
				Targa:          tg,
				Regione:        rp.Regione,
				DataCessazione: dc,
				CoIdx:          Normalizza(c),
			}
			ca, ok := inattivi[s]
			if ok {
				// prende la modifica più recente
				if strings.Compare(ca.DataCessazione, comuneattuale.DataCessazione) < 0 {
					inattivi[s] = comuneattuale
				}
			} else {
				inattivi[s] = comuneattuale
			}
		}
	}

	// appende tutti i comuni inattivi una volta sola
	for _, i := range inattivi {
		cc = append(cc, i)
	}

	fmt.Printf("Comuni attivi+inattivi: %d\n", len(cc))
	if datenonbuone > 0 {
		fmt.Printf("Trovate %d date non valide\n", datenonbuone)
	}
	sort.Sort(ByCoIdx(cc))

	f, err := os.Create("comuni.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	comuniTemplate.Execute(f, struct {
		Timestamp    time.Time
		URL          string
		URL1         string
		Comunecodice []Comunecodice
	}{
		Timestamp:    time.Now(),
		URL:          comuniURL,
		URL1:         comuniVariazioniURL,
		Comunecodice: cc,
	})

}

var comuniTemplate = template.Must(template.New("").Parse(`// go generate
// FILE GENERATO AUTOMATICAMENTE; NON MODIFICARE
// Generato il:
// {{ .Timestamp }}
// usando dati scaricati da:
// {{ .URL }}
// {{ .URL1 }}

package generacodicefiscale
// Comunecodice : array con il codice istat del comune,il nome
// Provincia, SiglaTarga (se esiste, '-' altrimenti), Regione,
// DataCessazione : data di cessazione del comune 9999-12-31 se attivo
// usare time.Parse("2006-01-02", ...)
// CoIdx:Nome comune normalizzato per indice
type Comunecodice struct {
	Codice, Comune, Provincia, Targa, Regione, CoIdx string
	DataCessazione string
}

// Comunecod : codici dei comuni
var Comunecod = []Comunecodice{
{{- range .Comunecodice}}
	{Codice:"{{ .Codice }}",Comune:"{{ .Comune }}", Provincia:"{{ .Provincia }}",
	Targa:"{{ .Targa }}", Regione:"{{ .Regione }}",
	DataCessazione: "{{.DataCessazione}}", CoIdx:"{{ .CoIdx }}"},
{{- end}}
}
`))
