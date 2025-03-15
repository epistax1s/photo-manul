package builder

import (
	"github.com/epistax1s/photo-manul/internal/server"
	"github.com/epistax1s/photo-manul/internal/statemachine/states/id"
	"github.com/epistax1s/photo-manul/internal/statemachine/states/idle"
	"github.com/epistax1s/photo-manul/internal/statemachine/states/photo"

	. "github.com/epistax1s/photo-manul/internal/statemachine/core"
)

func NewStateMachine(server *server.Server) *StateMachine {
	return &StateMachine{
		Server:       server,
		CurrentState: make(map[int64]State),
		StateFactory: map[StateType]StateFactory{
			Idle:          idle.NewIdleState,
			EmployeeID:    id.NewIDState,
			EmployeePhoto: photo.NewPhotoState,
		},
	}

}
