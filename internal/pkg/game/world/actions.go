package world

type Actions struct {
	planet *PlanetActions
	seed   *SeedActions
	plant  *PlantActions
}

func newActions(state *State) *Actions {
	return &Actions{
		planet: newPlanetActions(state),
		seed:   newSeedActions(state),
		plant:  newPlantActions(state),
	}
}

func (a *Actions) Planet() *PlanetActions {
	return a.planet
}

func (a *Actions) Seed() *SeedActions {
	return a.seed
}
