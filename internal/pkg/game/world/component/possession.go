package component

import (
	"github.com/dominati-one/backend/pkg/protocol/component"
	"github.com/rs/zerolog"
)

type Possession struct {
	OwnerEntity Entity
}

func (p Possession) Protobuf() *component.Possession {
	return &component.Possession{
		OwnerEntity: uint64(p.OwnerEntity),
	}
}

func (p Possession) MarshalZerologObject(e *zerolog.Event) {
	e.Str("ownerEntity", p.OwnerEntity.String())
}
