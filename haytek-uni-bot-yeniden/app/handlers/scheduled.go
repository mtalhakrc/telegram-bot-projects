package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/haytek-uni-bot-yeniden/app/service"
	"github.com/haytek-uni-bot-yeniden/common/model"
	"github.com/haytek-uni-bot-yeniden/pkg/app"
	"github.com/haytek-uni-bot-yeniden/pkg/utils"
	"github.com/uptrace/bun"
	"log"
	"strconv"
	"time"
)

type ScheduledHandler struct {
	userService        service.IUserService
	sheetsservice      service.ISheetsService
	gunlukraporService service.IGunlukRaporService
}

func NewScheduled(db *bun.DB, s service.ISheetsService) ScheduledHandler {
	return ScheduledHandler{
		userService:        service.NewUserService(db),
		gunlukraporService: service.NewGunlukRaporService(db),
		sheetsservice:      s,
	}
}

func (u ScheduledHandler) CeteleHatirlatmaMesaji() []app.ScheduledResponse {
	userids, err := u.userService.GetAllUserIDs(context.Background())
	if err != nil {
		log.Println(err)
		return nil
	}
	var response []app.ScheduledResponse
	for _, id := range userids {
		response = append(response, app.ScheduledResponse{
			UserID: id,
			Result: "Çetelelerinizi saat 01:00'dan önce doldurunuz!",
			Error:  nil,
		})
	}
	return response
}

func (u ScheduledHandler) GunlukRaporMesaji() []app.ScheduledResponse {
	gunlukRapor, kisierrs := newGunlukRapor(u.sheetsservice, "Dün Özet")
	err := u.gunlukraporService.Kaydet(gunlukRapor)
	if err != nil {
		log.Println(err)
	}
	str := gunlukRaporFormat(gunlukRapor, kisierrs)
	return []app.ScheduledResponse{
		{
			UserID: -1001721260016, //haytek uni
			Result: str,
			Error:  nil,
		},
	}
}

func (u ScheduledHandler) GunlukErkenKontrolMesaji() []app.ScheduledResponse {
	gunlukrapor, _ := newGunlukRapor(u.sheetsservice, "Dün Özet")
	var donecek []app.ScheduledResponse
	for _, kisisonuc := range gunlukrapor.KisilerSonuc {
		user, err := u.userService.GetByName(context.Background(), kisisonuc.Isim)
		if err != nil {
			log.Println(err)
			continue
		}
		str := gunSonuErkenKontrolFormat(kisisonuc)
		donecek = append(donecek, app.ScheduledResponse{
			UserID: user.UserID,
			Result: str,
			Error:  nil,
		})
	}
	return donecek
}

func newGunlukRapor(s service.ISheetsService, range_ string) (model.GunlukRapor, []error) {
	resp, err := s.GetFromSheet(range_)
	if err != nil {
		log.Println(err)
		return model.GunlukRapor{}, []error{err}
	}
	kisiler, err := parseKisilerFromSheetResponse(resp)
	if err != nil {
		log.Println(err)
		return model.GunlukRapor{}, []error{err}
	}

	var kisisonuclar []model.KisiSonuc
	var kisierrs []error

	for _, kisi := range kisiler {
		kisisonuc, errs := kisiSonucHesapla(kisi)
		kisierrs = append(kisierrs, errs...)
		kisisonuclar = append(kisisonuclar, kisisonuc)
	}

	var gunlukRapor model.GunlukRapor
	gunlukRapor.KisilerSonuc = kisisonuclar
	gunlukRapor.Tarih = time.Now()
	return gunlukRapor, kisierrs
}

func parseKisilerFromSheetResponse(resp [][]interface{}) ([]model.Kisi, error) {
	var kisiler []model.Kisi
	var programlar = make(map[string]int)
	var adet int
	var err error
	// col0 [Kişiler Tarih Kuran-ı kerim Mütalaa Cevşen Yazı] bu yüzden 1. coldan başlıyoruz
	for i := 1; i < len(resp); i++ {
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
		for index := 2; index < len(resp[i]); index++ {
			if resp[i][index] == "" {
				resp[i][index] = "0"
			}
			adet, err = strconv.Atoi(resp[i][index].(string))
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
			Isim:       resp[i][0].(string),
			Programlar: programlar,
		})
	}
	return kisiler, nil
}

const cezaAmount = 5

func kisiSonucHesapla(kisi model.Kisi) (model.KisiSonuc, []error) {
	var errs []error
	var donecek = model.KisiSonuc{
		Isim:        kisi.Isim,
		Eksiklikler: map[string]int{},
		Ceza:        0,
	}
	for program, adet := range kisi.Programlar {
		if adet < 0 {
			adet = 0
			//str := fmt.Sprintf("Adet negatif olamaz. %s kişisi %s programı için adet 0 olarak ayarlandı.", kisi.Isim, program)
			//err := errors.New(str)
			//errs = append(errs, err)
			errs = append(errs, errors.New(fmt.Sprintf("Adet negatif olamaz. %s kişisi %s programı için adet 0 olarak ayarlandı.", kisi.Isim, program)))
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
func gunlukRaporFormat(rapor model.GunlukRapor, errs []error) string {
	var str string
	var eksikMi bool
	loc, _ := time.LoadLocation("Europe/Istanbul")
	str += fmt.Sprintf("HaytekUni %s Tarihli Gun Sonu Ozeti\n", utils.GetTarih(rapor.Tarih))
	str += fmt.Sprintf("Saat: %s \n", utils.GetSaat(rapor.Tarih.In(loc)))
	str += "------------------------- \n"

	for _, kisi := range rapor.KisilerSonuc {
		if len(kisi.Eksiklikler) > 0 {
			eksikMi = true
			str += fmt.Sprintf("%s\n", kisi.Isim)
			for program, adet := range kisi.Eksiklikler {
				str += fmt.Sprintf("%s: %d\n", program, adet)
			}
			str += " ------------------------- \n"
		}
	}
	if !eksikMi {
		str += "Bugün için eksiklik yok! 🎉🎉🎉"
	}
	for _, err := range errs {
		str += fmt.Sprintf("%s\n", err.Error())
	}
	return str
}

func gunSonuErkenKontrolFormat(kisi model.KisiSonuc) string {
	var str string
	loc, _ := time.LoadLocation("Europe/Istanbul")
	str += fmt.Sprintf("Gun Sonu Erken Kontrol Mesajı  \nSaat: %s \t Tarih: %s\n", utils.GetSaat(time.Now().In(loc)), utils.GetTarih(time.Now()))
	var eksikMi = false
	if len(kisi.Eksiklikler) > 0 {
		eksikMi = true
		str += fmt.Sprintf("%s\n", kisi.Isim)
		for program, adet := range kisi.Eksiklikler {
			str += fmt.Sprintf("%s: %d\n", program, adet)
		}
		str += fmt.Sprintf("CEZANIZ: %d TL\n", kisi.Ceza)
	}
	if !eksikMi {
		str += "Çeteleniz tam!\n"
	}
	return str
}
