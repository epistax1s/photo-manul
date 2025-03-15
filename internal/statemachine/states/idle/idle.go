package idle

import (
	"github.com/epistax1s/photo-manul/internal/server"
	. "github.com/epistax1s/photo-manul/internal/statemachine/core"
)

const (
	cmdStart  = "start"
	cmdHelp   = "help"
	cmdPhoto  = "photo"
	cmdInvite = "invite"
)

type IdleState struct {
	server       *server.Server
	stateMachine *StateMachine
	handlers     map[string]StateHandler
}

func NewIdleState(server *server.Server, stateMachine *StateMachine, data *StateContext) State {
	state := &IdleState{
		server:       server,
		stateMachine: stateMachine,
	}

	state.handlers = map[string]StateHandler{
		cmdStart:  state.helpHandler,
		cmdHelp:   state.helpHandler,
		cmdPhoto:  state.photoHandler,
		cmdInvite: state.inviteHandler,
	}

	return state
}
