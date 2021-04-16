package component

import (
	"fmt"
	"github.com/rs/zerolog"
)

type PlantKind uint8

const (
	PlantKindEmpty PlantKind = iota
	PlantKindOakTree
	PlantKindPineTree
	PlantKindWheat
	PlantKindCannabis
	PlantKindCorn
)

type Plant struct {
	Kind               PlantKind
	Maturity           float32
	AnemochoryMaturity float32
}

func (p PlantKind) String() string {
	switch p {
	case PlantKindEmpty:
		return "PlantKindEmpty"
	case PlantKindOakTree:
		return "PlantKindOakTree"
	case PlantKindPineTree:
		return "PlantKindPineTree"
	case PlantKindWheat:
		return "PlantKindWheat"
	case PlantKindCannabis:
		return "PlantKindCannabis"
	case PlantKindCorn:
		return "PlantKindCorn"
	default:
		panic(fmt.Sprintf("missing PlantKind to string conversion for %d", p))
	}
}

func (p Plant) MarshalZerologObject(e *zerolog.Event) {
	e.Str("plantKind", p.Kind.String())
	e.Float32("plantMaturity", p.Maturity)
	e.Float32("plantAnemochoryMaturity", p.AnemochoryMaturity)
}
