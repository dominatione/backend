package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanetSystem_Count(t *testing.T) {
	state := NewState()

	planetCount := state.planet.Count()
	assert.EqualValues(t, 0, planetCount)
}
