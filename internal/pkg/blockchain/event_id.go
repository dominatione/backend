package blockchain

import (
	"crypto/sha256"
	"fmt"
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type EventId [32]byte

func MustEventId(event *blockchain.Event) *EventId {
	eventId, err := NewEventId(event)
	if err != nil {
		panic(err)
	}

	return &eventId
}

func NewEventId(event *blockchain.Event) (EventId, error) {
	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return EmptyEventId, ErrEventIdInvalidBytes
	}

	hash := sha256.Sum256(eventBytes)
	eventId := EventId{}

	if len(hash) != len(eventId) {
		return EmptyEventId, ErrEventIdInvalidHashSize
	}

	copy(eventId[:], hash[:])

	return eventId, nil
}

func (e EventId) String() string {
	buffer := [32]byte(e)
	return fmt.Sprintf("%x", buffer)
}

func (b EventId) Bytes() []byte {
	return b[:]
}

var (
	EmptyEventId = EventId{}
)

var (
	ErrEventIdInvalidBytes    = errors.New("event id invalid bytes")
	ErrEventIdInvalidHashSize = errors.New("event id invalid hash size")
)
