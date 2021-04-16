package security

import "crypto"

type PublicKeysBag struct {
	keys []crypto.PublicKey
}

func NewPublicKeysBag(keys []crypto.PublicKey) *PublicKeysBag {
	return &PublicKeysBag{
		keys: keys,
	}
}

func (b *PublicKeysBag) VerifySignature(signature Signature) bool {
	return false
}
