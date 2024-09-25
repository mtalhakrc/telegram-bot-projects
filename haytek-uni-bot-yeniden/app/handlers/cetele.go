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
	"strings"
	"time"
)

type CeteleHandler struct {
	gunlukRaporService service.IGunlukRaporService
	userService        service.IUserService
}

func NewCeteleHandler(db *bun.DB) CeteleHandler {
	return CeteleHandler{
		gunlukRaporService: service.NewGunlukRaporService(db),
		userService:        service.NewUserService(db),
	}
}

func (CeteleHandler) Start(ctx *app.Ctx, params []string) (string, error) {
	return `
HaytekUni Cetele Kontrol Botu çetele kontrolü yapar.
Her gün 11:00 da bir hatırlatma mesajı gönderir.
Her gün 11:00 da kişilere erken kontrol mesajı gönderir.
Her gün 01:00 da çetele kontrolü yapar raporun çıktısını gönderir ve kaydeder.
Kaydolmak için: /kaydol {çeteledeki isminiz}
Kaydınızı silmek için: /kayitsil
Spesifik bir günün ozeti: /gunlukozet {tarih} (only admin)
Hafralik ozet: /haftalikozet (only admin)
`, nil
}
func (s CeteleHandler) GetSpecificRecord(ctx *app.Ctx, params []string) (string, error) {
	if !s.IsAdmin(ctx.SentFrom().ID) {
		return "", errors.New("bu komut için yetkiniz yok")
	}
	if len(params) != 1 {
		return "", errors.New("komuttan sonra tarih berlirtin(ör: /gunlukozet 10.10.2022)")
	}
	parts := strings.Split(params[0], ".")
	if len(parts) != 3 {
		return "", errors.New("hatalı tarih formatı")
	}
	date := parts[2] + "-" + parts[1] + "-" + parts[0]

	res, err := s.gunlukRaporService.GetSpecificDayRecord(date)
	if err != nil {
		log.Println(err)
	}
	str := gunlukRaporFormat(res, nil)
	return str, nil
}
func (s CeteleHandler) GetHaftalikOzet(ctx *app.Ctx, params []string) (string, error) {
	if !s.IsAdmin(ctx.SentFrom().ID) {
		return "", errors.New("bu komut için yetkiniz yok")
	}
	res, err := s.gunlukRaporService.GetLastWeekRecords()
	if err != nil {
		log.Println(err)
		return "", err
	}
	str := haftalikRaporFormat(res)
	return str, nil
}
func (s CeteleHandler) Admins(ctx *app.Ctx, params []string) (string, error) {
	names, err := s.userService.GetAdminsNames(context.Background())
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(names, "\n"), err
}

func haftalikRaporFormat(raporlar []model.GunlukRapor) string {
	var str string
	loc, _ := time.LoadLocation("Europe/Istanbul")
	if len(raporlar) == 0 {
		str = "Bu hafta için kayıtlar bulunamadı."
		return str
	}

	str += fmt.Sprintf("HaytekUni Haftalık Çetele Raporu\t Saat: %s \n", utils.GetSaat(time.Now().In(loc)))
	str += fmt.Sprintf("Bu rapor şu tarihler için çıkarılmıştır: ")
	for _, gunlukraporlar := range raporlar {
		str += fmt.Sprintf("%s - ", utils.GetTarih(gunlukraporlar.Tarih))
	}

	str += "\n\n"

	//bu hafta yapılmayanları bir gunluk modelde toplayacağım
	var raporsonuc = make(map[string]map[string]int)

	for _, gunlukrapor := range raporlar { //raporlarda dön
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

func (s CeteleHandler) IsAdmin(id int64) bool {
	user, err := s.userService.GetByUserID(context.Background(), id)
	if err != nil {
		log.Println(err)
		return false
	}
	return user.Type == model.UserTypeAdmin
}
