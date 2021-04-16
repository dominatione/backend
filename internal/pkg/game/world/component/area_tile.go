package component

import (
	"github.com/dominati-one/backend/pkg/protocol/component"
)

type AreaTileKind uint8

const (
	AreaTileKindEmpty AreaTileKind = iota
	AreaTileKindWater
	AreaTileKindShallowWater
	AreaTileKindSand
	AreaTileKindGround
	AreaTileKindFertileGround
	AreaTileKindGravel
	AreaTileKindLava
	AreaTileKindStone
	AreaTileKindSnow
)

type AreaTile struct {
	Kind        AreaTileKind // 1 byte
	OwnerEntity Entity       // 8 bytes
}

type AreaTiles []AreaTile

func (k AreaTileKind) Protobuf() component.AreaTileKind {
	switch k {
	case AreaTileKindEmpty:
		return component.AreaTileKind_AREA_TILE_EMPTY
	case AreaTileKindWater:
		return component.AreaTileKind_AREA_TILE_WATER
	case AreaTileKindShallowWater:
		return component.AreaTileKind_AREA_TILE_SHALLOW_WATER
	case AreaTileKindSand:
		return component.AreaTileKind_AREA_TILE_SAND
	case AreaTileKindGround:
		return component.AreaTileKind_AREA_TILE_GROUND
	case AreaTileKindFertileGround:
		return component.AreaTileKind_AREA_TILE_FERTILE_GROUND
	case AreaTileKindGravel:
		return component.AreaTileKind_AREA_TILE_GRAVEL
	case AreaTileKindLava:
		return component.AreaTileKind_AREA_TILE_LAVA
	case AreaTileKindStone:
		return component.AreaTileKind_AREA_TILE_STONE
	case AreaTileKindSnow:
		return component.AreaTileKind_AREA_TILE_SNOW
	default:
		panic("invalid AreaTileKind")
	}
}
