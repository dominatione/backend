package blockchain

import (
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalEventBacklog_Add(t *testing.T) {
	var err error
	var eventId EventId

	eventBacklog := NewLocalEventBacklog(NewEventValidator())

	eventId, err = eventBacklog.Add("invalid event")
	assert.EqualValues(t, EmptyEventId, eventId)
	assert.Error(t, err)

	eventId, err = eventBacklog.Add(&blockchain.EventCreatePlanet{})
	assert.NotEqualValues(t, EmptyEventId, eventId)
	assert.NoError(t, err)

	eventId, err = eventBacklog.Add(&blockchain.EventCreatePlayer{})
	assert.NotEqualValues(t, EmptyEventId, eventId)
	assert.NoError(t, err)
}

func TestLocalEventBacklog_Exists(t *testing.T) {
	var err error
	var eventId EventId

	eventBacklog := NewLocalEventBacklog(NewEventValidator())

	eventId, err = eventBacklog.Add(&blockchain.EventCreatePlanet{})
	assert.NotEqualValues(t, EmptyEventId, eventId)
	assert.NoError(t, err)

	eventExists := eventBacklog.Exists(eventId)
	assert.True(t, eventExists)
}

func TestLocalEventBacklog_MarkAsConfirmed(t *testing.T) {
	var err error
	var eventId EventId

	eventBacklog := NewLocalEventBacklog(NewEventValidator())

	eventId, err = eventBacklog.Add(&blockchain.EventCreatePlanet{})
	assert.NotEqualValues(t, EmptyEventId, eventId)
	assert.NoError(t, err)

	unconfirmedEvents := eventBacklog.Unreceived()
	assert.Len(t, unconfirmedEvents, 1)

	err = eventBacklog.MarkAsReceived(eventId)
	assert.NoError(t, err)

	unconfirmedEvents = eventBacklog.Unreceived()
	assert.Empty(t, unconfirmedEvents)
}

func TestLocalEventBacklog_Unconfirmed(t *testing.T) {
	eventBacklog := NewLocalEventBacklog(NewEventValidator())

	unconfirmedEvents := eventBacklog.Unreceived()
	assert.Empty(t, unconfirmedEvents)

	eventId, err := eventBacklog.Add(&blockchain.EventCreatePlanet{})
	assert.NotEqualValues(t, EmptyEventId, eventId)
	assert.NoError(t, err)

	unconfirmedEvents = eventBacklog.Unreceived()
	assert.Len(t, unconfirmedEvents, 1)
	assert.Contains(t, unconfirmedEvents, eventId)
}

func TestLocalEventBacklog_MarkAsSent(t *testing.T) {
	var err error
	var eventId EventId

	eventBacklog := NewLocalEventBacklog(NewEventValidator())

	eventId, err = eventBacklog.Add(&blockchain.EventCreatePlanet{})
	assert.NotEqualValues(t, EmptyEventId, eventId)
	assert.NoError(t, err)

	unsentEvents := eventBacklog.Unsent()
	assert.Len(t, unsentEvents, 1)

	err = eventBacklog.MarkAsSent(eventId)
	assert.NoError(t, err)

	unsentEvents = eventBacklog.Unsent()
	assert.Empty(t, unsentEvents)
}

func TestLocalEventBacklog_Unsent(t *testing.T) {
	eventBacklog := NewLocalEventBacklog(NewEventValidator())

	unsentEvents := eventBacklog.Unsent()
	assert.Empty(t, unsentEvents)

	eventId, err := eventBacklog.Add(&blockchain.EventCreatePlanet{})
	assert.NotEqualValues(t, EmptyEventId, eventId)
	assert.NoError(t, err)

	unsentEvents = eventBacklog.Unsent()
	assert.Len(t, unsentEvents, 1)
	assert.Contains(t, unsentEvents, eventId)
}
