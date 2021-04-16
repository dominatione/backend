package blockchain

import (
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type localBlockBacklogItem struct {
	block    *blockchainProtocol.Block
	sent     bool
	received bool
}

type LocalBlockBacklog struct {
	log            zerolog.Logger
	blockValidator *BlockValidator

	blocks      map[BlockId]*localBlockBacklogItem
	state       sync.Mutex
	latestBlock *localBlockBacklogItem
}

func NewLocalBlockBacklog(blockValidator *BlockValidator) *LocalBlockBacklog {
	return &LocalBlockBacklog{
		log:            log.With().Str("applicationComponent", "blockchain").Str("blockchainComponent", "localBlockBacklog").Logger(),
		blockValidator: blockValidator,
		blocks:         map[BlockId]*localBlockBacklogItem{},
	}
}

func (b *LocalBlockBacklog) All() map[BlockId]*blockchainProtocol.Block {
	defer b.state.Unlock()
	b.state.Lock()

	allBlocks := map[BlockId]*blockchainProtocol.Block{}

	for blockId, block := range b.blocks {
		blockCopy := proto.Clone(block.block).(*blockchainProtocol.Block)
		allBlocks[blockId] = blockCopy
	}

	return allBlocks
}

func (b *LocalBlockBacklog) Unreceived() map[BlockId]*blockchainProtocol.Block {
	defer b.state.Unlock()
	b.state.Lock()

	confirmedBlocks := map[BlockId]*blockchainProtocol.Block{}

	for blockId, block := range b.blocks {
		if block.received {
			continue
		}

		blockCopy := proto.Clone(block.block).(*blockchainProtocol.Block)
		confirmedBlocks[blockId] = blockCopy
	}

	return confirmedBlocks
}

func (b *LocalBlockBacklog) Unsent() map[BlockId]*blockchainProtocol.Block {
	defer b.state.Unlock()
	b.state.Lock()

	unsentBlocks := map[BlockId]*blockchainProtocol.Block{}

	for blockId, block := range b.blocks {
		if block.sent {
			continue
		}

		blockCopy := proto.Clone(block.block).(*blockchainProtocol.Block)
		unsentBlocks[blockId] = blockCopy
	}

	return unsentBlocks
}

func (b *LocalBlockBacklog) MarkAsSent(blockId BlockId) error {
	defer b.state.Unlock()
	b.state.Lock()

	if localBlock, exists := b.blocks[blockId]; !exists {
		return errors.Wrap(ErrLocalBacklogBlockNotFound, "unable to mark local block as sent")
	} else {
		b.log.Trace().
			Str("blockId", blockId.String()).
			Msg("Local block marked as send.")

		localBlock.sent = true
	}

	return nil
}

func (b *LocalBlockBacklog) MarkAsConfirmed(blockId BlockId) error {
	defer b.state.Unlock()
	b.state.Lock()

	if localBlock, exists := b.blocks[blockId]; !exists {
		return errors.Wrap(ErrLocalBacklogBlockNotFound, "unable to mark local block as received")
	} else {
		b.log.Trace().
			Str("blockId", blockId.String()).
			Msg("Local block marked as received.")

		localBlock.received = true
	}

	return nil
}

func (b *LocalBlockBacklog) Exists(blockId BlockId) bool {
	defer b.state.Unlock()
	b.state.Lock()

	_, exists := b.blocks[blockId]
	return exists
}

func (b *LocalBlockBacklog) Add(localBlock blockchainProtocol.Block) (*BlockId, error) {
	defer b.state.Unlock()
	b.state.Lock()

	blockId, err := NewBlockId(&localBlock)
	if err != nil {
		log.Warn().Err(err)
		return nil, errors.Wrap(err, "unable to generate block id")
	}

	if _, exists := b.blocks[*blockId]; exists {
		log.Warn().Err(ErrLocalBacklogBlockAlreadyExists)
		return nil, ErrLocalBacklogBlockAlreadyExists
	}

	blockCopy := proto.Clone(&localBlock).(*blockchainProtocol.Block)

	b.blocks[*blockId] = &localBlockBacklogItem{
		block:    blockCopy,
		received: false,
		sent:     false,
	}

	log.Debug().Str("blockId", blockId.String()).Msg("Block added to local backlog.")

	return blockId, nil
}

var (
	ErrLocalBacklogBlockAlreadyExists = errors.New("local backlog block already exists")
	ErrLocalBacklogBlockNotFound      = errors.New("local backlog block not found")
)
