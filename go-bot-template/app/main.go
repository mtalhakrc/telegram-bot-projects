package main

import (
	"github.com/go-bot-template/app/handlers"
	"github.com/go-bot-template/pkg/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	ceteleapp := app.New()

	commands := ceteleapp.Commands
	scheduled := ceteleapp.Scheduled

	userhandler := handlers.NewUserHandler(ceteleapp.DB)
	scheduledHandler := handlers.NewScheduled(ceteleapp.DB)

	commands.RegisterCommand("kaydol", userhandler.Kaydol)
	commands.RegisterCommand("isimdegistir", userhandler.UpdateName)
	commands.RegisterCommand("kayitsil", userhandler.DeleteUser)
	commands.RegisterCommand("deneme", userhandler.Deneme)

	scheduled.RegisterScheduled("23:16:00", scheduledHandler.ScheduledDeneme)

	ceteleapp.Start()
}
