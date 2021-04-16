package world

import (
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	"github.com/ojrac/opensimplex-go"
	"github.com/pkg/errors"
	"math"
	"math/rand"
)

type planetGenerator struct {
	firstOctave  opensimplex.Noise
	secondOctave opensimplex.Noise
	thirdOctave  opensimplex.Noise
	frequency    float64
	width        uint32
	height       uint32
}

type planetGenerators struct {
	surfaceGenerator planetGenerator
	fertileGenerator planetGenerator
	stoneGenerator   planetGenerator
	gravelGenerator  planetGenerator
}

type PlanetActions struct {
	state *State
}

func newPlanetActions(state *State) *PlanetActions {
	return &PlanetActions{
		state: state,
	}
}

func newPlanetGenerator(seed int64, frequency float64, width, height uint32) planetGenerator {
	return planetGenerator{
		firstOctave:  opensimplex.NewNormalized(seed * int64(frequency)),
		secondOctave: opensimplex.NewNormalized(seed * int64(frequency) * 2),
		thirdOctave:  opensimplex.NewNormalized(seed * int64(frequency) * 4),
		frequency:    frequency,
		width:        width,
		height:       height,
	}
}

func (g *planetGenerator) get(x, y uint32) float64 {
	xFloat := float64(x) / float64(g.width)
	yFloat := float64(y) / float64(g.height)

	value := 0.0

	value += 1 * g.firstOctave.Eval2(xFloat*g.frequency, yFloat*g.frequency)
	value += 0.5 * g.secondOctave.Eval2(xFloat*4*g.frequency, yFloat*4*g.frequency)
	value += 0.25 * g.thirdOctave.Eval2(xFloat*8*g.frequency, yFloat*8*g.frequency)

	return value / 1.75
}

func (f *PlanetActions) Create() (*component.Entity, error) {
	planetEntity := f.state.Create(component.EntityKindPlanet)

	seed := int64(f.state.Planet().Count() + 1)

	width, height := f.createDimensionsFromSeed(seed)

	name, err := f.createNameFromSeed(seed)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create planetComponent name from seed")
	}

	planetComponent := component.Planet{
		Seed: seed,
		Name: name,
	}

	areaComponent := component.Area{
		Width:  width,
		Height: height,
	}

	areaTiles := f.createAreaTilesFromSeed(width, height, seed, planetEntity)

	if err := f.state.planet.add(planetEntity, planetComponent); err != nil {
		return nil, errors.Wrap(err, "planet component add to planet system failed")
	}

	if err := f.state.area.addArea(planetEntity, areaComponent, areaTiles); err != nil {
		return nil, errors.Wrap(err, "area component add to area system failed")
	}

	return &planetEntity, nil
}

func (f *PlanetActions) createDimensionsFromSeed(seed int64) (uint32, uint32) {
	source := rand.New(rand.NewSource(seed))

	width := uint32(1000 + (source.Float64() * 4000))
	height := uint32(1000 + (source.Float64() * 4000))

	return width, height
}

func (f *PlanetActions) createNameFromSeed(seed int64) (string, error) {
	names := []string{
		"New Ganymede",
		"Tatlon",
		"Aertan",
		"New Kenya",
		"Satai",
		"Callisto",
		"9733 Sagittae III",
		"Ru-Shou Prime",
		"New Earth",
	}

	if seed < 1 || seed > int64(len(names)) {
		return "", ErrPlanetNameOutOfBounds
	}

	return names[seed-1], nil
}

func (f *PlanetActions) createAreaTilesFromSeed(width, height uint32, seed int64, planetEntity component.Entity) component.AreaTiles {
	generators := &planetGenerators{
		surfaceGenerator: newPlanetGenerator(seed, 8.0, width, height),
		fertileGenerator: newPlanetGenerator(seed, 16.0, width, height),
		stoneGenerator:   newPlanetGenerator(seed, 32.0, width, height),
	}

	tiles := make(component.AreaTiles, width*height)

	var x, y uint32

	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			index := x + (y * width)

			areaTileKind := f.createSurfaceFromPosition(x, y, width, height, generators)

			tiles[index] = component.AreaTile{
				Kind:        areaTileKind,
				OwnerEntity: planetEntity,
			}
		}
	}

	return tiles
}

func (f *PlanetActions) createSurfaceFromPosition(x, y, width, height uint32, generators *planetGenerators) component.AreaTileKind {
	shallowWaterLevel := 0.30
	waterLevel := 0.40
	sandLevel := 0.43
	stoneLevel := 0.80

	wrapDistance := 0.85

	surfaceLevel := generators.surfaceGenerator.get(x, y)

	xDistance := math.Abs(float64(width/2)-float64(x)) / float64(width/2)
	yDistance := math.Abs(float64(height/2)-float64(y)) / float64(height/2)

	distance := math.Max(xDistance, yDistance)

	if distance > wrapDistance {
		deepLevelModifier := (distance - wrapDistance) / (1 - wrapDistance)
		surfaceLevel -= deepLevelModifier
		if surfaceLevel < 0 {
			surfaceLevel = 0
		}
	}

	if surfaceLevel < shallowWaterLevel {
		return component.AreaTileKindShallowWater
	} else if surfaceLevel < waterLevel {
		return component.AreaTileKindWater
	} else if surfaceLevel < sandLevel {
		return component.AreaTileKindSand
	} else if surfaceLevel > stoneLevel {
		if generators.stoneGenerator.get(x, y) > 0.7 {
			return component.AreaTileKindLava
		} else {
			return component.AreaTileKindStone
		}
	} else {
		if generators.fertileGenerator.get(x, y) > 0.25 {
			return component.AreaTileKindGround
		} else {
			return component.AreaTileKindFertileGround
		}

	}

}

var (
	ErrPlanetNameOutOfBounds = errors.New("planet name out of bounds")
)
