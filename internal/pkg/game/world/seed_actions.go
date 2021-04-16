package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
	"math/rand"
)

type SeedActions struct {
	state *State
}

func newSeedActions(state *State) *SeedActions {
	return &SeedActions{state: state}
}

func (b *SeedActions) CreateRandomSeedsOnPlanet(planetEntity component.Entity) ([]component.Entity, error) {
	planetSystem := b.state.Planet()
	areaSystem := b.state.Area()

	planet, err := planetSystem.Get(planetEntity)
	if err != nil {
		return []component.Entity{}, errors.Wrap(err, "unable to get planet")
	}

	area, err := areaSystem.GetArea(planetEntity)
	if err != nil {
		return []component.Entity{}, errors.Wrap(err, "unable to get area")
	}

	areaTiles, err := areaSystem.GetAreaTiles(planetEntity, AreaTilesExtent{
		Left:   0,
		Top:    0,
		Right:  area.Width - 1,
		Bottom: area.Height - 1,
	})
	if err != nil {
		return []component.Entity{}, errors.Wrap(err, "unable to get planet tiles")
	}

	generator := rand.New(rand.NewSource(planet.Seed))

	var seedEntities []component.Entity
	var x, y uint32

	for y = 0; y < area.Width; y++ {
		for x = 0; x < area.Height; x++ {
			index := x + (y * area.Width)

			if areaTiles[index].OwnerEntity != planetEntity {
				continue
			}

			if !((areaTiles[index].Kind == component.AreaTileKindGround) ||
				(areaTiles[index].Kind == component.AreaTileKindFertileGround)) {
				continue
			}

			generatorValue := generator.Float64()

			if generatorValue > 0.99999 {
				seedEntity, err := b.CreateOakSeed(planetEntity, planetEntity, x, y)
				if err != nil {
					return []component.Entity{}, errors.Wrapf(err, "unable to create oak seed at %d,%d", x, y)
				}
				seedEntities = append(seedEntities, *seedEntity)
				continue
			}

			if generatorValue > 0.99995 {
				seedEntity, err := b.CreatePineSeed(planetEntity, planetEntity, x, y)
				if err != nil {
					return []component.Entity{}, errors.Wrapf(err, "unable to create pine seed at %d,%d", x, y)
				}
				seedEntities = append(seedEntities, *seedEntity)
				continue
			}

			if generatorValue > 0.99990 {
				seedEntity, err := b.CreateWheatSeed(planetEntity, planetEntity, x, y)
				if err != nil {
					return []component.Entity{}, errors.Wrapf(err, "unable to create wheat seed at %d,%d", x, y)
				}
				seedEntities = append(seedEntities, *seedEntity)
				continue
			}

		}
	}

	return seedEntities, nil
}

func (b *SeedActions) CreateWheatSeed(owner, planet component.Entity, x, y uint32) (*component.Entity, error) {
	entity := b.state.Create(component.EntityKindSeedWheat)

	seed := component.Seed{
		Kind:     component.SeedKindWheat,
		Maturity: 0,
	}

	return b.create(entity, owner, planet, x, y, seed)
}

func (b *SeedActions) CreatePineSeed(owner, planet component.Entity, x, y uint32) (*component.Entity, error) {
	entity := b.state.Create(component.EntityKindSeedPineTree)

	seed := component.Seed{
		Kind:     component.SeedKindPineTree,
		Maturity: 0,
	}

	return b.create(entity, owner, planet, x, y, seed)
}

func (b *SeedActions) CreateOakSeed(owner, planet component.Entity, x, y uint32) (*component.Entity, error) {
	entity := b.state.Create(component.EntityKindSeedOakTree)

	seed := component.Seed{
		Kind:     component.SeedKindOakTree,
		Maturity: 0,
	}

	return b.create(entity, owner, planet, x, y, seed)
}

func (b *SeedActions) create(entity, owner, planet component.Entity, x, y uint32, seed component.Seed) (*component.Entity, error) {
	position := component.AreaPosition{
		Entity: planet,
		X:      x,
		Y:      y,
		Layer:  component.AreaPositionLayerSurface,
		Width:  1,
		Height: 1,
	}

	possession := component.Possession{
		OwnerEntity: owner,
	}

	if err := b.state.area.addPosition(entity, position); err != nil {
		return nil, errors.Wrap(err, "unable to add area position component to seed entity")
	}

	if err := b.state.seed.add(entity, seed); err != nil {
		return nil, errors.Wrap(err, "unable to add seed component to seed entity")
	}

	if err := b.state.possession.add(entity, possession); err != nil {
		return nil, errors.Wrap(err, "unable to add possessions component to seed entity")
	}

	return &entity, nil
}
