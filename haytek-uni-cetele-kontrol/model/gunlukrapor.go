package model

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/uptrace/bun"
	"haytekuni-cetele-kontrol/config"
	"haytekuni-cetele-kontrol/database"
	"haytekuni-cetele-kontrol/utils"
	"os"
	"path/filepath"
	"time"
)

var db *bun.DB

type GunlukRapor struct {
	Tarih        time.Time   `json:"tarih"`
	KisilerSonuc []KisiSonuc `json:"kisiler_sonuc"`
}

func (s GunlukRapor) Kaydet() error {
	cfg := config.Get()
	name := utils.GetTarih(s.Tarih) + ".json"
	path := filepath.Join(cfg.Cetele.RaporlarPath, name)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err = errors.New("bu tarih için rapor zaten kaydedilmiş")
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(s)
	if err != nil {
		return err
	}

	//dosya yolunu da databaseye yazıcaz.
	db = database.Get()

	rapordbmodel := GunlukRaporDBModel{
		DosyaYolu: path,
	}
	_, err = db.NewInsert().Model(&rapordbmodel).Exec(context.Background())
	return err

	/*
		SELECT *
		FROM public.users AS "Users"
		WHERE "Users"."created_date" >= NOW() - INTERVAL '24 HOURS'
		ORDER BY "Users"."created_date" DESC
	*/
}

// GunlukRaporDBModel gunluk raporları dosyada tutucaz. dbde dosyanın yolunu tutucaz. güne göre kayıtları alıcaz.
type GunlukRaporDBModel struct {
	BaseModel
	DosyaYolu string `json:"dosya_yolu"`
}
