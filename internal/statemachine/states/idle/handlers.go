package idle

import (
	"fmt"

	. "github.com/epistax1s/photo-manul/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *IdleState) Init(update *tgbotapi.Update) {
	cmdHandler, _ := state.handlers[cmdHelp]
	cmdHandler(update)
}

func (state *IdleState) Handle(update *tgbotapi.Update) {
	manul := state.server.Manul

	chatID := update.FromChat().ID
	cmd := update.Message.Command()

	cmdHandler, exits := state.handlers[cmd]
	if exits {
		cmdHandler(update)
	} else {
		manul.SendMessage(chatID, "Неизвестная команда 🧐")
		state.Init(update)
	}
}

func (state *IdleState) helpHandler(update *tgbotapi.Update) {
	manul := state.server.Manul
	chatID := update.FromChat().ID

	help := fmt.Sprintf("" +
		"/help - вывести список доступных команд 💡\n" +
		"/photo - загрузить фото сотрудника 🖼️\n" +
		"/cancel - сбросить состояние бота, если что-то идет не так 🗘",
	)

	manul.SendMessage(chatID, help)
}

func (state *IdleState) photoHandler(update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	state.stateMachine.
		Set(EmployeeID, chatID, &StateContext{}).
		Init(update)
}
