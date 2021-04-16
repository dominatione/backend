package component

import (
	"github.com/dominati-one/backend/pkg/protocol/component"
	"github.com/rs/zerolog"
)

type Planet struct {
	Seed int64
	Name string
}

func (p Planet) Protobuf() *component.Planet {
	return &component.Planet{
		Seed: p.Seed,
		Name: p.Name,
	}
}

func (p Planet) MarshalZerologObject(e *zerolog.Event) {
	e.Int64("planetSeed", p.Seed)
	e.Str("planetName", p.Name)
}
