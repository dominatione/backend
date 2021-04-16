package local

import (
	"context"
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewConnector(t *testing.T) {
	connector := NewConnector()
	assert.NotNil(t, connector)
}

func TestConnector_SendEventToBacklog(t *testing.T) {
	connector := NewConnector()

	backlogEvent := &blockchain.Event{}

	err := connector.SendEventToBacklog(backlogEvent)
	assert.NoError(t, err)
}

func TestConnector_SendBlockToBacklog(t *testing.T) {
	connector := NewConnector()

	backlogBlock := &blockchain.Block{}

	err := connector.SendBlockToBacklog(backlogBlock)
	assert.NoError(t, err)
}

func TestConnector_GetBacklogBlock(t *testing.T) {
	connector := NewConnector()

	backlogBlock := &blockchain.Block{}

	err := connector.SendBlockToBacklog(backlogBlock)
	assert.NoError(t, err)

	ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*10)

	receivedBacklogBlock, err := connector.GetBacklogBlock(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, receivedBacklogBlock)
	assert.True(t, proto.Equal(backlogBlock, receivedBacklogBlock))
}

func TestConnector_GetBacklogEvent(t *testing.T) {
	connector := NewConnector()

	backlogEvent := &blockchain.Event{}

	err := connector.SendEventToBacklog(backlogEvent)
	assert.NoError(t, err)

	ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*10)

	receivedBacklogEvent, err := connector.GetBacklogEvent(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, receivedBacklogEvent)
	assert.True(t, proto.Equal(backlogEvent, receivedBacklogEvent))
}
