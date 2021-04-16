package blockchain

import (
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalBlockBacklog_Add(t *testing.T) {
	var err error
	var blockId *BlockId

	blockBacklog := NewLocalBlockBacklog(NewBlockValidator(NewEventValidator()))

	blockId, err = blockBacklog.Add(blockchain.Block{})
	assert.NotNil(t, blockId)
	assert.NoError(t, err)
}

func TestLocalBlockBacklog_Exists(t *testing.T) {
	var err error
	var blockId *BlockId

	blockBacklog := NewLocalBlockBacklog(NewBlockValidator(NewEventValidator()))

	blockId, err = blockBacklog.Add(blockchain.Block{})
	assert.NotNil(t, blockId)
	assert.NoError(t, err)

	eventExists := blockBacklog.Exists(*blockId)
	assert.True(t, eventExists)
}

func TestLocalBlockBacklog_MarkAsConfirmed(t *testing.T) {
	var err error
	var blockId *BlockId

	blockBacklog := NewLocalBlockBacklog(NewBlockValidator(NewEventValidator()))

	blockId, err = blockBacklog.Add(blockchain.Block{})
	assert.NotNil(t, blockId)
	assert.NoError(t, err)

	unconfirmedEvents := blockBacklog.Unreceived()
	assert.Len(t, unconfirmedEvents, 1)

	err = blockBacklog.MarkAsConfirmed(*blockId)
	assert.NoError(t, err)

	unconfirmedEvents = blockBacklog.Unreceived()
	assert.Empty(t, unconfirmedEvents)
}

func TestLocalBlockBacklog_Unconfirmed(t *testing.T) {
	blockBacklog := NewLocalBlockBacklog(NewBlockValidator(NewEventValidator()))

	unconfirmedEvents := blockBacklog.Unreceived()
	assert.Empty(t, unconfirmedEvents)

	blockId, err := blockBacklog.Add(blockchain.Block{})
	assert.NotNil(t, blockId)
	assert.NoError(t, err)

	unconfirmedEvents = blockBacklog.Unreceived()
	assert.Len(t, unconfirmedEvents, 1)
	assert.Contains(t, unconfirmedEvents, *blockId)
}

func TestLocalBlockBacklog_MarkAsSent(t *testing.T) {
	var err error
	var blockId *BlockId

	blockBacklog := NewLocalBlockBacklog(NewBlockValidator(NewEventValidator()))

	blockId, err = blockBacklog.Add(blockchain.Block{})
	assert.NotNil(t, blockId)
	assert.NoError(t, err)

	unsentEvents := blockBacklog.Unsent()
	assert.Len(t, unsentEvents, 1)

	err = blockBacklog.MarkAsSent(*blockId)
	assert.NoError(t, err)

	unsentEvents = blockBacklog.Unsent()
	assert.Empty(t, unsentEvents)
}

func TestLocalBlockBacklog_Unsent(t *testing.T) {
	blockBacklog := NewLocalBlockBacklog(NewBlockValidator(NewEventValidator()))

	unsentBlocks := blockBacklog.Unsent()
	assert.Empty(t, unsentBlocks)

	blockId, err := blockBacklog.Add(blockchain.Block{})
	assert.NotNil(t, blockId)
	assert.NoError(t, err)

	unsentBlocks = blockBacklog.Unsent()
	assert.Len(t, unsentBlocks, 1)
	assert.Contains(t, unsentBlocks, *blockId)
}
