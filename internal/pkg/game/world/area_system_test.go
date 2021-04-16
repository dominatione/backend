package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAreaSystem_Clone(t *testing.T) {
	state := NewState()

	stateClone := state.Clone()

	assert.NotNil(t, stateClone.area.state)
	assert.NotNil(t, stateClone.area.areas)
	assert.NotNil(t, stateClone.area.areasOccupancy)
	assert.NotNil(t, stateClone.area.areasPositions)
	assert.NotNil(t, stateClone.area.areasTiles)
}

func TestAreaSystem_GetArea(t *testing.T) {
	var area *component.Area
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entity := state.Create(component.EntityKindUnknown)

	area, err = state.area.GetArea(entity)
	assert.Nil(t, area)
	assert.ErrorIs(t, err, ErrAreaComponentNotFound)

	err = state.area.addArea(entity, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	area, err = state.area.GetArea(entity)
	assert.NoError(t, err)
	assert.NotNil(t, area)
	assert.EqualValues(t, width, area.Width)
	assert.EqualValues(t, height, area.Height)
}

func TestAreaSystem_GetAreaTiles(t *testing.T) {
	var areaTiles component.AreaTiles
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entity := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entity, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	areaTiles, err = state.area.GetAreaTiles(entity, AreaTilesExtent{0, 0, 0, 0})
	assert.NoError(t, err)
	assert.Len(t, areaTiles, 1)

	areaTiles, err = state.area.GetAreaTiles(entity, AreaTilesExtent{width, 0, 0, 0})
	assert.ErrorIs(t, err, ErrAreaTileOutOfBounds)
	assert.Len(t, areaTiles, 0)

	areaTiles, err = state.area.GetAreaTiles(entity, AreaTilesExtent{0, height, 0, 0})
	assert.ErrorIs(t, err, ErrAreaTileOutOfBounds)
	assert.Len(t, areaTiles, 0)

	areaTiles, err = state.area.GetAreaTiles(entity, AreaTilesExtent{0, 0, width, 0})
	assert.ErrorIs(t, err, ErrAreaTileOutOfBounds)
	assert.Len(t, areaTiles, 0)

	areaTiles, err = state.area.GetAreaTiles(entity, AreaTilesExtent{0, 0, 0, height})
	assert.ErrorIs(t, err, ErrAreaTileOutOfBounds)
	assert.Len(t, areaTiles, 0)

	areaTiles, err = state.area.GetAreaTiles(entity, AreaTilesExtent{0, 0, width - 1, height - 1})
	assert.NoError(t, err)
	assert.Len(t, areaTiles, int(width*height))
}

func TestAreaSystem_GetTile(t *testing.T) {
	var areaTile *component.AreaTile
	var err error

	width := uint32(2)
	height := uint32(2)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithoutArea := state.Create(component.EntityKindUnknown)

	areaTiles := createAreaTiles(width, height, component.AreaTileKindGround)
	areaTiles[0].Kind = component.AreaTileKindGravel

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, areaTiles)
	assert.NoError(t, err)

	areaTile, err = state.area.GetTile(entityWithoutArea, 0, 0)
	assert.ErrorIs(t, err, ErrAreaComponentNotFound)
	assert.Nil(t, areaTile)

	areaTile, err = state.area.GetTile(entityWithArea, width, 0)
	assert.ErrorIs(t, err, ErrAreaTileOutOfBounds)
	assert.Nil(t, areaTile)

	areaTile, err = state.area.GetTile(entityWithArea, 0, height)
	assert.ErrorIs(t, err, ErrAreaTileOutOfBounds)
	assert.Nil(t, areaTile)

	areaTile, err = state.area.GetTile(entityWithArea, 0, 0)
	assert.NoError(t, err)
	assert.EqualValues(t, component.AreaTileKindGravel, areaTile.Kind)
	areaTile.Kind = component.AreaTileKindGround

	areaTile, err = state.area.GetTile(entityWithArea, 0, 0)
	assert.NoError(t, err)
	assert.EqualValues(t, component.AreaTileKindGravel, areaTile.Kind)
}

func TestAreaSystem_GetPosition(t *testing.T) {
	var areaPosition *component.AreaPosition
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithPosition := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	areaPosition, err = state.area.GetPosition(entityWithPosition)
	assert.ErrorIs(t, err, ErrAreaPositionComponentNotFound)
	assert.Nil(t, areaPosition)

	err = state.area.addPosition(entityWithPosition, component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      1,
		Y:      2,
		Width:  3,
		Height: 4,
	})
	assert.NoError(t, err)

	areaPosition, err = state.area.GetPosition(entityWithPosition)
	assert.NoError(t, err)
	assert.NotNil(t, areaPosition)
	assert.EqualValues(t, 1, areaPosition.X)
	assert.EqualValues(t, 2, areaPosition.Y)
	assert.EqualValues(t, 3, areaPosition.Width)
	assert.EqualValues(t, 4, areaPosition.Height)
}

