package component

import (
	"fmt"
	"github.com/rs/zerolog"
)

type EntityKind uint8

const (
	EntityKindUnknown EntityKind = iota
	EntityKindPlanet
	EntityKindPlayer
	EntityKindSeedOakTree
	EntityKindSeedPineTree
	EntityKindSeedWheat
	EntityKindSeedCannabis
	EntityKindSeedCorn
	EntityKindPlantOakTree
	EntityKindPlantPineTree
	EntityKindPlantWheat
	EntityKindPlantCannabis
	EntityKindPlantCorn
)

type Entity uint64

func (entity Entity) MarshalZerologObject(e *zerolog.Event) {
	e.Uint64("component", uint64(entity))
}

func (entity Entity) String() string {
	return fmt.Sprintf("%d", entity)
}
