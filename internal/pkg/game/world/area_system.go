package world

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AreaPositionUpdateFn func(areaPosition component.AreaPosition) (*component.AreaPosition, error)

type AreaTilesExtent struct {
	Left   uint32
	Top    uint32
	Right  uint32
	Bottom uint32
}

type AreaSystem struct {
	log   zerolog.Logger
	state *State

	areas          map[component.Entity]component.Area
	areasOccupancy map[component.Entity]map[component.AreaPositionLayer]*roaring64.Bitmap
	areasTiles     map[component.Entity]component.AreaTiles
	areasPositions map[component.Entity]component.AreaPosition
}

func NewAreaSystem(entities *State) *AreaSystem {
	return &AreaSystem{
		log:            log.With().Str("applicationComponent", "game").Str("gameComponent", "AreaSystem").Logger(),
		state:          entities,
		areas:          map[component.Entity]component.Area{},
		areasOccupancy: map[component.Entity]map[component.AreaPositionLayer]*roaring64.Bitmap{},
		areasTiles:     map[component.Entity]component.AreaTiles{},
		areasPositions: map[component.Entity]component.AreaPosition{},
	}
}

func (s *AreaSystem) Clone(newState *State) *AreaSystem {
	areasClone := map[component.Entity]component.Area{}
	areasPositionsClone := map[component.Entity]component.AreaPosition{}
	areasTilesClone := map[component.Entity]component.AreaTiles{}
	areasOccupancyClone := map[component.Entity]map[component.AreaPositionLayer]*roaring64.Bitmap{}

	for entity, area := range s.areas {
		areasClone[entity] = area
	}

	for entity, areaPosition := range s.areasPositions {
		areasPositionsClone[entity] = areaPosition
	}

	for entity, areaTiles := range s.areasTiles {
		areasTilesClone[entity] = append(areaTiles[:0:0], areaTiles...)
	}

	for entity, areaOccupancy := range s.areasOccupancy {
		areasOccupancyClone[entity] = map[component.AreaPositionLayer]*roaring64.Bitmap{
			component.AreaPositionLayerSurface: areaOccupancy[component.AreaPositionLayerSurface].Clone(),
			component.AreaPositionLayerPlayer:  areaOccupancy[component.AreaPositionLayerPlayer].Clone(),
		}
	}

	return &AreaSystem{
		log:            zerolog.Nop(),
		state:          newState,
		areas:          areasClone,
		areasPositions: areasPositionsClone,
		areasTiles:     areasTilesClone,
		areasOccupancy: areasOccupancyClone,
	}
}

func (s *AreaSystem) ValidatePosition(entity component.Entity, component component.AreaPosition) error {
	area, exists := s.areas[component.Entity]
	if !exists {
		return ErrAreaPositionEntityHasNoArea
	}

	if component.Width == 0 || component.Height == 0 {
		return ErrAreaPositionWithoutDimensions
	}

	if component.X+uint32(component.Width) > area.Width {
		return ErrAreaPositionOverflow
	}
	if component.Y+uint32(component.Height) > area.Height {
		return ErrAreaPositionOverflow
	}

	return nil
}

func (s *AreaSystem) ValidateArea(entity component.Entity, area component.Area, areaTiles component.AreaTiles) error {
	if area.Width == 0 || area.Height == 0 {
		return ErrAreaWithoutDimensions
	}
	if int(area.Width*area.Height) != len(areaTiles) {
		return ErrAreaTilesInvalidCount
	}

	return nil
}

func (s *AreaSystem) updatePosition(entity component.Entity, update AreaPositionUpdateFn) error {
	component, exists := s.areasPositions[entity]
	if !exists {
		return ErrAreaPositionComponentNotFound
	}

	componentAfterUpdate, err := update(component)
	if componentAfterUpdate == nil {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "update error")
	}

	if componentAfterUpdate.Layer != component.Layer {
		return ErrAreaPositionLayerImmutable
	}
	if componentAfterUpdate.Width != component.Width || componentAfterUpdate.Height != component.Height {
		return ErrAreaPositionDimensionsImmutable
	}

	if err := s.ValidatePosition(entity, *componentAfterUpdate); err != nil {
		return errors.Wrap(err, "unable to validate")
	}

	if err := s.movePosition(component, *componentAfterUpdate); err != nil {
		return errors.Wrap(err, "unable to move component")
	}

	s.areasPositions[entity] = *componentAfterUpdate

	return nil
}

