package blockchain

import (
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkEventBacklog_Add(t *testing.T) {
	var err error

	eventBacklog := NewNetworkEventBacklog(NewEventValidator())

	err = eventBacklog.Add(&blockchain.Event{})
	assert.NoError(t, err)
}

func TestNetworkEventBacklog_MarkAsConfirmed(t *testing.T) {
	var err error

	eventBacklog := NewNetworkEventBacklog(NewEventValidator())

	event := &blockchain.Event{}
	eventId := MustEventId(event)

	err = eventBacklog.Add(event)
	assert.NoError(t, err)

	unconfirmedEvents := eventBacklog.Unconfirmed()
	assert.Len(t, unconfirmedEvents, 1)

	err = eventBacklog.MarkAsConfirmed(*eventId)
	assert.NoError(t, err)

	unconfirmedEvents = eventBacklog.Unconfirmed()
	assert.Empty(t, unconfirmedEvents)
}
