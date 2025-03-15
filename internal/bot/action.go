package manul

import (
	"github.com/epistax1s/photo-manul/internal/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (gomer *Manul) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	gomer.Send(msg)
}

func (gomer *Manul) SendMessageWithKeyboard(chatID int64, text string, markup *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	gomer.Send(msg)
}

func (gomer *Manul) EditMessageWithKeyboard(chatID int64, messageID int, text string, markup *tgbotapi.InlineKeyboardMarkup) {
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ReplyMarkup = markup
	gomer.Send(editMsg)
}

func (gomer *Manul) SendCallbackResponse(callback *tgbotapi.CallbackQuery, text string) error {
	callbackConfig := tgbotapi.CallbackConfig{
		CallbackQueryID: callback.ID,
		Text:            text,
		ShowAlert:       false,
	}

	_, err := gomer.Request(callbackConfig)
	return err
}

func (gomer *Manul) RemoveMarkup(callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, callback.Message.Text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	// TODO fix log
	if _, err := gomer.Send(msg); err != nil {
		log.Error(
			"error",
			"err", err)
	}
}