func (s *AreaSystem) addPosition(entity component.Entity, areaPosition component.AreaPosition) error {
	if s.hasPosition(entity) {
		return ErrAreaPositionComponentAlreadyExists
	}

	if err := s.ValidatePosition(entity, areaPosition); err != nil {
		return errors.Wrap(err, "unable to validate area position")
	}

	if err := s.takePosition(areaPosition); err != nil {
		return errors.Wrap(err, "unable to take position")
	}

	s.areasPositions[entity] = areaPosition

	s.log.Info().EmbedObject(entity).EmbedObject(areaPosition).Msg("Added area position component.")

	return nil
}

func (s *AreaSystem) addArea(entity component.Entity, area component.Area, areaTiles component.AreaTiles) error {
	if s.hasArea(entity) {
		return ErrAreaComponentAlreadyExists
	}

	if err := s.ValidateArea(entity, area, areaTiles); err != nil {
		return errors.Wrap(err, "unable to validate area")
	}

	s.areas[entity] = area
	s.areasTiles[entity] = areaTiles
	s.areasOccupancy[entity] = map[component.AreaPositionLayer]*roaring64.Bitmap{
		component.AreaPositionLayerSurface: roaring64.New(),
		component.AreaPositionLayerPlayer:  roaring64.New(),
	}

	s.log.Info().EmbedObject(entity).EmbedObject(area).Msg("Added area component.")

	return nil
}

func (s *AreaSystem) GetPosition(entity component.Entity) (*component.AreaPosition, error) {
	if !s.hasPosition(entity) {
		return nil, ErrAreaPositionComponentNotFound
	}

	areaPositionCopy := s.areasPositions[entity]

	return &areaPositionCopy, nil
}

func (s *AreaSystem) GetArea(entity component.Entity) (*component.Area, error) {
	area, exists := s.areas[entity]
	if !exists {
		return nil, ErrAreaComponentNotFound
	}

	areaCopy := area

	return &areaCopy, nil
}

func (s *AreaSystem) GetAreaTiles(entity component.Entity, extent AreaTilesExtent) ([]component.AreaTile, error) {
	area, exists := s.areas[entity]
	if !exists {
		return []component.AreaTile{}, ErrAreaComponentNotFound
	}

	_, exists = s.areasTiles[entity]
	if !exists {
		return []component.AreaTile{}, ErrAreaComponentTilesNotFound
	}

	var tiles []component.AreaTile

	if extent.Left >= area.Width || extent.Right >= area.Width {
		return []component.AreaTile{}, ErrAreaTileOutOfBounds
	}
	if extent.Top >= area.Height || extent.Bottom >= area.Height {
		return []component.AreaTile{}, ErrAreaTileOutOfBounds
	}

	for y := extent.Top; y <= extent.Bottom; y++ {
		tileLineIndex := extent.Left + (y * area.Width)
		tileLineWidth := extent.Right - extent.Left + 1
		tiles = append(
			tiles,
			s.areasTiles[entity][tileLineIndex:tileLineIndex+tileLineWidth]...,
		)
	}

	return tiles, nil
}

func (s *AreaSystem) GetTile(entity component.Entity, x uint32, y uint32) (*component.AreaTile, error) {
	area, exists := s.areas[entity]
	if !exists {
		return nil, ErrAreaComponentNotFound
	}

	_, exists = s.areasTiles[entity]
	if !exists {
		return nil, ErrAreaComponentTilesNotFound
	}

	if x >= area.Width || y >= area.Height {
		return nil, ErrAreaTileOutOfBounds
	}

	tileIndex := x + (y * area.Width)

	tileCopy := s.areasTiles[entity][tileIndex]

	return &tileCopy, nil
}

func (s *AreaSystem) hasPosition(entity component.Entity) bool {
	if _, exists := s.areasPositions[entity]; exists {
		return true
	}

	return false
}

func (s *AreaSystem) removePosition(entity component.Entity) error {
	if !s.hasPosition(entity) {
		return ErrAreaPositionComponentNotFound
	}

	delete(s.areasPositions, entity)

	return nil
}

func (s *AreaSystem) hasArea(entity component.Entity) bool {
	if _, exists := s.areas[entity]; exists {
		return true
	}

	return false
}

func (s *AreaSystem) removeArea(entity component.Entity) error {
	if !s.hasArea(entity) {
		return ErrAreaComponentNotFound
	}

	delete(s.areas, entity)
	delete(s.areasTiles, entity)
	delete(s.areasOccupancy, entity)

	return nil
}

