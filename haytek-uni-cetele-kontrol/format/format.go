package format

import (
	"fmt"
	"haytekuni-cetele-kontrol/model"
	"haytekuni-cetele-kontrol/utils"
	"time"
)

//paketin amacı herhangi bir veriyi basabilmek için illaki onu bi string formatına çevirmek lazım. bu paket verileri stringe formatlar

func GunlukRaporFormat(rapor model.GunlukRapor, errss map[string][]error) string {
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
	//if errs != nil {
	//	for _, err := range errs {
	//		str += fmt.Sprintf("%s\n", err.Error())
	//	}
	//}
	if len(errss) > 0 {
		for isim, errs := range errss {
			for _, err := range errs {
				str += fmt.Sprintf("%s: %s\n", isim, err.Error())
			}
		}
	}

	return str
}

func HaftalikRaporFormat(haftalikrapor model.HaftalikRapor) string {
	var str string
	loc, _ := time.LoadLocation("Europe/Istanbul")
	if len(haftalikrapor.GunlukRaporlar) == 0 {
		str = "Bu hafta için kayıtlar bulunamadı."
		return str
	}

	str += fmt.Sprintf("HaytekUni Haftalık Çetele Raporu\t Saat: %s \n", utils.GetSaat(time.Now().In(loc)))
	str += fmt.Sprintf("Bu rapor şu tarihler için çıkarılmıştır: ")
	for _, gunlukraporlar := range haftalikrapor.GunlukRaporlar {
		str += fmt.Sprintf("%s - ", utils.GetTarih(gunlukraporlar.Tarih))
	}
	str += "\n\n"

	//bu hafta yapılmayanları bir gunluk modelde toplayacağım
	var raporsonuc = make(map[string]map[string]int)

	for _, gunlukrapor := range haftalikrapor.GunlukRaporlar { //raporlarda dön
		for _, kisi := range gunlukrapor.KisilerSonuc { //kisilerde dön
			//kişi için eksiklik map başlat
			if _, ok := raporsonuc[kisi.Isim]; !ok {
				raporsonuc[kisi.Isim] = make(map[string]int)
			}
			if len(kisi.Eksiklikler) > 0 {
				for program, adet := range kisi.Eksiklikler { //kisiye ait eksikliklerde dön
					raporsonuc[kisi.Isim][program] += adet
				}
			}
		}
	}
	//fmt.Println(raporsonuc)

	var eksikMi = false
	str += "------------------------- \n"
	for kisi, eksiklikler := range raporsonuc {
		if len(eksiklikler) > 0 {
			eksikMi = true
			str += fmt.Sprintf("%s\n", kisi)
			for program, adet := range eksiklikler {
				str += fmt.Sprintf("%s: %d\n", program, adet)
			}
			str += " ------------------------- \n"
		}
	}
	if !eksikMi {
		str += "Bu hafta eksiklik yok! MaşellahBea"
	}

	return str
}

func PersonalRaporFormat(kisi model.KisiSonuc, errs ...error) string {
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

	if errs != nil {
		for _, err := range errs {
			str += fmt.Sprintf("%s\n", err.Error())
		}
	}

	if !eksikMi {
		str += "Çeteleniz tam!\n"
	}
	return str
}
