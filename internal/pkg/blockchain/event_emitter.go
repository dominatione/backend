package blockchain

import (
	"context"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventEmitter struct {
	log         zerolog.Logger
	eventsQueue chan *blockchainProtocol.Event
	blocksQueue chan *blockchainProtocol.Block
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		log:         log.With().Str("applicationComponent", "blockchain").Str("blockchainComponent", "eventEmitter").Logger(),
		eventsQueue: make(chan *blockchainProtocol.Event, 0),
		blocksQueue: make(chan *blockchainProtocol.Block, 0),
	}
}

func (e *EventEmitter) emitEvent(event *blockchainProtocol.Event) error {
	log := e.log.With().Str("eventData", event.Body.String()).Logger()

	eventId, err := NewEventId(event)
	if err != nil {
		return errors.Wrap(err, "unable to calculate event id")
	}
	log = log.With().Str("eventId", eventId.String()).Logger()

	e.eventsQueue <- proto.Clone(event).(*blockchainProtocol.Event)

	log.Trace().Msg("Event added to queue.")

	return nil
}

func (e *EventEmitter) emitBlock(block *blockchainProtocol.Block) error {
	log := e.log.With().Logger()

	blockId, err := NewBlockId(block)
	if err != nil {
		return errors.Wrap(err, "unable to calculate block id")
	}
	log = log.With().Str("blockId", blockId.String()).Logger()

	e.blocksQueue <- proto.Clone(block).(*blockchainProtocol.Block)

	log.Trace().Msg("Block added to queue.")

	return nil
}

func (e *EventEmitter) WaitForEvent(ctx context.Context) (*blockchainProtocol.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ErrCanceledEventEmitterWait
	case event := <-e.eventsQueue:
		return event, nil
	}
}

func (e *EventEmitter) WaitForBlock(ctx context.Context) (*blockchainProtocol.Block, error) {
	select {
	case <-ctx.Done():
		return nil, ErrCanceledEventEmitterWait
	case block := <-e.blocksQueue:
		return block, nil
	}
}

var (
	ErrCanceledEventEmitterWait = errors.New("canceled event emitter wait")
)
