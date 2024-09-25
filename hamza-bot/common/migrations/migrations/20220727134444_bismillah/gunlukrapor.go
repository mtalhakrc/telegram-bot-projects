package migration

import (
	"github.com/haytek-uni-bot-yeniden/pkg/model"
	"time"
)

type GunlukRapor struct {
	model.BaseModel
	Tarih        time.Time   `json:"tarih"`
	KisilerSonuc []KisiSonuc `json:"kisiler_sonuc"`
}

func (GunlukRapor) Model() string {
	return "gunlukrapor"
}

type KisiSonuc struct {
	Isim        string         `json:"isim,omitempty"`
	Eksiklikler map[string]int `json:"eksiklikler,omitempty"`
	Ceza        int            `json:"ceza,omitempty"`
}
