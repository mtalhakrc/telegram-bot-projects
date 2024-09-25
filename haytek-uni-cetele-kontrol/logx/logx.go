package logx

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var bot *tgbotapi.BotAPI
var me int64

func InitLogx(botx *tgbotapi.BotAPI) {
	//	log.Printf("Authorized on account %s\n", bot.Self.UserName)
	bot = botx
	me = 952363491
}
func SendLog(text string) {
	text = fmt.Sprintf("%v : %v", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"), text)
	msg := tgbotapi.NewMessage(me, text)
	bot.Send(msg)
}
