package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PlantSystem struct {
	log   zerolog.Logger
	state *State

	plants map[component.Entity]component.Plant
}

func newPlantSystem(state *State) *PlantSystem {
	return &PlantSystem{
		log:    log.With().Str("applicationComponent", "game").Str("gameComponent", "PlantSystem").Logger(),
		state:  state,
		plants: map[component.Entity]component.Plant{},
	}
}

func (s *PlantSystem) clone(newState *State) *PlantSystem {
	plantsClone := map[component.Entity]component.Plant{}

	for entity, plant := range s.plants {
		plantsClone[entity] = plant
	}

	return &PlantSystem{
		log:    zerolog.Nop(),
		state:  newState,
		plants: plantsClone,
	}
}

func (s *PlantSystem) validate(entity component.Entity, plant component.Plant) error {
	return nil
}

func (s *PlantSystem) add(entity component.Entity, plant component.Plant) error {
	if s.exists(entity) {
		return ErrPlantComponentAlreadyExists
	}

	if err := s.validate(entity, plant); err != nil {
		return errors.Wrap(err, "unable to validate")
	}

	s.plants[entity] = plant

	s.log.Info().EmbedObject(entity).EmbedObject(plant).Msg("Added plant component.")

	return nil
}

func (s *PlantSystem) Entities() []component.Entity {
	entities := []component.Entity{}

	for entity := range s.plants {
		entities = append(entities, entity)
	}

	return entities
}

func (s *PlantSystem) exists(entity component.Entity) bool {
	_, exists := s.plants[entity]

	return exists
}

func (s *PlantSystem) remove(entity component.Entity) error {
	if !s.exists(entity) {
		return ErrPlantComponentNotFound
	}

	delete(s.plants, entity)

	return nil
}

func (s *PlantSystem) applyDeltaTime(delta uint64) error {
	return nil
}

var (
	ErrPlantComponentNotFound      = errors.New("plant component not found")
	ErrPlantComponentAlreadyExists = errors.New("plant component already hasPosition")
)
