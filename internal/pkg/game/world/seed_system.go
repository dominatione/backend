package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type SeedUpdateFn func(seed component.Seed) (*component.Seed, error)

type SeedSystem struct {
	log   zerolog.Logger
	state *State

	seedsMutex sync.Mutex
	seeds      map[component.Entity]component.Seed
}

func newSeedSystem(state *State) *SeedSystem {
	return &SeedSystem{
		log:   log.With().Str("applicationComponent", "game").Str("gameComponent", "SeedSystem").Logger(),
		state: state,
		seeds: map[component.Entity]component.Seed{},
	}
}

func (s *SeedSystem) clone(newState *State) *SeedSystem {
	seedsClone := map[component.Entity]component.Seed{}

	for entity, seed := range s.seeds {
		seedsClone[entity] = seed
	}

	return &SeedSystem{
		log:   zerolog.Nop(),
		state: newState,
		seeds: seedsClone,
	}
}

func (s *SeedSystem) validate(entity component.Entity, seed component.Seed) error {
	if seed.Maturity > 1.0 {
		return ErrSeedComponentMaturityOverflow
	}

	return nil
}

func (s *SeedSystem) add(entity component.Entity, seed component.Seed) error {
	if s.exists(entity) {
		return ErrSeedComponentAlreadyExists
	}

	if err := s.validate(entity, seed); err != nil {
		return errors.Wrap(err, "unable to validate")
	}

	s.seedsMutex.Lock()
	s.seeds[entity] = seed
	s.seedsMutex.Unlock()

	s.log.Info().EmbedObject(entity).EmbedObject(seed).Msg("Added seed component.")

	return nil
}

func (s *SeedSystem) update(entity component.Entity, update SeedUpdateFn) error {
	seed, exists := s.seeds[entity]
	if !exists {
		return ErrSeedComponentNotFound
	}

	updatedSeed, err := update(seed)
	if err != nil {
		return errors.Wrap(err, "updatePosition function failed")
	}
	if updatedSeed == nil {
		return nil
	}

	if err := s.validate(entity, *updatedSeed); err != nil {
		return errors.Wrap(err, "unable to validate after updatePosition")
	}

	s.seedsMutex.Lock()
	s.seeds[entity] = *updatedSeed
	s.seedsMutex.Unlock()

	return nil
}

func (s *SeedSystem) Entities() []component.Entity {
	entities := []component.Entity{}

	s.seedsMutex.Lock()
	for entity := range s.seeds {
		entities = append(entities, entity)
	}
	s.seedsMutex.Unlock()

	return entities
}

func (s *SeedSystem) Get(entity component.Entity) (*component.Seed, error) {
	defer s.seedsMutex.Unlock()
	s.seedsMutex.Lock()

	seed, exists := s.seeds[entity]
	if !exists {
		return nil, ErrSeedComponentNotFound
	}

	return &seed, nil
}

func (s *SeedSystem) applyDeltaTime(delta uint64) error {
	deltaSeconds := float32(delta) / 1000

	for entity := range s.seeds {
		err := s.update(entity, func(seed component.Seed) (*component.Seed, error) {
			if !s.state.area.hasPosition(entity) {
				return nil, nil
			}

			switch seed.Kind {
			case component.SeedKindOakTree:
				seed.Maturity += ThreeDaysDeltaFactor * deltaSeconds
			case component.SeedKindPineTree:
				seed.Maturity += TwoDaysDeltaFactor * deltaSeconds
			case component.SeedKindWheat:
				seed.Maturity += TwoDaysDeltaFactor * deltaSeconds
			case component.SeedKindCorn:
				seed.Maturity += TwoDaysDeltaFactor * deltaSeconds
			case component.SeedKindCannabis:
				seed.Maturity += DayDeltaFactor * deltaSeconds
			default:
				s.log.Panic().Msg("Unsupported seed.")
			}

			if seed.Maturity > 1 {
				_, err := s.state.actions.plant.CreateFromSeedAndRemoveSeed(entity)
				if err != nil {
					return nil, errors.Wrap(err, "unable to create plant component from seed and removePosition seed component")
				}
				return nil, nil
			} else {
				return &seed, nil
			}
		})
		if err != nil {
			return errors.Wrapf(err, "unable to apply delta time on seed %s", entity)
		}
	}

	return nil
}

func (s *SeedSystem) remove(entity component.Entity) error {
	if !s.exists(entity) {
		return ErrSeedComponentNotFound
	}

	seed := s.seeds[entity]

	delete(s.seeds, entity)

	s.log.Info().EmbedObject(entity).EmbedObject(seed).Msg("Remove seed component.")

	return nil
}

func (s *SeedSystem) exists(entity component.Entity) bool {
	_, exists := s.seeds[entity]

	return exists
}

var (
	ErrSeedComponentAlreadyExists    = errors.New("seed component already hasPosition")
	ErrSeedComponentNotFound         = errors.New("seed component not found")
	ErrSeedComponentMaturityOverflow = errors.New("seed component maturity overflow")
)
