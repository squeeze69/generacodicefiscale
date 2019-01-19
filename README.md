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