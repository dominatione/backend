package blockchain

import blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"

type BlockStorage interface {
	Count() int
	Add(block *blockchainProtocol.Block) (*BlockId, error)
	GetLatestBlock() (*blockchainProtocol.Block, error)
	Exists(blockId BlockId) bool
}
