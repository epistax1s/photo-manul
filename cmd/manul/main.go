package main

import (
	"github.com/epistax1s/photo-manul/internal/interceptor"
	"github.com/epistax1s/photo-manul/internal/server"
	"github.com/epistax1s/photo-manul/internal/statemachine/builder"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	server := server.InitServer()

	stateMachine := builder.NewStateMachine(server)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan := server.Manul.GetUpdatesChan(updateConfig)

	chain := interceptor.NewChainBuilder().
		Add(&interceptor.RecoverInterceptor{}).
		Add(&interceptor.LogInterceptor{}).
		Add(&interceptor.RegistrationInterceptor{
			Server: server,
		}).
		Add(&interceptor.SecurityInterceptor{
			Server: server,
		}).
		Add(&interceptor.CancelInterceptor{
			Server:       server,
			StateMachine: stateMachine,
		}).
		Add(&interceptor.HandlerInterceptor{
			Server:       server,
			StateMachine: stateMachine,
		}).
		Build()

	for update := range updateChan {
		chain.Handle(&update)
	}
}
