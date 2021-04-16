package component

import (
	"fmt"
	"github.com/dominati-one/backend/pkg/protocol/component"
	"github.com/rs/zerolog"
)

type SeedKind uint8

const (
	SeedKindEmpty SeedKind = iota
	SeedKindOakTree
	SeedKindPineTree
	SeedKindWheat
	SeedKindCannabis
	SeedKindCorn
)

type Seed struct {
	Kind     SeedKind
	Maturity float32
}

func (s SeedKind) String() string {
	switch s {
	case SeedKindOakTree:
		return "SeedKindOakTree"
	case SeedKindPineTree:
		return "SeedKindPineTree"
	case SeedKindWheat:
		return "SeedKindWheat"
	case SeedKindCannabis:
		return "SeedKindCanabis"
	case SeedKindCorn:
		return "SeedKindCorn"
	default:
		panic(fmt.Sprintf("missing SeedKind to string conversion for %d", s))
	}
}

func (s SeedKind) Protobuf() component.SeedKind {
	switch s {
	case SeedKindEmpty:
		return component.SeedKind_SEED_KIND_EMPTY
	case SeedKindOakTree:
		return component.SeedKind_SEED_KIND_OAK_TREE
	case SeedKindPineTree:
		return component.SeedKind_SEED_KIND_PINE_TREE
	case SeedKindWheat:
		return component.SeedKind_SEED_KIND_WHEAT
	case SeedKindCannabis:
		return component.SeedKind_SEED_KIND_CANNABIS
	case SeedKindCorn:
		return component.SeedKind_SEED_KIND_CORN
	default:
		panic(fmt.Sprintf("missing SeedKind to component conversion for %d", s))
	}
}

func (s Seed) MarshalZerologObject(e *zerolog.Event) {
	e.Str("seedKind", s.Kind.String())
	e.Float32("seedMaturity", s.Maturity)
}

func (s Seed) Protobuf() *component.Seed {
	return &component.Seed{
		Kind:     s.Kind.Protobuf(),
		Maturity: s.Maturity,
	}
}
