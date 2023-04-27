// ricerca all'interno degli array i comuni e nazioni

package generacodicefiscale

import (
	"regexp"
	"sort"
	"strings"
)

// CFSearchError errore nella ricerca
type CFSearchError struct {
	msg string
}

func (r *CFSearchError) Error() string {
	return r.msg
}

//inizializza per "Normalizza"
var ns1=regexp.MustCompile("è|é")
var ns2=regexp.MustCompile("à")
var ns3=regexp.MustCompile("ù")
var ns4=regexp.MustCompile("ò")
var ns5=regexp.MustCompile("ì")
var ns6=regexp.MustCompile("[^a-z]")

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

// CercaComune all'interno dell'array - normalizza prima, per evitare problemi con spazi, simboli od altro
// in ingresso: nome del comune
// in uscita: voce dell'array relativa o nil se non trovato, errore: nil o CFSearchError se non trovato
func CercaComune(a string) (*Comunecodice, *CFSearchError) {
	na := Normalizza(a)
	r := sort.Search(len(Comunecod), func(i int) bool { return strings.Compare(Comunecod[i].CoIdx, na) >= 0 })
	if r < len(Comunecod) {
		if strings.Compare(Comunecod[r].CoIdx, na) == 0 {
			return &Comunecod[r], nil
		}
	}
	er := new(CFSearchError)
	er.msg = "Non trovato"
	return nil, er

}
