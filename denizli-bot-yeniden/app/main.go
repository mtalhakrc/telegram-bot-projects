package main

import (
	"github.com/go-bot-template/app/handlers"
	"github.com/go-bot-template/app/service"
	"github.com/go-bot-template/pkg/app"
	"github.com/go-bot-template/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	appp := app.New()

	service.InitSheetsService(config.Get().Sheets)
	scheduled := appp.Scheduled

	scheduledHandler := handlers.NewScheduled(service.Get())
	scheduled.RegisterScheduled("12:02:30", scheduledHandler.GunSonuOzetMesaji)

	appp.Start()
}
