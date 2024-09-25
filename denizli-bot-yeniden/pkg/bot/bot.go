package bot

import (
	"github.com/go-bot-template/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func New(config config.BotConfig) {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.ID)
	if err != nil {
		panic(err)
	}
	bot.Debug = config.DebugMode
}

func Get() *tgbotapi.BotAPI {
	return bot
}