func TestAreaSystem_addArea(t *testing.T) {
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entity := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entity, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	err = state.area.addArea(entity, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.ErrorIs(t, err, ErrAreaComponentAlreadyExists)
}

func TestAreaSystem_takePosition(t *testing.T) {
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithoutArea := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	err = state.area.takePosition(component.AreaPosition{
		Entity: entityWithoutArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.ErrorIs(t, err, ErrAreaComponentNotFound)

	err = state.area.takePosition(component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.NoError(t, err)

	err = state.area.takePosition(component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.ErrorIs(t, err, ErrAreaPositionAlreadyTaken)
}

func TestAreaSystem_releasePosition(t *testing.T) {
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithoutArea := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	err = state.area.releasePosition(component.AreaPosition{
		Entity: entityWithoutArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.ErrorIs(t, err, ErrAreaComponentNotFound)

	err = state.area.releasePosition(component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.ErrorIs(t, err, ErrAreaPositionNotTaken)

	err = state.area.takePosition(component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.NoError(t, err)

	err = state.area.releasePosition(component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.NoError(t, err)

	err = state.area.takePosition(component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	})
	assert.NoError(t, err)
}

func TestAreaSystem_addPosition(t *testing.T) {
	var err error
	var areaPosition component.AreaPosition

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithPosition := state.Create(component.EntityKindUnknown)
	entityWithSamePosition := state.Create(component.EntityKindUnknown)
	entityWithoutArea := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	areaPosition = component.AreaPosition{
		Entity: entityWithoutArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	}
	err = state.area.addPosition(entityWithPosition, areaPosition)
	assert.Error(t, err)

	areaPosition = component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	}
	err = state.area.addPosition(entityWithPosition, areaPosition)
	assert.NoError(t, err)

	err = state.area.addPosition(entityWithSamePosition, areaPosition)
	assert.Error(t, err)

	areaPosition = component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	}
	err = state.area.addPosition(entityWithPosition, areaPosition)
	assert.ErrorIs(t, err, ErrAreaPositionComponentAlreadyExists)
}

func TestAreaSystem_ValidatePosition(t *testing.T) {
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithoutArea := state.Create(component.EntityKindUnknown)
	entityWithPosition := state.Create(component.EntityKindUnknown)

	areaTiles := createAreaTiles(width, height, component.AreaTileKindGround)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, areaTiles)
	assert.NoError(t, err)

	err = state.area.ValidatePosition(entityWithPosition, component.AreaPosition{
		Entity: entityWithoutArea,
		Layer:  0,
		X:      0,
		Y:      0,
		Width:  0,
		Height: 0,
	})
	assert.ErrorIs(t, err, ErrAreaPositionEntityHasNoArea)

	err = state.area.ValidatePosition(entityWithPosition, component.AreaPosition{
		Entity: entityWithArea,
		Layer:  0,
		X:      0,
		Y:      0,
		Width:  0,
		Height: 0,
	})
	assert.ErrorIs(t, err, ErrAreaPositionWithoutDimensions)

	err = state.area.ValidatePosition(entityWithPosition, component.AreaPosition{
		Entity: entityWithArea,
		Layer:  0,
		X:      0,
		Y:      0,
		Width:  20,
		Height: 1,
	})
	assert.ErrorIs(t, err, ErrAreaPositionOverflow)

	err = state.area.ValidatePosition(entityWithPosition, component.AreaPosition{
		Entity: entityWithArea,
		Layer:  0,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 20,
	})
	assert.ErrorIs(t, err, ErrAreaPositionOverflow)
}

func TestAreaSystem_updatePosition(t *testing.T) {
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithPosition := state.Create(component.EntityKindUnknown)
	entityCollider := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	err = state.area.updatePosition(entityWithPosition, func(areaPosition component.AreaPosition) (*component.AreaPosition, error) {
		return &areaPosition, nil
	})
	assert.ErrorIs(t, err, ErrAreaPositionComponentNotFound)

	err = state.area.addPosition(entityCollider, component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      2,
		Y:      2,
		Width:  1,
		Height: 1,
	})
	assert.NoError(t, err)

	err = state.area.addPosition(entityWithPosition, component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  2,
		Height: 2,
	})
	assert.NoError(t, err)

	err = state.area.updatePosition(entityWithPosition, func(areaPosition component.AreaPosition) (*component.AreaPosition, error) {
		areaPosition.Layer = component.AreaPositionLayerEmpty
		return &areaPosition, nil
	})
	assert.ErrorIs(t, err, ErrAreaPositionLayerImmutable)

	err = state.area.updatePosition(entityWithPosition, func(areaPosition component.AreaPosition) (*component.AreaPosition, error) {
		areaPosition.Width = 3
		return &areaPosition, nil
	})
	assert.ErrorIs(t, err, ErrAreaPositionDimensionsImmutable)

	err = state.area.updatePosition(entityWithPosition, func(areaPosition component.AreaPosition) (*component.AreaPosition, error) {
		areaPosition.Height = 3
		return &areaPosition, nil
	})
	assert.ErrorIs(t, err, ErrAreaPositionDimensionsImmutable)

	err = state.area.updatePosition(entityWithPosition, func(areaPosition component.AreaPosition) (*component.AreaPosition, error) {
		areaPosition.X = width + 1
		areaPosition.Y = height + 1
		return &areaPosition, nil
	})
	assert.Error(t, err)

	err = state.area.updatePosition(entityWithPosition, func(areaPosition component.AreaPosition) (*component.AreaPosition, error) {
		areaPosition.X = 1
		areaPosition.Y = 1
		return &areaPosition, nil
	})
	assert.Error(t, err, ErrAreaPositionAlreadyTaken)

	err = state.area.updatePosition(entityWithPosition, func(areaPosition component.AreaPosition) (*component.AreaPosition, error) {
		areaPosition.X = 4
		areaPosition.Y = 4
		return &areaPosition, nil
	})
	assert.NoError(t, err)

	areaPosition, err := state.area.GetPosition(entityWithPosition)
	assert.NoError(t, err)
	assert.NotNil(t, areaPosition)
	assert.EqualValues(t, 4, areaPosition.X)
	assert.EqualValues(t, 4, areaPosition.Y)
}

func TestAreaSystem_ValidateArea(t *testing.T) {
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)

	areaTiles := createAreaTiles(width, height, component.AreaTileKindGround)

	err = state.area.ValidateArea(entityWithArea, component.Area{
		Width:  width + 1,
		Height: height,
	}, areaTiles)
	assert.ErrorIs(t, err, ErrAreaTilesInvalidCount)

	err = state.area.ValidateArea(entityWithArea, component.Area{
		Width:  0,
		Height: 0,
	}, areaTiles)
	assert.ErrorIs(t, err, ErrAreaWithoutDimensions)

	err = state.area.ValidateArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, areaTiles)
	assert.NoError(t, err)
}

func TestAreaSystem_removePosition(t *testing.T) {
	var err error

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)
	entityWithPosition := state.Create(component.EntityKindUnknown)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	err = state.area.removePosition(entityWithPosition)
	assert.ErrorIs(t, err, ErrAreaPositionComponentNotFound)

	areaPosition := component.AreaPosition{
		Entity: entityWithArea,
		Layer:  component.AreaPositionLayerSurface,
		X:      0,
		Y:      0,
		Width:  1,
		Height: 1,
	}
	err = state.area.addPosition(entityWithPosition, areaPosition)
	assert.NoError(t, err)

	err = state.area.removePosition(entityWithPosition)
	assert.NoError(t, err)

}

func TestAreaSystem_removeArea(t *testing.T) {
	var err error
	var area *component.Area

	width := uint32(10)
	height := uint32(10)

	state := NewState()
	entityWithArea := state.Create(component.EntityKindUnknown)

	err = state.area.removeArea(entityWithArea)
	assert.ErrorIs(t, err, ErrAreaComponentNotFound)

	err = state.area.addArea(entityWithArea, component.Area{
		Width:  width,
		Height: height,
	}, createAreaTiles(width, height, component.AreaTileKindGround))
	assert.NoError(t, err)

	area, err = state.area.GetArea(entityWithArea)
	assert.NoError(t, err)
	assert.NotNil(t, area)

	err = state.area.removeArea(entityWithArea)
	assert.NoError(t, err)

	area, err = state.area.GetArea(entityWithArea)
	assert.ErrorIs(t, err, ErrAreaComponentNotFound)
	assert.Nil(t, area)
}

func TestAreaSystem_applyDeltaTime(t *testing.T) {
	var err error

	state := NewState()

	err = state.area.applyDeltaTime(1000)
	assert.NoError(t, err)
}

func createAreaTiles(width, height uint32, kind component.AreaTileKind) component.AreaTiles {
	size := uint32(width * height)
	tiles := make(component.AreaTiles, size)

	for index := uint32(0); index < size; index++ {
		tiles[index].Kind = kind
	}

	return tiles
}
