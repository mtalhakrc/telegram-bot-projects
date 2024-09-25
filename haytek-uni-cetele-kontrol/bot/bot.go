package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"haytekuni-cetele-kontrol/app"
	"haytekuni-cetele-kontrol/config"
	"haytekuni-cetele-kontrol/format"
	"haytekuni-cetele-kontrol/logx"
	"haytekuni-cetele-kontrol/utils"
	"log"
	"time"
)

var cfg *config.Config

func InitBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI("5325031941:AAHSdWLZKX-2yobRnXIW9rRUH64tDObEEsc")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)
	return bot
}

func ListenForUpdates(bot *tgbotapi.BotAPI) {
	cfg = config.Get()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "Merhaba, ben HaytekUni Cetele Kontrol Botuyum. Daha fazla bilgi için /help yazabilirsin."
		case "help":
			msg.Text = `HaytekUni Cetele Kontrol Botu çetele kontrolü yapar.
Her gün 00:00 da bir hatırlatma mesajı gönderir.
Her gün 01:00 da çetele kontrolü yapar raporun çıktısını gönderir ve kaydeder.
/admins komutu ile adminlerin listesini görebilirsiniz.
/gunlukozet $tarih ile spesifik bir günün kaydını (only admin),
/haftalikozet ile geçen haftanın kaydını görebilirsiniz.(only admin)`

		case "admins":
			for admin, _ := range cfg.Cetele.AllowedUsers {
				msg.Text += "@" + admin + "\n"
			}
		case "test":
			if !isAuthenticated(update.Message.From.UserName) {
				msg.Text = "Bu komutu kullanmak için yetkiniz yok."
				break
			}
			rapor, errs := app.GunlukRaporHazirla("Dün Özet")
			msg.Text = format.GunlukRaporFormat(rapor, errs)
		case "gunlukozet":
			if !isAuthenticated(update.Message.From.UserName) {
				msg.Text = "Bu komutu kullanmak için yetkiniz yok."
				break
			}
			tarih, err := utils.ParseTarihFromCommandArguments(update.Message.CommandArguments())
			if err != nil {
				msg.Text = err.Error()
				break
			}
			res, err := app.GunlukOzetAl(tarih)
			if err != nil {
				msg.Text = "Belirtilen tarihe ait kayıt bulunamadı"
			} else {
				msg.Text = format.GunlukRaporFormat(res, nil)
			}

		case "getid":
			msg.Text = fmt.Sprintf("%s %s kişisinin chatID'si: %d", update.SentFrom().FirstName, update.SentFrom().LastName, update.Message.Chat.ID)
		case "haftalikozet":
			if !isAuthenticated(update.Message.From.UserName) {
				msg.Text = "Bu komutu kullanmak için yetkiniz yok."
				break
			}
			haftalikozet, err := app.HaftalikOzetAl()
			if err != nil {
				msg.Text = "Haftalık özet alınırken hata oluştu"
			}
			msg.Text = format.HaftalikRaporFormat(haftalikozet)

		default:
			msg.Text = "Bu komutu bilmiyorum."
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func SetUpScheduledMessages(bot *tgbotapi.BotAPI) {
	go GunSonuHatirlatmaScheduledMessage(bot)
	go GunSonuKontrolScheduledMessage(bot)
	go GunSonuPersonalScheduledMessage(bot)
}

func GunSonuHatirlatmaScheduledMessage(bot *tgbotapi.BotAPI) {
	loc, _ := time.LoadLocation("Europe/Istanbul")
	cfg = config.Get()
	ticker := time.NewTicker(10 * time.Minute)
	for {
		t := <-ticker.C
		if t.In(loc).Hour() == 0 {
			ticker.Stop()
			txt := "HATIRLATMA: Cetelerinizi 01:00 dan once doldurunuz"
			s := tgbotapi.NewMessage(cfg.Cetele.Gruplar["HaytekUni"], txt)
			_, err := bot.Send(s)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(23 * time.Hour)

			ticker.Reset(10 * time.Minute)
		}
	}
}

func GunSonuKontrolScheduledMessage(bot *tgbotapi.BotAPI) {
	//token saat baiı doluyor. bunun için kontrol etmeden önce servisi restart etmek laızm
	loc, _ := time.LoadLocation("Europe/Istanbul")
	var err error
	cfg = config.Get()
	//ticker := time.NewTicker(time.Second)
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		if t.In(loc).Hour() == 1 {
			ticker.Stop()
			rapor, errs := app.GunlukRaporHazirla("Dün Özet")
			err = rapor.Kaydet()
			if err != nil {
				logx.SendLog(err.Error())
				log.Println(err)
			}
			str := format.GunlukRaporFormat(rapor, errs)

			s := tgbotapi.NewMessage(cfg.Cetele.Gruplar["HaytekUni"], str)
			//s := tgbotapi.NewMessage(cfg.Cetele.Gruplar["mtalhakrc"], str)
			_, err = bot.Send(s)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(23 * time.Hour)
			ticker.Reset(10 * time.Minute)
		}
	}

}
func isAuthenticated(username string) bool {
	return cfg.Cetele.AllowedUsers[username]
}

func GunSonuPersonalScheduledMessage(bot *tgbotapi.BotAPI) {
	loc, _ := time.LoadLocation("Europe/Istanbul")
	ticker := time.NewTicker(10 * time.Minute)
	//ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		if t.In(loc).Hour() == 23 {
			ticker.Stop()
			rapor, errs := app.GunlukRaporHazirla("Bugün Özet")
			for i := 0; i < len(rapor.KisilerSonuc); i++ {
				if cfg.Cetele.Kisiler[rapor.KisilerSonuc[i].Isim] == 0 {
					time.Sleep(2 * time.Second)
					continue
				}
				s := tgbotapi.NewMessage(cfg.Cetele.Kisiler[rapor.KisilerSonuc[i].Isim], format.PersonalRaporFormat(rapor.KisilerSonuc[i], errs[rapor.KisilerSonuc[i].Isim]...))
				_, err := bot.Send(s)
				if err != nil {
					logx.SendLog(err.Error())
					log.Println(err)
				}
				time.Sleep(3 * time.Second)
			}
			time.Sleep(23 * time.Hour)
			ticker.Reset(10 * time.Minute)
		}
	}
}
