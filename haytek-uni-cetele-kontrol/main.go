package main

import (
	"haytekuni-cetele-kontrol/bot"
	"haytekuni-cetele-kontrol/config"
	"haytekuni-cetele-kontrol/database"
	"haytekuni-cetele-kontrol/logx"
	"haytekuni-cetele-kontrol/service"
)

func main() {
	config.SetupConfig()
	cfg := config.Get()
	telegrambot := bot.InitBot()
	logx.InitLogx(telegrambot)
	service.InitSheetService()
	database.New(cfg.Database)
	bot.SetUpScheduledMessages(telegrambot)
	bot.ListenForUpdates(telegrambot)
}
