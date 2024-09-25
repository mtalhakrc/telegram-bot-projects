package model

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

type HaftalikRapor struct {
	GunlukRaporlar []GunlukRapor `json:"gunluk_raporlar"`
}
