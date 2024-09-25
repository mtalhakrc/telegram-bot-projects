package migration

// GunlukRaporDBModel gunluk raporları dosyada tutucaz. dbde dosyanın yolunu tutucaz. güne göre kayıtları alıcaz.
type GunlukRaporDBModel struct {
	BaseModel
	DosyaYolu string `json:"dosya_yolu"`
}
