package blockchain

import "github.com/dominati-one/backend/pkg/protocol/blockchain"

type EventValidator struct {
}

func NewEventValidator() *EventValidator {
	return &EventValidator{}
}

func (v *EventValidator) Validate(event *blockchain.Event) error {
	return nil
}
