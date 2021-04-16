package blockchain

import (
	"crypto/sha256"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type BlockBuilder struct {
	previousBlock *blockchainProtocol.Block
	events        []*blockchainProtocol.Event
}

func NewBlockBuilder(previousBlock *blockchainProtocol.Block, events []*blockchainProtocol.Event) *BlockBuilder {
	return &BlockBuilder{
		previousBlock: previousBlock,
		events:        events,
	}
}

func (b *BlockBuilder) Build(blockTimestamp BlockTimestamp) (*blockchainProtocol.Block, error) {
	previousBlockId, err := NewBlockId(b.previousBlock)
	if err != nil {
		return nil, errors.Wrap(err, "unable to calculate previous block id")
	}

	blockEvents := []*blockchainProtocol.Block_Body_BlockEvent{}

	for _, event := range b.events {
		eventId, err := NewEventId(event)
		if err != nil {
			return nil, errors.Wrap(err, "unable to calculate event id")
		}

		blockEvents = append(blockEvents, &blockchainProtocol.Block_Body_BlockEvent{
			Id:    eventId.Bytes(),
			Event: event,
		})
	}

	blockBody := &blockchainProtocol.Block_Body{
		PreviousBlockId: previousBlockId.Bytes(),
		Timestamp:       blockTimestamp.UnixMilliseconds(),
		Events:          blockEvents,
	}

	blockBodyBytes, err := proto.Marshal(blockBody)
	if err != nil {
		return nil, ErrInvalidBlockBodyBytes
	}

	blockBodyHash := sha256.Sum256(blockBodyBytes)

	block := &blockchainProtocol.Block{
		Body:     blockBody,
		Checksum: blockBodyHash[:],
	}

	return block, nil
}

var (
	ErrInvalidBlockBodyBytes = errors.New("invalid block body bytes")
)
