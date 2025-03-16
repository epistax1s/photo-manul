package interceptor

import (
	"fmt"

	"github.com/epistax1s/photo-manul/internal/log"
	"github.com/epistax1s/photo-manul/internal/model"
	"github.com/epistax1s/photo-manul/internal/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegistrationInterceptor struct {
	Server *server.Server
	BaseInterceptor
}

func (i *RegistrationInterceptor) Handle(update *tgbotapi.Update) {
	i.tryRegister(update)
	i.Next(update)
}

func (i *RegistrationInterceptor) tryRegister(update *tgbotapi.Update) {
	userService := i.Server.UserService

	if update == nil || update.Message == nil || !update.FromChat().IsPrivate() {
		return
	}

	if update.Message.Command() == "start" {
		chatID := update.FromChat().ID
		username := update.FromChat().UserName

		key := update.Message.CommandArguments()

		log.Info("Try register", "chatID", chatID, "username", username, "key", key)

		if key == "" {
			log.Error("Received the start command without a key")
			return
		}

		// Если ключ валидный, создаем пользователя
		if i.isValidKey(key) {
			user := &model.User{
				ChatID:   chatID,
				Username: username,
				Name:     fmt.Sprintf("%s %s", update.FromChat().FirstName, update.FromChat().LastName),
			}
			userService.CreateUser(user)

			log.Info("The user has successfully registered", "user", user, "key", key)
		} else {
			log.Error("Invalid key", "chatID", chatID, "username", username, "key", key)
		}
	}

}

func (i *RegistrationInterceptor) isValidKey(key string) bool {
	validKeys := []string{
		"f7b6b1a5-1188-4b3a-8f1a-35356107c251",
		"6753a735-7aea-41ed-aec2-cf8559af9201",
		"a77b6c82-57fa-4665-b95d-0df9f161f0dc",
		"b398442b-023f-432d-a8fd-527265cf46b7",
	}

	// Проверка, находится ли ключ в массиве
	for _, validKey := range validKeys {
		if key == validKey {
			return true
		}
	}

	return false
}
