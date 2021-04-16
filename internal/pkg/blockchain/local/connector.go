package local

import (
	"context"
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type Connector struct {
	eventBacklog chan *blockchain.Event
	blockBacklog chan *blockchain.Block
}

func NewConnector() *Connector {
	return &Connector{
		eventBacklog: make(chan *blockchain.Event, 2048),
		blockBacklog: make(chan *blockchain.Block, 64),
	}
}

func (c *Connector) SendEventToBacklog(backlogEvent *blockchain.Event) error {
	backlogEventCopy := proto.Clone(backlogEvent).(*blockchain.Event)

	c.eventBacklog <- backlogEventCopy

	return nil
}

func (c *Connector) SendBlockToBacklog(backlogBlock *blockchain.Block) error {
	backlogBlockCopy := proto.Clone(backlogBlock).(*blockchain.Block)

	c.blockBacklog <- backlogBlockCopy

	return nil
}

func (c *Connector) GetBacklogBlock(ctx context.Context) (*blockchain.Block, error) {
	select {
	case <-ctx.Done():
		return nil, ErrCanceledReadBacklogLocalConnector
	case backlogBlock := <-c.blockBacklog:
		return backlogBlock, nil
	}
}

func (c *Connector) GetBacklogEvent(ctx context.Context) (*blockchain.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ErrCanceledReadBacklogLocalConnector
	case backlogEvent := <-c.eventBacklog:
		return backlogEvent, nil
	}
}

var (
	ErrCanceledReadBacklogLocalConnector = errors.New("canceled read backlog local connector")
)
