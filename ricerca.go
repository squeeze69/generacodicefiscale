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

// Normalizza : esegue alcune operazioni per permettere di confrontare i nomi in maniera agnostica dalle vocali
func Normalizza(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile("è|é").ReplaceAllString(s, "e")
	s = regexp.MustCompile("à").ReplaceAllString(s, "a")
	s = regexp.MustCompile("ù").ReplaceAllString(s, "u")
	s = regexp.MustCompile("ò").ReplaceAllString(s, "o")
	s = regexp.MustCompile("ì").ReplaceAllString(s, "i")
	return regexp.MustCompile("[^a-z]").ReplaceAllString(s, "")
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
