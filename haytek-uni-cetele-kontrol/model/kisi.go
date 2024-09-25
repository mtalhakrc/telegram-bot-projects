package model

import (
	"errors"
	"fmt"
)

const cezaAmount = 5

type Kisi struct {
	Isim       string         `json:"isim"`
	Programlar map[string]int `json:"programlar"`
}

func (kisi Kisi) CezaHesapla() (KisiSonuc, []error) {
	var errs []error
	var donecek = KisiSonuc{
		Isim:        kisi.Isim,
		Eksiklikler: map[string]int{},
		Ceza:        0,
	}
	for program, adet := range kisi.Programlar {
		if adet < 0 {
			adet = 0
			str := fmt.Sprintf("Adet negatif olamaz. %s kişisi %s programı için adet 0 olarak ayarlandı.", kisi.Isim, program)
			err := errors.New(str)
			errs = append(errs, err)
		}
		switch program {
		case "Kuran-ı kerim":
			if adet < 1 {
				donecek.Ceza += (1 - adet) * cezaAmount
				donecek.Eksiklikler["Kuran-ı kerim"] = 1 - adet
			}
		case "Mütalaa":
			if adet < 3 {
				donecek.Ceza += (3 - adet) * cezaAmount
				donecek.Eksiklikler["Mütalaa"] = 3 - adet
			}

		case "Cevşen":
			if adet < 5 {
				donecek.Ceza += (5 - adet) * cezaAmount
				donecek.Eksiklikler["Cevşen"] = 5 - adet
			}

		case "Yazı":
			if adet < 11 {
				donecek.Ceza += (11 - adet) * cezaAmount
				donecek.Eksiklikler["Yazı"] = 11 - adet
			}

		}
	}
	return donecek, errs
}

type KisiSonuc struct {
	Isim        string         `json:"isim,omitempty"`
	Eksiklikler map[string]int `json:"eksiklikler,omitempty"`
	Ceza        int            `json:"ceza,omitempty"`
}
