package blockchain

import (
	"context"
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
)

type Connector interface {
	// SendEventToBacklog send event to network backlog
	SendEventToBacklog(backlogEvent *blockchain.Event) error

	// SendBlockToBacklog send block proposal to network backlog
	SendBlockToBacklog(backlogBlock *blockchain.Block) error

	GetBacklogBlock(ctx context.Context) (*blockchain.Block, error)
	GetBacklogEvent(ctx context.Context) (*blockchain.Event, error)
}
