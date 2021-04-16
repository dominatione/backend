package blockchain

import (
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type networkEventBacklogItem struct {
	event     *blockchain.Event
	confirmed bool
}

type NetworkEventBacklog struct {
	log            zerolog.Logger
	eventValidator *EventValidator
	state          sync.Mutex
	events         map[EventId]*networkEventBacklogItem
}

func NewNetworkEventBacklog(eventValidator *EventValidator) *NetworkEventBacklog {
	return &NetworkEventBacklog{
		log:            log.With().Str("applicationComponent", "blockchain").Str("blockchainComponent", "networkEventBacklog").Logger(),
		eventValidator: eventValidator,
		events:         map[EventId]*networkEventBacklogItem{},
	}
}

// Exists checks if event is in network backlog.
func (b *NetworkEventBacklog) Exists(eventId EventId) bool {
	defer b.state.Unlock()
	b.state.Lock()

	_, exists := b.events[eventId]

	return exists
}

func (b *NetworkEventBacklog) MarkAsConfirmed(eventId EventId) error {
	defer b.state.Unlock()
	b.state.Lock()

	if networkEvent, exists := b.events[eventId]; !exists {
		return errors.Wrap(ErrLocalBacklogEventNotFound, "unable to mark network event as confirmed")
	} else {
		b.log.Trace().
			Str("eventId", eventId.String()).
			Str("eventData", networkEvent.event.Body.String()).
			Msg("Network event marked as confirmed.")

		networkEvent.confirmed = true
	}

	return nil
}

// Unconfirmed returns map with unconfirmed events. Events are copy of original backlog item.
func (b *NetworkEventBacklog) Unconfirmed() map[EventId]*blockchain.Event {
	defer b.state.Unlock()
	b.state.Lock()

	unconfirmed := map[EventId]*blockchain.Event{}

	for eventId, event := range b.events {
		if event.confirmed {
			continue
		}

		unconfirmed[eventId] = proto.Clone(event.event).(*blockchain.Event)
	}

	return unconfirmed
}

func (b *NetworkEventBacklog) All() map[EventId]*blockchain.Event {
	defer b.state.Unlock()
	b.state.Lock()

	allEvents := map[EventId]*blockchain.Event{}

	for eventId, event := range b.events {
		allEvents[eventId] = event.event
	}

	return allEvents
}

func (b *NetworkEventBacklog) Add(networkEvent *blockchain.Event) error {
	defer b.state.Unlock()
	b.state.Lock()

	log := b.log.With().Str("eventData", networkEvent.String()).Logger()

	err := b.eventValidator.Validate(networkEvent)
	if err != nil {
		log.Err(err)
		return errors.Wrap(err, "unable to add event to network backlog")
	}

	eventId, err := NewEventId(networkEvent)
	if err != nil {
		log.Err(err)
		return errors.Wrap(err, "unable to add event to network backlog")
	}
	log = b.log.With().Str("eventId", eventId.String()).Logger()

	if _, exists := b.events[eventId]; exists {
		return errors.Wrap(ErrNetworkBacklogEventAlreadyExists, "unable to add event to network backlog")
	}

	b.events[eventId] = &networkEventBacklogItem{
		event:     networkEvent,
		confirmed: false,
	}

	log.Debug().Str("eventId", eventId.String()).Msg("Event added to network backlog.")

	return nil
}

var (
	ErrNetworkBacklogEventAlreadyExists = errors.New("network backlog event already exists")
)
