package photo

import (
	"github.com/epistax1s/photo-manul/internal/server"
	. "github.com/epistax1s/photo-manul/internal/statemachine/core"
)

type PhotoState struct {
	server       *server.Server
	stateMachine *StateMachine
	context      *StateContext
}

func NewPhotoState(server *server.Server, stateMachine *StateMachine, context *StateContext) State {
	state := &PhotoState{
		server:       server,
		stateMachine: stateMachine,
		context:      context,
	}

	return state
}
