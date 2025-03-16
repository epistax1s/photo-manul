package id

import (
	"fmt"
	"strconv"

	"github.com/epistax1s/photo-manul/internal/log"
	. "github.com/epistax1s/photo-manul/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *IDState) Init(update *tgbotapi.Update) {
	manul := state.server.Manul
	chatID := update.FromChat().ID

	manul.SendMessage(chatID, "Введите ID сотрудника 🪪")
}

func (state *IDState) Handle(update *tgbotapi.Update) {
	manul := state.server.Manul
	employeeService := state.server.EmployeeService

	chatID := update.FromChat().ID

	if update.Message == nil {
		log.Error(
			"Message is nil",
			"chatID", chatID, "state", EmployeeID, "step", "Handle")

		state.Init(update)
		return
	}

	text := update.Message.Text

	employeeID, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		log.Error(
			"The employeeID has an incorrect format",
			"chatID", chatID, "state", EmployeeID, "step", "Handle", "Message.Text", text)

		manul.SendMessage(chatID, "ID имеет некорректный формат ✖")

		state.Init(update)
		return
	}

	employee, err := employeeService.GetEmployeeByID(employeeID)
	if err != nil {
		log.Error(
			"The employee with ID was not found",
			"chatID", chatID, "state", EmployeeID, "step", "Handle", "err", err)

		manul.SendMessage(chatID, fmt.Sprintf("Сотрудник с ID %d не найден 🔎", employeeID))
		state.Init(update)
		return
	}

	state.stateMachine.
		Set(EmployeePhoto, chatID, &StateContext{Employee: employee}).
		Init(update)
}
