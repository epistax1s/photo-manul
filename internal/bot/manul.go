package manul

import (
	"github.com/epistax1s/photo-manul/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Manul struct {
	tgbotapi.BotAPI
}

func InitTelegramBot(botConfig *config.BotConfig) (*Manul, error) {
	botApi, err := tgbotapi.NewBotAPI(botConfig.Token)
	if err != nil {
		return nil, err
	}
	return &Manul{BotAPI: *botApi}, nil
}
