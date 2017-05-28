# Generazione codice fiscale in  [GO](https://golang.org)
[![Build Status](https://travis-ci.org/squeeze69/generacodicefiscale.svg?branch=master)](https://travis-ci.org/squeeze69/generacodicefiscale)
## **package**: github.com/squeeze69/generacodicefiscale
## --- NON si danno garanzie ---
### dipende dal package [codicefiscale](https://github.com/squeeze69/codicefiscale)

Uso:

```
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