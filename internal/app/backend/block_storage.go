package backend

import (
	"github.com/dominati-one/backend/internal/pkg/blockchain"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type BlockStorage struct {
	blocks     []*blockchainProtocol.Block
	blocksById map[blockchain.BlockId]*blockchainProtocol.Block
	lastBlock  *blockchainProtocol.Block
}

func NewBlockStorage() *BlockStorage {
	return &BlockStorage{
		blocks:     []*blockchainProtocol.Block{},
		blocksById: map[blockchain.BlockId]*blockchainProtocol.Block{},
	}
}

func (s *BlockStorage) Count() int {
	return len(s.blocks)
}

func (s *BlockStorage) Add(block *blockchainProtocol.Block) (*blockchain.BlockId, error) {
	blockId, err := blockchain.NewBlockId(block)
	if err != nil {
		return nil, errors.Wrap(err, "unable to add block")
	}

	blockCopy := proto.Clone(block).(*blockchainProtocol.Block)

	s.blocks = append(s.blocks, blockCopy)
	s.blocksById[*blockId] = blockCopy

	s.lastBlock = blockCopy

	return blockId, nil
}

func (s *BlockStorage) GetLatestBlock() (*blockchainProtocol.Block, error) {
	if s.lastBlock == nil || len(s.blocks) == 0 {
		return nil, ErrNoBlockInStorage
	}

	return s.lastBlock, nil
}

func (s *BlockStorage) Exists(blockId blockchain.BlockId) bool {
	_, exists := s.blocksById[blockId]
	return exists
}

var (
	ErrNoBlockInStorage = errors.New("no blocks in storage")
)
