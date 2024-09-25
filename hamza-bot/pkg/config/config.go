package config

import "os"

var cfg *Config

type Config struct {
	Database DbConfig
	Bot      BotConfig
	Sheets   SheetsServiceConfig
}

type DbConfig struct {
	Path string
}
type SheetsServiceConfig struct {
	CredentialsPath string
}

type BotConfig struct {
	ID        string
	DebugMode bool
}

func Setup() {
	cfg = &Config{
		Database: DbConfig{
			Path: "/Users/talha/GolandProjects/hamza-bot/hamza-bot.db",
		},
		Bot: BotConfig{
			ID: "5714183726:AAGIc7aijlKxXzV8fTeimJZWUPYOMeAWTRk", //talha test bot
			//ID:        "5688574923:AAETBPEqD26jM4relx3Fu8j7CX7ndYgjDqg",
			DebugMode: false,
		},
		Sheets: SheetsServiceConfig{
			//CredentialsPath: "/Users/talha/go/src/github.com/haytek-uni-bot-yeniden/credentials/fluted-ranger-364116-ea4e986f9ca1.json",
		},
	}

	if os.Getenv("IS_DEVELOPMENT") == "true" {
		cfg.Bot.ID = "5688574923:AAETBPEqD26jM4relx3Fu8j7CX7ndYgjDqg"
		cfg.Bot.DebugMode = true
		cfg.Database.Path = "/home/ubuntu/hamza/hamza-bot.db"
		//cfg.Sheets.CredentialsPath = "/home/ubuntu/credentials/fluted-ranger-364116-ea4e986f9ca1.json"
	}
}

func Get() *Config {
	return cfg
}
