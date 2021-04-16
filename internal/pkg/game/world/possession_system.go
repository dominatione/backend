package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type PossessionSystemFilterFn func(possession component.Possession) bool

type PossessionSystem struct {
	log   zerolog.Logger
	state *State

	possessionsMutex sync.Mutex
	possessions      map[component.Entity]component.Possession
}

func newPossessionSystem(state *State) *PossessionSystem {
	return &PossessionSystem{
		log:         log.With().Str("applicationComponent", "game").Str("gameComponent", "PossessionSystem").Logger(),
		state:       state,
		possessions: map[component.Entity]component.Possession{},
	}
}

func (s *PossessionSystem) clone(newState *State) *PossessionSystem {
	possessionClone := map[component.Entity]component.Possession{}

	for entity, possession := range s.possessions {
		possessionClone[entity] = possession
	}

	return &PossessionSystem{
		log:         zerolog.Nop(),
		state:       newState,
		possessions: possessionClone,
	}
}

func (s *PossessionSystem) validate(entity component.Entity, possession component.Possession) error {
	ownerKind, err := s.state.GetKind(possession.OwnerEntity)
	if err != nil {
		return errors.Wrap(err, "unable to get owner component kind")
	}

	allowedOwnerKind := (*ownerKind == component.EntityKindPlanet) ||
		(*ownerKind == component.EntityKindPlayer)

	if !allowedOwnerKind {
		return ErrPossessionOwnerInvalidEntity
	}

	return nil
}

func (s *PossessionSystem) add(entity component.Entity, possession component.Possession) error {
	if s.exists(entity) {
		return ErrPossessionAlreadyExists
	}

	if err := s.validate(entity, possession); err != nil {
		return errors.Wrap(err, "unable to validate")
	}

	s.possessionsMutex.Lock()
	s.possessions[entity] = possession
	s.possessionsMutex.Unlock()

	s.log.Info().EmbedObject(entity).EmbedObject(possession).Msg("Added possessions component.")

	return nil
}

func (s *PossessionSystem) Filter(entities []component.Entity, filter PossessionSystemFilterFn) []component.Entity {
	filteredEntities := []component.Entity{}

	for _, entity := range entities {
		s.possessionsMutex.Lock()
		possession, exists := s.possessions[entity]
		s.possessionsMutex.Unlock()
		if !exists {
			continue
		}

		if !filter(possession) {
			continue
		}

		entityCopy := entity
		filteredEntities = append(filteredEntities, entityCopy)
	}

	return filteredEntities
}

func (s *PossessionSystem) Get(entity component.Entity) (*component.Possession, error) {
	defer s.possessionsMutex.Unlock()
	s.possessionsMutex.Lock()

	possesion, exists := s.possessions[entity]
	if !exists {
		return nil, ErrPossessionComponentNotFound
	}

	return &possesion, nil
}

func (s *PossessionSystem) remove(entity component.Entity) error {
	possession, exists := s.possessions[entity]
	if !exists {
		return ErrPossessionComponentNotFound
	}

	s.possessionsMutex.Lock()
	delete(s.possessions, entity)
	s.possessionsMutex.Unlock()

	s.log.Info().EmbedObject(entity).EmbedObject(possession).Msg("Removed possessions component.")

	return nil
}

func (s *PossessionSystem) exists(entity component.Entity) bool {
	s.possessionsMutex.Lock()
	_, exists := s.possessions[entity]
	s.possessionsMutex.Unlock()

	return exists
}

func (s *PossessionSystem) applyDeltaTime(delta uint64) error {
	return nil
}

var (
	ErrPossessionAlreadyExists      = errors.New("possessions already hasPosition")
	ErrPossessionComponentNotFound  = errors.New("possessions component not found")
	ErrPossessionOwnerInvalidEntity = errors.New("possessions owner invalid component")
)
