package app

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/uptrace/bun"
	"google.golang.org/api/sheets/v4"
	"haytekuni-cetele-kontrol/config"
	"haytekuni-cetele-kontrol/database"
	"haytekuni-cetele-kontrol/logx"
	"haytekuni-cetele-kontrol/model"
	"haytekuni-cetele-kontrol/service"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var srv *sheets.Service
var cfg *config.Config
var db *bun.DB

func getKisilerFromSheet(range_ string) ([]model.Kisi, error) {
	//kişilerin dünkü özetlerini (aws 3 saat ileride olduğu için ona göre bugünün özeti dünün özeti oluyor) al
	//özet sayfası
	cfg = config.Get()
	srv = service.Get()
	resp, err := srv.Spreadsheets.Values.Get(cfg.Cetele.SpreadSheetID, range_).Do()
	if err != nil {
		return nil, err
	}
	var kisiler []model.Kisi
	var programlar = make(map[string]int)
	var adet int
	// col0 [Kişiler Tarih Kuran-ı kerim Mütalaa Cevşen Yazı] bu yüzden 1. coldan başlıyoruz
	for i := 1; i < len(resp.Values); i++ {
		//önce hiç yapılmamış olarak kabul edilecek. sonra yaptıkları varsa güncellenecek.
		programlar = map[string]int{
			"Kuran-ı kerim": 0,
			"Mütalaa":       0,
			"Cevşen":        0,
			"Yazı":          0,
		}
		// [Talha Karaca 04.11.2022 1 3 5 11]
		// 0 -> isim
		// 1 -> tarih
		// 2 -> Kuran-ı kerim
		// 3 -> Mütalaa
		// 4 -> Cevşen
		// 5 -> Yazı

		//kişi girmezse index 2(isim, tarih) + girdiği kadar index geliyor. hepsini giren kişi 6 tane oluyor.
		//programlar 2. indexten başlıyor.
		for index := 2; index < len(resp.Values[i]); index++ {
			if resp.Values[i][index] == "" {
				resp.Values[i][index] = "0"
			}
			adet, err = strconv.Atoi(resp.Values[i][index].(string))
			if err != nil {
				return nil, err
			}
			switch index {
			case 2:
				programlar["Kuran-ı kerim"] = adet
			case 3:
				programlar["Mütalaa"] = adet
			case 4:
				programlar["Cevşen"] = adet
			case 5:
				programlar["Yazı"] = adet
			}
		}
		kisiler = append(kisiler, model.Kisi{
			Isim:       resp.Values[i][0].(string),
			Programlar: programlar,
		})
	}
	return kisiler, nil
}

func GunlukRaporHazirla(range_ string) (model.GunlukRapor, map[string][]error) {
	kisiler, err := getKisilerFromSheet(range_)
	if err != nil {
		logx.SendLog(err.Error())
		log.Fatal(err)
	}
	kisilersonuc, errs := getKisilerSonuc(kisiler)

	var gunlukRapor model.GunlukRapor
	gunlukRapor.KisilerSonuc = kisilersonuc
	gunlukRapor.Tarih = time.Now()
	return gunlukRapor, errs
}

//func getKisiData(isim string) ([][]interface{}, error) {
//	resp, err := srv.Spreadsheets.Values.Get(cfg.Cetele.SpreadSheetID, isim).Do()
//	if err != nil {
//		log.Printf("Unable to retrieve data from sheet: %v", err)
//		return nil, err
//	}
//	if len(resp.Values) == 0 {
//		return nil, nil
//	}
//	return resp.Values, nil
//}

func getKisilerSonuc(kisiler []model.Kisi) ([]model.KisiSonuc, map[string][]error) {
	var result []model.KisiSonuc
	var ResultErrs = make(map[string][]error)
	for _, kisi := range kisiler {
		res, errs := kisi.CezaHesapla()
		result = append(result, res)
		//ResultErrs = append(ResultErrs, errs...)
		ResultErrs[kisi.Isim] = errs

	}
	return result, ResultErrs
}

// GunlukOzetAl geçmiş zamandaki spesifik bir günün  kaydını getirir.
func GunlukOzetAl(tarih string) (model.GunlukRapor, error) {
	db = database.Get()
	var raporDBModel model.GunlukRaporDBModel
	var rapor model.GunlukRapor
	err := db.NewSelect().Model(&raporDBModel).Where("created_at::TIMESTAMP::DATE = ?", tarih).Scan(context.Background())
	if err != nil || raporDBModel.DosyaYolu == "" {
		//kaydı alamazsan dön
		return rapor, err
	}

	//eğer kayıt varsa kayıttaki dosyayı aç ve gönder
	f, err := os.Open(raporDBModel.DosyaYolu)
	if err != nil {
		return rapor, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return rapor, err
	}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&rapor)
	if err != nil {
		return rapor, err
	}

	return rapor, nil
}

func HaftalikOzetAl() (model.HaftalikRapor, error) {
	db = database.Get()
	var raporlarDBModel []model.GunlukRaporDBModel
	var haftalikRapor model.HaftalikRapor
	//SELECT *
	//	FROM public.users AS "Users"
	//WHERE "Users"."created_date" >= NOW() - INTERVAL '24 HOURS'
	//ORDER BY "Users"."created_date" DESC
	err := db.NewSelect().Model(&raporlarDBModel).Where("created_at  >= NOW() - INTERVAL '1 WEEKS'").Scan(context.Background())
	if err != nil {
		return haftalikRapor, err
	}

	//eğer kayıt varsa kayıttaki dosyayı aç ve gönder
	var f *os.File
	for _, raporDBModel := range raporlarDBModel {
		f, err = os.Open(raporDBModel.DosyaYolu)
		if err != nil {
			return haftalikRapor, err
		}
		b, err := io.ReadAll(f)
		if err != nil {
			return haftalikRapor, err
		}
		var rapor model.GunlukRapor
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&rapor)
		if err != nil {
			return haftalikRapor, err
		}
		haftalikRapor.GunlukRaporlar = append(haftalikRapor.GunlukRaporlar, rapor)
	}
	return haftalikRapor, nil
}
