package security

import blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"

type Signature struct {
}

func CreateSignatureFromBody(eventBody *blockchainProtocol.Event_Body) (*Signature, error) {
	return NewSignature([]byte{})
}

func NewSignature(buffer []byte) (*Signature, error) {
	return &Signature{}, nil
}

func (s *Signature) Verify(eventBody *blockchainProtocol.Event_Body) error {
	return nil
}
