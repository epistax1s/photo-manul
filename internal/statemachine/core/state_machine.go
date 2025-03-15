package core

import (
	"fmt"

	"github.com/epistax1s/photo-manul/internal/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type State interface {
	Init(*tgbotapi.Update)
	Handle(*tgbotapi.Update)
}

type StateFactory func(*server.Server, *StateMachine, *StateContext) State

type StateHandler func(*tgbotapi.Update)

type StateMachine struct {
	Server       *server.Server
	CurrentState map[int64]State
	StateFactory map[StateType]StateFactory
}

type StateType string

const (
	Idle          StateType = "idle"
	EmployeeID    StateType = "employeeID"
	EmployeePhoto StateType = "employeePhoto"
	InviteAdmin   StateType = "inviteAdmin"
)

type StateContext struct {
	EmployeeID string
	NextState  StateType
}

func (stateMachine *StateMachine) Set(stateType StateType, chatID int64, data *StateContext) State {
	stateFactory, exists := stateMachine.StateFactory[stateType]
	if !exists {
		panic(
			fmt.Sprintf("no state factory was found for a state with type = %s", stateType),
		)
	}

	stateMachine.CurrentState[chatID] = stateFactory(stateMachine.Server, stateMachine, data)
	return stateMachine.CurrentState[chatID]
}

func (stateMachine *StateMachine) Get(chatID int64) State {
	state, exists := stateMachine.CurrentState[chatID]
	if !exists {
		state = stateMachine.Set(Idle, chatID, &StateContext{})
	}
	return state
}
