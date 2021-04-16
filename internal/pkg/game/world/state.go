package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/pkg/errors"
	"sync"
)

type State struct {
	freeEntityId  uint64
	entitiesMutex sync.Mutex
	entities      map[component.Entity]component.EntityKind
	actions       *Actions

	area       *AreaSystem
	seed       *SeedSystem
	plant      *PlantSystem
	planet     *PlanetSystem
	possession *PossessionSystem
}

func NewState() *State {
	state := &State{
		freeEntityId: 1,
		entities:     map[component.Entity]component.EntityKind{},
	}

	state.area = NewAreaSystem(state)
	state.seed = newSeedSystem(state)
	state.plant = newPlantSystem(state)
	state.planet = newPlanetSystem(state)
	state.possession = newPossessionSystem(state)

	state.actions = newActions(state)

	return state
}

func (m *State) Clone() *State {
	entitiesClone := map[component.Entity]component.EntityKind{}

	m.entitiesMutex.Lock()
	for entity, entityKind := range m.entities {
		entitiesClone[entity] = entityKind
	}
	m.entitiesMutex.Unlock()

	stateClone := &State{
		freeEntityId: m.freeEntityId,
		entities:     entitiesClone,
	}

	stateClone.area = m.area.Clone(stateClone)
	stateClone.seed = m.seed.clone(stateClone)
	stateClone.plant = m.plant.clone(stateClone)
	stateClone.planet = m.planet.clone(stateClone)
	stateClone.possession = m.possession.clone(stateClone)

	stateClone.actions = newActions(stateClone)

	return stateClone
}

func (m *State) Create(kind component.EntityKind) component.Entity {
	defer m.entitiesMutex.Unlock()

	m.entitiesMutex.Lock()

	entity := component.Entity(m.freeEntityId)
	m.freeEntityId++

	m.entities[entity] = kind

	return entity
}

func (m *State) GetKind(entity component.Entity) (*component.EntityKind, error) {
	defer m.entitiesMutex.Unlock()

	m.entitiesMutex.Lock()
	kind, exists := m.entities[entity]
	if !exists {
		return nil, ErrEntityNotExists
	}

	return &kind, nil
}

func (m *State) Exists(entity component.Entity) bool {
	defer m.entitiesMutex.Unlock()

	m.entitiesMutex.Lock()
	_, exists := m.entities[entity]

	return exists
}

func (m *State) Remove(entity component.Entity) error {
	if !m.Exists(entity) {
		return ErrEntityNotFound
	}

	if m.area.hasPosition(entity) {
		if err := m.area.removePosition(entity); err != nil {
			return errors.Wrap(err, "unable to remove area position components from area system")
		}
	}

	if m.area.hasArea(entity) {
		if err := m.area.removeArea(entity); err != nil {
			return errors.Wrap(err, "unable to remove area components from area system")
		}
	}

	if m.seed.exists(entity) {
		if err := m.seed.remove(entity); err != nil {
			return errors.Wrap(err, "unable to removePosition components from seed system")
		}
	}

	if m.plant.exists(entity) {
		if err := m.plant.remove(entity); err != nil {
			return errors.Wrap(err, "unable to removePosition components from plant system")
		}
	}

	if m.planet.exists(entity) {
		if err := m.planet.remove(entity); err != nil {
			return errors.Wrap(err, "unable to removePosition components from planet system")
		}
	}

	if m.possession.exists(entity) {
		if err := m.possession.remove(entity); err != nil {
			return errors.Wrap(err, "unable to removePosition components from possessions system")
		}
	}

	m.entitiesMutex.Lock()
	delete(m.entities, entity)
	m.entitiesMutex.Unlock()

	return nil
}

func (m *State) ApplyDeltaTime(delta uint64) error {
	if err := m.area.applyDeltaTime(delta); err != nil {
		return errors.Wrap(err, "unable to apply delta time on area system")
	}

	if err := m.seed.applyDeltaTime(delta); err != nil {
		return errors.Wrap(err, "unable to apply delta time on seed system")
	}

	if err := m.plant.applyDeltaTime(delta); err != nil {
		return errors.Wrap(err, "unable to apply delta time on plant system")
	}

	if err := m.planet.applyDeltaTime(delta); err != nil {
		return errors.Wrap(err, "unable to apply delta time on planet system")
	}

	if err := m.possession.applyDeltaTime(delta); err != nil {
		return errors.Wrap(err, "unable to apply delta time on possessions system")
	}

	return nil
}

func (m *State) Actions() *Actions {
	return m.actions
}

func (m *State) Area() *AreaSystem {
	return m.area
}

func (m *State) Seed() *SeedSystem {
	return m.seed
}

func (m *State) Plant() *PlantSystem {
	return m.Plant()
}

func (m *State) Planet() *PlanetSystem {
	return m.planet
}

func (m *State) Possession() *PossessionSystem {
	return m.possession
}

var (
	ErrEntityNotExists = errors.New("component not hasPosition")
)
