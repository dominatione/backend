package blockchain

import (
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type localEventBacklogItem struct {
	event     *blockchainProtocol.Event
	received  bool
	sent      bool
	confirmed bool
}

type LocalEventBacklog struct {
	log            zerolog.Logger
	eventValidator *EventValidator

	state  sync.Mutex
	events map[EventId]*localEventBacklogItem
}

func NewLocalEventBacklog(eventValidator *EventValidator) *LocalEventBacklog {
	return &LocalEventBacklog{
		log:            log.With().Str("applicationComponent", "blockchain").Str("blockchainComponent", "localEventBacklog").Logger(),
		eventValidator: eventValidator,
		events:         map[EventId]*localEventBacklogItem{},
	}
}

// Exists checks if event is in local backlog.
func (b *LocalEventBacklog) Exists(eventId EventId) bool {
	defer b.state.Unlock()
	b.state.Lock()

	_, exists := b.events[eventId]

	return exists
}

// MarkAsReceived should be called, when local event is received from network. This means, that event was validated by nodes
// and will be placed on block.
func (b *LocalEventBacklog) MarkAsReceived(eventId EventId) error {
	defer b.state.Unlock()
	b.state.Lock()

	if localEvent, exists := b.events[eventId]; !exists {
		return errors.Wrap(ErrLocalBacklogEventNotFound, "unable to mark local event as received")
	} else {
		b.log.Trace().
			Str("eventId", eventId.String()).
			Str("eventData", localEvent.event.Body.String()).
			Msg("Local event marked as received.")

		localEvent.received = true
	}

	return nil
}

// MarkAsConfirmed should be called, when local event is received from block. This means, that event was placed in block and
// confirmed.
func (b *LocalEventBacklog) MarkAsConfirmed(eventId EventId) error {
	defer b.state.Unlock()
	b.state.Lock()

	if localEvent, exists := b.events[eventId]; !exists {
		return errors.Wrap(ErrLocalBacklogEventNotFound, "unable to mark local event as confirmed")
	} else {
		b.log.Trace().
			Str("eventId", eventId.String()).
			Str("eventData", localEvent.event.Body.String()).
			Msg("Local event marked as confirmed.")

		localEvent.confirmed = true
	}

	return nil
}

// MarkAsSent should be called, when local event is sent to network.
func (b *LocalEventBacklog) MarkAsSent(eventId EventId) error {
	defer b.state.Unlock()
	b.state.Lock()

	if localEvent, exists := b.events[eventId]; !exists {
		return errors.Wrap(ErrLocalBacklogEventNotFound, "unable to mark local event as sent")
	} else {
		b.log.Trace().
			Str("eventId", eventId.String()).
			Str("eventData", localEvent.event.String()).
			Msg("Local event marked as send.")

		localEvent.sent = true
	}

	return nil
}

// Unreceived returns map with unconfirmed events. Events are copy of original backlog item.
func (b *LocalEventBacklog) Unreceived() map[EventId]*blockchainProtocol.Event {
	defer b.state.Unlock()
	b.state.Lock()

	unreceived := map[EventId]*blockchainProtocol.Event{}

	for eventId, event := range b.events {
		if event.received {
			continue
		}

		unreceived[eventId] = proto.Clone(event.event).(*blockchainProtocol.Event)
	}

	return unreceived
}

// Unconfirmed returns map with unconfirmed events. Events are copy of original backlog item.
func (b *LocalEventBacklog) Unconfirmed() map[EventId]*blockchainProtocol.Event {
	defer b.state.Unlock()
	b.state.Lock()

	unconfirmed := map[EventId]*blockchainProtocol.Event{}

	for eventId, event := range b.events {
		if event.confirmed {
			continue
		}

		unconfirmed[eventId] = proto.Clone(event.event).(*blockchainProtocol.Event)
	}

	return unconfirmed
}

// Unsent returns map with unsent events. Events are copy of original backlog item.
func (b *LocalEventBacklog) Unsent() map[EventId]*blockchainProtocol.Event {
	defer b.state.Unlock()
	b.state.Lock()

	unsent := map[EventId]*blockchainProtocol.Event{}

	for eventId, event := range b.events {
		if event.sent {
			continue
		}

		unsent[eventId] = proto.Clone(event.event).(*blockchainProtocol.Event)
	}

	return unsent
}

// Add insert local emitted event in to backlog.
func (b *LocalEventBacklog) Add(localEvent interface{}) (EventId, error) {
	defer b.state.Unlock()
	b.state.Lock()

	log := b.log.With().Logger()

	backlogEvent := &blockchainProtocol.Event{
		Body: &blockchainProtocol.Event_Body{},
	}

	switch resolvedEvent := localEvent.(type) {
	case *blockchainProtocol.EventCreatePlanet:
		backlogEvent.Body.Event = &blockchainProtocol.Event_Body_CreatePlanet{CreatePlanet: resolvedEvent}
	case *blockchainProtocol.EventCreatePlayer:
		backlogEvent.Body.Event = &blockchainProtocol.Event_Body_CreatePlayer{CreatePlayer: resolvedEvent}
	default:
		return EmptyEventId, ErrLocalBacklogUnsupportedEvent
	}

	backlogEvent.Timestamp = CreateBlockTimestampFromNow().UnixMilliseconds()

	log = log.With().Str("eventData", backlogEvent.String()).Logger()

	err := b.eventValidator.Validate(backlogEvent)
	if err != nil {
		log.Warn().Err(err)
		return EmptyEventId, errors.Wrap(err, "unable to add event to local backlog")
	}

	eventId, err := NewEventId(backlogEvent)
	if err != nil {
		log.Warn().Err(err)
		return EmptyEventId, errors.Wrap(err, "unable to add event to local backlog")
	}

	if _, exists := b.events[eventId]; exists {
		return EmptyEventId, errors.Wrap(ErrLocalBacklogEventAlreadyExists, "unable to add event to local backlog")
	}

	b.events[eventId] = &localEventBacklogItem{
		event:     backlogEvent,
		received:  false,
		sent:      false,
		confirmed: false,
	}

	log.Debug().Str("eventId", eventId.String()).Msg("Event added to local backlog.")

	return eventId, nil
}

var (
	ErrLocalBacklogUnsupportedEvent   = errors.New("local backlog unsupported event")
	ErrLocalBacklogEventNotFound      = errors.New("local backlog event not found")
	ErrLocalBacklogEventAlreadyExists = errors.New("local backlog event already exists")
)
