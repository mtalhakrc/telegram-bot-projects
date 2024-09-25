package model

const cezaAmount = 5

type Kisi struct {
	Isim       string         `json:"isim"`
	Programlar map[string]int `json:"programlar"`
}

type KisiSonuc struct {
	Isim        string         `json:"isim,omitempty"`
	Eksiklikler map[string]int `json:"eksiklikler,omitempty"`
	Ceza        int            `json:"ceza,omitempty"`
}
