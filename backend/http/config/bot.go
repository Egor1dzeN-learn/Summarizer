package config

import (
	"os"
)

type BotConfig struct {
	Token string
}

func LoadBotConfig() *BotConfig {
	return &BotConfig{
		Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}
