package generacodicefiscale

// Generazione codice fiscale 2019 - Squeeze69

//go:generate go run scaricacomuni.go

//go:generate go fmt comuni.go

//go:generate go run scaricanazioni.go

//go:generate go fmt nazioni.go

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/squeeze69/codicefiscale"
)

var vocali = map[rune]bool{
	'A': true, 'E': true, 'I': true, 'U': true, 'O': true,
}

//mappatura numero=LETTERA, da specifiche CF - anno 2017
var mesedinascita = map[int]string{
	1: "A", 2: "B", 3: "C", 4: "D", 5: "E", 6: "H", 7: "L",
	8: "M", 9: "P", 10: "R", 11: "S", 12: "T",
}

//CFGenError : errore di generazione codice fiscale
type CFGenError struct {
	msg string
}

func (c *CFGenError) Error() string {
	return c.msg
}

//EliminaAccenti : elimina in maniera semplice gli accenti - solo sulle minuscole
func EliminaAccenti(s string) string {
	s = regexp.MustCompile("è|é").ReplaceAllString(s, "e")
	s = regexp.MustCompile("à").ReplaceAllString(s, "a")
	s = regexp.MustCompile("ù").ReplaceAllString(s, "u")
	s = regexp.MustCompile("ò").ReplaceAllString(s, "o")
	return regexp.MustCompile("ì").ReplaceAllString(s, "i")
}

//EstrazioneLettere : Estrae le lettere (3) per il cognome ed il nome
func EstrazioneLettere(s string) string {
	var r, c, v string
	rx, _ := regexp.Compile("[^a-zA-Z]")
	s = strings.ToUpper(rx.ReplaceAllString(EliminaAccenti(strings.ToLower(s)), ""))
	for _, l := range s {
		if l == ' ' {
			continue
		}
		if _, ok := vocali[rune(l)]; ok {
			v = v + string(l)
		} else {
			c = c + string(l)
		}
	}
	switch {
	case len(c) < 3:
		r = c
		for _, l := range v {
			r = r + string(l)
			if len(r) == 3 {
				break
			}
		}
		for len(r) < 3 {
			r = r + "X"
		}

	case len(c) > 3:
		r = string(c[0]) + string(c[2]) + string(c[3])
	default:
		r = c[0:3]
	}
	return r[0:3]
}

//genera un errore di tipo CFGenError
func errCFGenError(s string) *CFGenError {
	er := new(CFGenError)
	er.msg = s
	return er
}

//Genera : genera il codice fiscale
//Ingresso: cognome,nome,sesso (M/F),istatcitta:codice ISTAT della città,datadinascita in formato "AAAA-MM-DD"
func Genera(cognome, nome, sesso, istatcitta, datadinascita string) (string, *CFGenError) {
	data, errtime := time.Parse("2006-1-2", datadinascita)
	if errtime != nil {
		return "", errCFGenError("data non valida")
	}
	giorno := data.Day()
	switch {
	case sesso == "F" || sesso == "f":
		giorno = giorno + 40
	case sesso == "M" || sesso == "m":
	default:
		return "", errCFGenError("Genere non valido")
	}
	cf := fmt.Sprintf("%3s%3s%2s%s%02d%4s",
		EstrazioneLettere(cognome), EstrazioneLettere(nome),
		data.Format("06"), mesedinascita[int(data.Month())],
		giorno, strings.ToUpper(istatcitta))
	cc, err := codicefiscale.Codicedicontrollo(cf)
	if err != nil {
		return "", errCFGenError(err.Error() + " " + cf)
	}
	return cf + cc, nil
}
