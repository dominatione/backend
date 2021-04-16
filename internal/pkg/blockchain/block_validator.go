package blockchain

import blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"

type BlockValidator struct {
	eventValidator *EventValidator
}

func NewBlockValidator(eventValidator *EventValidator) *BlockValidator {
	return &BlockValidator{
		eventValidator: eventValidator,
	}
}

func (v *BlockValidator) Validate(previousBlock *blockchainProtocol.Block, currentBlock *blockchainProtocol.Block) error {
	return nil
}
