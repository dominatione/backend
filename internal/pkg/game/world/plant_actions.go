package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
)

type PlantActions struct {
	state *State
}

func newPlantActions(state *State) *PlantActions {
	return &PlantActions{
		state: state,
	}
}

func (f *PlantActions) CreateFromSeedAndRemoveSeed(seedEntity component.Entity) (*component.Entity, error) {
	areaPosition, err := f.state.area.GetPosition(seedEntity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get area position component from seed entity")
	}

	seed, err := f.state.seed.Get(seedEntity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get seed component from seed entity")
	}

	possession, err := f.state.possession.Get(seedEntity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get possession component from seed entity")
	}

	var plantKind component.PlantKind
	var entityKind component.EntityKind

	switch seed.Kind {
	case component.SeedKindOakTree:
		plantKind = component.PlantKindOakTree
		entityKind = component.EntityKindPlantOakTree
	case component.SeedKindPineTree:
		plantKind = component.PlantKindPineTree
		entityKind = component.EntityKindPlantPineTree
	case component.SeedKindWheat:
		plantKind = component.PlantKindWheat
		entityKind = component.EntityKindPlantWheat
	case component.SeedKindCorn:
		plantKind = component.PlantKindCorn
		entityKind = component.EntityKindPlantCorn
	case component.SeedKindCannabis:
		plantKind = component.PlantKindCannabis
		entityKind = component.EntityKindPlantCannabis
	default:
		return nil, ErrUnsupportedSeed
	}

	plant := component.Plant{
		Kind:               plantKind,
		Maturity:           0,
		AnemochoryMaturity: 0,
	}

	err = f.state.Remove(seedEntity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to remove seed entity")
	}

	plantEntity := f.state.Create(entityKind)

	err = f.state.area.addPosition(plantEntity, *areaPosition)
	if err != nil {
		return nil, errors.Wrap(err, "unable to add area position component to plant entity")
	}

	err = f.state.plant.add(plantEntity, plant)
	if err != nil {
		return nil, errors.Wrap(err, "unable to add plant component to plant entit")
	}

	err = f.state.possession.add(plantEntity, *possession)
	if err != nil {
		return nil, errors.Wrap(err, "unable to add possession component to plant entity")
	}

	return &plantEntity, nil
}

var (
	ErrUnsupportedSeed = errors.New("unsupported seed")
)
