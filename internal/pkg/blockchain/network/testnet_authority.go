package network

import (
	"crypto"
	"github.com/dominati-one/backend/internal/pkg/security"
)

func CreateTestNetAuthority() *security.PublicKeysBag {
	return security.NewPublicKeysBag([]crypto.PublicKey{})
}
