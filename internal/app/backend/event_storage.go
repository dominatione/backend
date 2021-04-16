package backend

import (
	"github.com/dominati-one/backend/internal/pkg/blockchain"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
)

type EventStorage struct {
	events map[blockchain.EventId]*blockchainProtocol.Event
}

func NewEventStorage() *EventStorage {
	return &EventStorage{
		events: map[blockchain.EventId]*blockchainProtocol.Event{},
	}
}

func (s *EventStorage) Exists(eventId blockchain.EventId) bool {
	_, exists := s.events[eventId]

	return exists
}
