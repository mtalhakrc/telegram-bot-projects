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
			Path: "/Users/talha/go/src/github.com/haytek-uni-bot-yeniden/pkg/database/deneme.db",
		},
		Bot: BotConfig{
			ID:        "5714183726:AAGIc7aijlKxXzV8fTeimJZWUPYOMeAWTRk", //talha test bot
			DebugMode: false,
		},
		Sheets: SheetsServiceConfig{
			CredentialsPath: "/Users/talha/go/src/denizli-bot-yeniden/denizli-cetele-kontrol-c3ed78490729.json",
		},
	}

	if os.Getenv("IS_DEVELOPMENT") == "true" {
		cfg.Bot.ID = "5556282721:AAGUUDk8ivh1DHoXWFCZmyk9wtHYtEGUh68" //denizli bot
		cfg.Bot.DebugMode = true
		cfg.Sheets.CredentialsPath = "/home/ubuntu/denizli/credentials/denizli-cetele-kontrol-c3ed78490729.json"
	}
}

func Get() *Config {
	return cfg
}
