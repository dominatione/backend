package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PlanetSystemFilterFn func(planet component.Planet) bool

type PlanetSystem struct {
	log   zerolog.Logger
	state *State

	planets map[component.Entity]component.Planet
}

func newPlanetSystem(state *State) *PlanetSystem {
	return &PlanetSystem{
		log:     log.With().Str("applicationComponent", "game").Str("gameComponent", "PlanetSystem").Logger(),
		state:   state,
		planets: map[component.Entity]component.Planet{},
	}
}

func (s *PlanetSystem) clone(newState *State) *PlanetSystem {
	planetsClone := map[component.Entity]component.Planet{}

	for entity, planet := range s.planets {
		planetsClone[entity] = component.Planet{
			Seed: planet.Seed,
			Name: planet.Name,
		}
	}

	return &PlanetSystem{
		log:     zerolog.Nop(),
		state:   newState,
		planets: planetsClone,
	}
}

func (s *PlanetSystem) remove(entity component.Entity) error {
	if !s.exists(entity) {
		return ErrPlanetComponentNotFound
	}

	delete(s.planets, entity)

	return nil
}

func (s *PlanetSystem) exists(entity component.Entity) bool {
	_, exists := s.planets[entity]

	return exists
}

func (s *PlanetSystem) validate(entity component.Entity, planet component.Planet) error {
	return nil
}

func (s *PlanetSystem) add(entity component.Entity, planet component.Planet) error {
	if s.exists(entity) {
		return ErrPlanetAlreadyExists
	}

	if err := s.validate(entity, planet); err != nil {
		return errors.Wrap(err, "unable to validate")
	}

	s.planets[entity] = planet

	s.log.Info().EmbedObject(entity).EmbedObject(planet).Msg("Added planet component.")

	return nil
}

func (s *PlanetSystem) Get(entity component.Entity) (*component.Planet, error) {
	if !s.exists(entity) {
		return nil, ErrPlanetComponentNotFound
	}

	planetCopy := s.planets[entity]

	return &planetCopy, nil
}

func (s *PlanetSystem) Count() int {
	return len(s.planets)
}

func (s *PlanetSystem) Entities() []component.Entity {
	entities := []component.Entity{}

	for entity := range s.planets {
		entities = append(entities, entity)
	}

	return entities
}

func (s *PlanetSystem) applyDeltaTime(delta uint64) error {
	return nil
}

var (
	ErrPlanetAlreadyExists     = errors.New("planet already hasPosition")
	ErrPlanetComponentNotFound = errors.New("planet component not found")
)
