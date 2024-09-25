package main

import (
	"github.com/haytek-uni-bot-yeniden/app/handlers"
	"github.com/haytek-uni-bot-yeniden/app/service"
	"github.com/haytek-uni-bot-yeniden/pkg/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	ceteleapp := app.New()

	s := service.NewSheetsService(ceteleapp.Sheets, "1ZC1mad4ERjt9Ze2nS8Y1noeHiGx5VKA0wxyryIz8bRk")

	commands := ceteleapp.Commands
	scheduled := ceteleapp.Scheduled

	userhandler := handlers.NewUserHandler(ceteleapp.DB, s)
	scheduledHandler := handlers.NewScheduled(ceteleapp.DB, s)
	cetelehandler := handlers.NewCeteleHandler(ceteleapp.DB)

	commands.RegisterCommand("kaydol", userhandler.Kaydol)
	commands.RegisterCommand("kayitsil", userhandler.KayitSil)
	commands.RegisterCommand("makeadmin", userhandler.MakeAdmin)
	commands.RegisterCommand("chatid", cetelehandler.GetChatID)

	commands.RegisterCommand("start", cetelehandler.Start)
	commands.RegisterCommand("admins", cetelehandler.Admins)

	commands.RegisterCommand("gunlukozet", cetelehandler.GetSpecificRecord)
	commands.RegisterCommand("haftalikozet", cetelehandler.GetHaftalikOzet)

	scheduled.RegisterScheduled("22:00:10", scheduledHandler.CeteleHatirlatmaMesaji)
	scheduled.RegisterScheduled("23:00:10", scheduledHandler.GunlukErkenKontrolMesaji)
	scheduled.RegisterScheduled("01:00:10", scheduledHandler.GunlukRaporMesaji)

	ceteleapp.Start()
}
