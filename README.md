# Generazione codice fiscale in  [GO](https://golang.org)

[![Build Status](https://travis-ci.org/squeeze69/generacodicefiscale.svg?branch=master)](https://travis-ci.org/squeeze69/generacodicefiscale)

## **package**: github.com/squeeze69/generacodicefiscale

## --- NON si danno garanzie ---

### dipende dal package [codicefiscale](https://github.com/squeeze69/codicefiscale)

Uso:

``` go
package main

import (
    "github.com/squeeze69/generacodicefiscale"
    "fmt"
    "log"
)

func main() {
    codicecitta,erc := generacodicefiscale.CercaComune("Milano")
    if erc != nil {
        log.Fatal(erc)
    }
    fmt.Println(codicecitta)
    cf,erg := generacodicefiscale.Genera("Cognome","Nome","M",codicecitta.Codice,"2017-05-1")
    if erg != nil {
        log.Fatal(erg)
    }
    fmt.Println(cf)
}
```

**Note:**

- Per la corretta generazione del codice fiscale vengono rimossi gli accenti dalle vocali accentate, al momento sono supportate solo le "èéàùìò" (maiuscole o minuscole)

- La ricerca del comune avviene tramite CercaComune ed il nome esatto della città (non importano spazi, vocali accentate od altri simboli, la chiave di ricerca viene "normalizzata")

- Per ricerche più sofisticate, è a disposizione Cittacod []Cittacodice, definito in "comuni.go", è ordinato per "CoIdx", ottenuto rimuovendo gli accenti e tutto quel che non è un carattere alfabetico.

- Per cercare la nazione, prego, iterare su Nazionecod []Nazionecodice, definito in "nazioni.go"

- Se si vuole ri-scaricare l'elenco dei comuni e delle nazioni fate "go generate", se tutto andrà bene "nazioni.go" e "comuni.go" verranno rigenerati, anche se ==NON== è consigliato di farlo per la nazione, attualmente il file da cui venivano prese le informazioni è stato modificato togliendo i suddetti codici. Sto lavorando ad una soluzione.

- Un pezzo di codice interessante, Codici Fiscali a parte, per avere almeno un piccolo indizio di quanto sia comodo e poternte GO, la (brutta) funzione qua sotto scarica il file .zip (in memoria), poi fa una scansione dei file presenti finché non trova un file con ".csv" nel nome (potrebbe controllare che finisca per .csv, ma visto lo scopo molto preciso e l'uso di nomi stile windows, basta ed avanza), la funzione si trova in "scaricanazioni.go".

``` go
//leggiCSVinZIP : legge il primo csv contenuto in uno zip file scaricato al volo in memoria
func leggiCSVinZIP(url string) (data []byte, err error) {

    r, er := http.Get(url)
    if er != nil {
        log.Fatal(er)
    }
    defer r.Body.Close()
    tutto, er := ioutil.ReadAll(r.Body)
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
            return ioutil.ReadAll(zf)
        }
    }
    return nil, errors.New("non è stato trovato niente")
}

```