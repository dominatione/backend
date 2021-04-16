package event

import (
	"github.com/dominati-one/backend/internal/pkg/game/world"
	"github.com/dominati-one/backend/internal/pkg/security"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
)

type CreatePlayerHandler struct {
	state *world.State
}

func NewCreatePlayerHandler(state *world.State) *CreatePlayerHandler {
	return &CreatePlayerHandler{
		state: state,
	}
}

func (h *CreatePlayerHandler) Handle(event *blockchainProtocol.EventCreatePlayer, signature *security.Signature) error {
	return nil
}
