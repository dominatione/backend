package network

import (
	"github.com/dominati-one/backend/internal/pkg/blockchain"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
)

func CreateTestNetGenesisBlock() *blockchainProtocol.Block {
	blockBody := &blockchainProtocol.Block_Body{
		PreviousBlockId: []byte{0x0000000000000000000000000000000000000000000000000000000000000000},
		Timestamp:       blockchain.CreateBlockTimestampFromUnixMilliseconds(1616087057000).UnixMilliseconds(),
		Events:          []*blockchainProtocol.Block_Body_BlockEvent{},
	}

	blockChecksum := []byte{0x0000000000000000000000000000000000000000000000000000000000000000}

	return &blockchainProtocol.Block{
		Body:     blockBody,
		Checksum: blockChecksum,
	}
}