func (s *AreaSystem) applyDeltaTime(delta uint64) error {
	return nil
}

func (s *AreaSystem) movePosition(previousAreaPosition, areaPosition component.AreaPosition) error {
	previousAreaPositionBitmap, err := s.areaPositionToBitmap(previousAreaPosition)
	if err != nil {
		return err
	}

	areaPositionBitmap, err := s.areaPositionToBitmap(areaPosition)
	if err != nil {
		return err
	}

	bitmap := s.areasOccupancy[areaPosition.Entity][areaPosition.Layer]
	bitmapClone := bitmap.Clone()

	// Running on real
	for _, previousAreaPositionBitmapIndex := range previousAreaPositionBitmap.ToArray() {
		if !bitmap.Contains(previousAreaPositionBitmapIndex) {
			return ErrAreaPositionNotTaken
		}
	}

	// Running on clone
	for _, previousAreaPositionBitmapIndex := range previousAreaPositionBitmap.ToArray() {
		bitmapClone.Remove(previousAreaPositionBitmapIndex)
	}

	// Running on clone
	for _, areaPositionIndex := range areaPositionBitmap.ToArray() {
		if bitmapClone.Contains(areaPositionIndex) {
			return ErrAreaPositionAlreadyTaken
		}
	}

	// Running on real
	bitmap.AddMany(areaPositionBitmap.ToArray())

	return nil
}

func (s *AreaSystem) takePosition(areaPosition component.AreaPosition) error {
	bitmap := s.areasOccupancy[areaPosition.Entity][areaPosition.Layer]

	areaPositionBitmap, err := s.areaPositionToBitmap(areaPosition)
	if err != nil {
		return err
	}

	for _, areaPositionIndex := range areaPositionBitmap.ToArray() {
		if bitmap.Contains(areaPositionIndex) {
			return ErrAreaPositionAlreadyTaken
		}
	}

	bitmap.AddMany(areaPositionBitmap.ToArray())

	return nil
}

func (s *AreaSystem) releasePosition(areaPosition component.AreaPosition) error {
	bitmap := s.areasOccupancy[areaPosition.Entity][areaPosition.Layer]

	areaPositionBitmap, err := s.areaPositionToBitmap(areaPosition)
	if err != nil {
		return err
	}

	for _, areaPositionIndex := range areaPositionBitmap.ToArray() {
		if !bitmap.Contains(areaPositionIndex) {
			return ErrAreaPositionNotTaken
		}
	}

	for _, areaPositionIndex := range areaPositionBitmap.ToArray() {
		bitmap.Remove(areaPositionIndex)
	}

	return nil
}

func (s *AreaSystem) areaPositionToBitmap(areaPosition component.AreaPosition) (*roaring64.Bitmap, error) {
	area, exists := s.areas[areaPosition.Entity]
	if !exists {
		return nil, ErrAreaComponentNotFound
	}

	bitmap := roaring64.New()

	for y := areaPosition.Y; y < areaPosition.Y+uint32(areaPosition.Height); y++ {
		for x := areaPosition.X; x < areaPosition.X+uint32(areaPosition.Width); x++ {
			index := uint64((y * area.Width) + x)

			bitmap.Add(index)
		}
	}

	return bitmap, nil
}

var (
	ErrAreaComponentNotFound              = errors.New("area component not found")
	ErrAreaComponentTilesNotFound         = errors.New("area component tiles not found")
	ErrAreaWithoutDimensions              = errors.New("area without dimensions")
	ErrAreaPositionComponentNotFound      = errors.New("area position components not found")
	ErrAreaComponentAlreadyExists         = errors.New("area component already exists")
	ErrAreaTileOutOfBounds                = errors.New("area tile out of bounds")
	ErrAreaTilesInvalidCount              = errors.New("area tiles invalid count")
	ErrAreaPositionEntityHasNoArea        = errors.New("area position entity has no area")
	ErrAreaPositionAlreadyTaken           = errors.New("area position already taken")
	ErrAreaPositionNotTaken               = errors.New("area position not taken")
	ErrAreaPositionWithoutDimensions      = errors.New("area position without dimensions")
	ErrAreaPositionOverflow               = errors.New("area position overflow")
	ErrAreaPositionOverlapping            = errors.New("area position overlapping")
	ErrAreaPositionComponentAlreadyExists = errors.New("area position component already exists")
	ErrAreaPositionLayerImmutable         = errors.New("area position layer immutable")
	ErrAreaPositionDimensionsImmutable    = errors.New("area position dimensions immutable")
)
