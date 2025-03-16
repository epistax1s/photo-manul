package interceptor

import (
	"github.com/epistax1s/photo-manul/internal/log"
	"github.com/epistax1s/photo-manul/internal/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SecurityInterceptor struct {
	Server *server.Server
	BaseInterceptor
}

func (i *SecurityInterceptor) Handle(update *tgbotapi.Update) {
	manul := i.Server.Manul
	userService := i.Server.UserService

	chatID := update.FromChat().ID

	_, err := userService.GetUserByChatID(chatID)
	if err != nil {
		log.Error("user access denied", "chatID", chatID, "err", err)

		manul.SendMessage(chatID, "Воспользуйтесь ссылкой-приглашением для получения доступа к боту ⚠️")
		return
	}

	i.Next(update)
}
