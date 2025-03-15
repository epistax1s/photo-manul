package idle

import (
	"fmt"

	"github.com/epistax1s/photo-manul/internal/log"
	. "github.com/epistax1s/photo-manul/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *IdleState) Init(update *tgbotapi.Update) {
	helpHandler := state.handlers[cmdHelp]
	helpHandler(update, nil)
}

func (state *IdleState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer

	chatID := update.FromChat().ID
	cmd := update.Message.Command()

	cmdHandler, exits := state.handlers[cmd]
	if exits {
		cmdHandler(update, nil)
	} else {
		gomer.SendMessage(chatID, i18n.Localize("unsupportedCommand"))
		state.Init(update)
	}
}

func (state *IdleState) helpHandler(update *tgbotapi.Update) {
	manul := state.server.Manul
	chatID := update.FromChat().ID

	help := fmt.Sprintf("" +
		"/help		- help description\n" +
		"/photo		- photo description\n" +
		"/invite	- invite description",
	)

	manul.SendMessage(chatID, help)
}

func (state *IdleState) photoHandler(update *tgbotapi.Update) {
	userService := state.server.UserService

	chatID := update.FromChat().ID

	userExists, err := userService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdTrack, "err", err)

		return
	}

	if userExists {
		log.Info(
			"the user is already being tracked",
			"chatID", chatID)

		return
	}

	state.stateMachine.
		Set(TrackDepartment, chatID, &StateContext{}).
		Init(update)
}

func (state *IdleState) inviteHandler(update *tgbotapi.Update) {
	manul := state.server.Gomer
	userService := state.server.UserService

	chatID := update.FromChat().ID

	userExists, err := userService.UserExists(chatID)
	if err != nil {
		log.Error(
			"error checking the user's existence",
			"state", Idle, "cmd", cmdUntrack, "err", err)

		return
	}

	if !userExists {
		log.Info(
			"the user is not a tracked",
			"chatID", chatID)

		return
	}

	if err := userService.UntrackUser(chatID); err != nil {
		log.Error(
			"error when trying to stop tracking the user",
			"state", Idle, "cmd", cmdUntrack, "chatID", chatID, "err", err)

		gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	if err := userService.UntrackUser(chatID); err == nil {
		gomer.SendMessage(chatID, i18n.Localize("trackStopped"))
	} else {
		gomer.SendMessage(chatID, i18n.Localize("oops"))
	}
}
