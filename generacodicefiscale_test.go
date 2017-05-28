package generacodicefiscale

import (
	"fmt"
	"log"
	"testing"
)

func TestEstrazioneLettere(t *testing.T) {
	nomiTest := map[string]string{
		"MRC": "MARCò", "LCU": "LùCA", "LUX": "LU", "GFR": "GIAN FRANCO", "TZN": "TIZIANA",
		"LRT": "ALBERTO", "RSS": "ROSSI", "XXX": "", "DBR": "Adal Berto", "LCE": "LÈca",
		"LMR": "Laura, Maria", "FCS": "Felice Stefano", "CRL": "Carlo", "BCH": "Bianchi",
	}
	fmt.Println("- Test EstrazioneLettere")
	for r, n := range nomiTest {
		s := EstrazioneLettere(n)
		if s != r {
			t.Errorf("Da '%s' atteso '%s' ottenuto %s\n", n, r, s)
		}
		fmt.Printf("Ok risultato: \"%s\" da \"%s\"\n", s, n)
	}

}

type testStruct struct {
	CFAtteso, Cognome, Nome, Sesso, Istatcitta, Datadinascita string
}

func TestGenera(t *testing.T) {
	ts := []testStruct{
		{CFAtteso: "MRNMRT91R51G388N", Cognome: "Moroni", Nome: "Maruta", Sesso: "F", Istatcitta: "g388", Datadinascita: "1991-10-11"},
		{CFAtteso: "MROTRA92B01F205P", Cognome: "Mòro", Nome: "Tàru", Sesso: "M", Istatcitta: "F205", Datadinascita: "1992-2-1"},
	}
	fmt.Println("- Test Genera")
	for _, s := range ts {
		r, err := Genera(s.Cognome, s.Nome, s.Sesso, s.Istatcitta, s.Datadinascita)
		if err != nil {
			t.Errorf("Errore: %s\n", err)
		}
		if r != s.CFAtteso {
			t.Errorf("Ko, codice non corrisponde (ottenuto: \"%s\" atteso \"%s\"\n", r, s.CFAtteso)
		}
		fmt.Printf("Ok - corrisponde: %s da %s,%s\n", r, s.Cognome, s.Nome)
	}
}

func TestCercaComune(t *testing.T) {
	type TestCerca struct {
		Codice       string
		Comune       string
		ErroreAtteso bool
	}
	ts := []TestCerca{
		{Codice: "F205", Comune: "Milano", ErroreAtteso: false},
		{Codice: "A115", Comune: "Alà dei Sardi", ErroreAtteso: false},
		{Codice: "XXXX", Comune: "Inesistente", ErroreAtteso: true},
	}
	fmt.Println("- Test CercaComune")
	for _, n := range ts {
		i, err := CercaComune(n.Comune)
		if err != nil {
			if n.ErroreAtteso {
				fmt.Println("Ok - Errore atteso")
			} else {
				log.Fatalf("Ko - %s, errore: %s\n", n.Comune, err)
			}
		} else {
			if n.Codice == i.Codice {
				fmt.Printf("Ok - \"%s\" Ottenuto: \"%s\"\n", n.Comune, i.Codice)
			} else {
				log.Fatalf("Ko - da \"%s\" atteso \"%s\" ottenuto: \"%s\"\n", n.Comune, n.Codice, i.Codice)
			}
		}
	}
	
}
