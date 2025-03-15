package interceptor

import (
	"github.com/epistax1s/photo-manul/internal/log"
	"github.com/epistax1s/photo-manul/internal/server"
	"github.com/epistax1s/photo-manul/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerInterceptor struct {
	Server       *server.Server
	StateMachine *core.StateMachine
	BaseInterceptor
}

func (i *HandlerInterceptor) Handle(update *tgbotapi.Update) {
	if update == nil {
		return
	}

	if update.FromChat().IsPrivate() {
		chatID := update.FromChat().ID
		i.StateMachine.
			Get(chatID).
			Handle(update)
	} else {
		log.Error("This bot can only be used in private conversations and in groups")
	}
}
