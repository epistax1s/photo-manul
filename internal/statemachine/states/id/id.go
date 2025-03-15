package id

import (
	"github.com/epistax1s/photo-manul/internal/server"
	. "github.com/epistax1s/photo-manul/internal/statemachine/core"
)

type IDState struct {
	server       *server.Server
	stateMachine *StateMachine
}

func NewIDState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &IDState{
		server:       server,
		stateMachine: stateMachine,
	}

	return state
}
