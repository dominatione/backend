package blockchain

import (
	"crypto/sha256"
	"fmt"
	"github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type BlockId [32]byte

func NewBlockId(block *blockchain.Block) (*BlockId, error) {
	blockBytes, err := proto.Marshal(block)
	if err != nil {
		return nil, ErrBlockIdInvalidBytes
	}

	hash := sha256.Sum256(blockBytes)
	blockId := BlockId{}

	if len(hash) != len(blockId) {
		return nil, ErrBlockIdInvalidHashSize
	}

	copy(blockId[:], hash[:])

	return &blockId, nil
}

func (b BlockId) Bytes() []byte {
	return b[:]
}

func (e BlockId) String() string {
	buffer := [32]byte(e)
	return fmt.Sprintf("%x", buffer)
}

var (
	ErrBlockIdInvalidBytes    = errors.New("block id invalid bytes")
	ErrBlockIdInvalidHashSize = errors.New("block id invalid hash size")
)
