package event

import (
	"github.com/dominati-one/backend/internal/pkg/game/world"
	"github.com/dominati-one/backend/internal/pkg/security"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/pkg/errors"
)

type CreatePlanetHandler struct {
	state *world.State
}

func NewCreatePlanetHandler(state *world.State) *CreatePlanetHandler {
	return &CreatePlanetHandler{
		state: state,
	}
}

func (h *CreatePlanetHandler) Validate(event *blockchainProtocol.EventCreatePlanet, signature *security.Signature) error {
	stateClone := h.state.Clone()

	planetEntity, err := stateClone.Actions().Planet().Create()
	if err != nil {
		return errors.Wrap(err, "unable to create planet")
	}

	if _, err := stateClone.Actions().Seed().CreateRandomSeedsOnPlanet(*planetEntity); err != nil {
		return errors.Wrap(err, "unable to create random seeds on planet")
	}

	return nil
}

func (h *CreatePlanetHandler) Handle(event *blockchainProtocol.EventCreatePlanet, signature *security.Signature) error {
	if err := h.Validate(event, signature); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	planetEntity, err := h.state.Actions().Planet().Create()
	if err != nil {
		return errors.Wrap(err, "unable to create planet")
	}

	if _, err := h.state.Actions().Seed().CreateRandomSeedsOnPlanet(*planetEntity); err != nil {
		return errors.Wrap(err, "unable to create random seeds on planet")
	}

	return nil
}
