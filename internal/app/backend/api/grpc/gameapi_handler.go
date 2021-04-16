package grpc

import (
	"context"
	"github.com/dominati-one/backend/internal/pkg/blockchain"
	"github.com/dominati-one/backend/internal/pkg/game"
	"github.com/dominati-one/backend/internal/pkg/game/world"
	"github.com/dominati-one/backend/internal/pkg/game/world/component"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	protocolComponent "github.com/dominati-one/backend/pkg/protocol/component"
	protocolEntity "github.com/dominati-one/backend/pkg/protocol/entity"
	"github.com/dominati-one/backend/pkg/protocol/gameapi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type GameApiHandler struct {
	log          zerolog.Logger
	eventBacklog *blockchain.LocalEventBacklog
	game         *game.Game
}

func NewGameApiHandler(game *game.Game, eventBacklog *blockchain.LocalEventBacklog) *GameApiHandler {
	return &GameApiHandler{
		log:          log.With().Str("applicationComponent", "gameApiHandler").Logger(),
		game:         game,
		eventBacklog: eventBacklog,
	}
}

func (h *GameApiHandler) CreatePlanet(ctx context.Context, request *gameapi.CreatePlanetRequest) (*gameapi.CreatePlanetResponse, error) {
	createPlanetEvent := &blockchainProtocol.EventCreatePlanet{}

	eventId, err := h.eventBacklog.Add(createPlanetEvent)
	if err != nil {
		return nil, errors.Wrap(err, "unable to add event to backlog")
	}

	return &gameapi.CreatePlanetResponse{EventId: eventId.Bytes()}, nil
}

func (h *GameApiHandler) GetPlanet(ctx context.Context, request *gameapi.GetPlanetRequest) (*gameapi.GetPlanetResponse, error) {
	planetEntity := component.Entity(request.Entity)
	planet, err := h.game.State().Planet().Get(planetEntity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get planet")
	}

	area, err := h.game.State().Area().GetArea(planetEntity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get area")
	}

	return &gameapi.GetPlanetResponse{
		Planet: &protocolEntity.Planet{
			Entity: request.Entity,
			Planet: planet.Protobuf(),
			Area:   area.Protobuf(),
		},
	}, nil
}

func (h *GameApiHandler) GetAreaTiles(ctx context.Context, request *gameapi.GetAreaTilesRequest) (*gameapi.GetAreaTilesResponse, error) {
	areaEntity := component.Entity(request.Entity)
	areaTiles, err := h.game.State().Area().GetAreaTiles(areaEntity, world.AreaTilesExtent{
		Left:   request.Left,
		Top:    request.Top,
		Right:  request.Right,
		Bottom: request.Bottom,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to get area tiles")
	}

	responseTiles := make([]*protocolComponent.AreaTile, len(areaTiles))

	for tileIndex, areaTile := range areaTiles {
		responseTiles[tileIndex] = &protocolComponent.AreaTile{
			Kind:        areaTile.Kind.Protobuf(),
			OwnerEntity: uint64(areaTile.OwnerEntity),
		}
	}

	return &gameapi.GetAreaTilesResponse{
		TilesCount: uint64(len(responseTiles)),
		Tiles:      responseTiles,
	}, nil
}

func (h *GameApiHandler) GetSeeds(ctx context.Context, request *gameapi.GetSeedsRequest) (*gameapi.GetSeedsResponse, error) {
	entities := h.game.State().Seed().Entities()

	if request.QueryParams != nil {
		if request.QueryParams.Owner != nil {
			entities = h.game.State().Possession().Filter(entities, func(possession component.Possession) bool {
				return uint64(possession.OwnerEntity) == request.QueryParams.Owner.OwnerEntity
			})
		}
	}

	seeds := []*protocolEntity.Seed{}

	for _, seedEntity := range entities {
		seedComponent, err := h.game.State().Seed().Get(seedEntity)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get seed component")
		}

		areaPositionComponent, _ := h.game.State().Area().GetPosition(seedEntity)
		possessionComponent, _ := h.game.State().Possession().Get(seedEntity)

		seedItem := &protocolEntity.Seed{
			Entity: uint64(seedEntity),
			Seed:   seedComponent.Protobuf(),
		}

		if areaPositionComponent != nil {
			seedItem.AreaPosition = areaPositionComponent.Protobuf()
		}
		if possessionComponent != nil {
			seedItem.Possession = possessionComponent.Protobuf()
		}

		seeds = append(seeds, seedItem)
	}

	return &gameapi.GetSeedsResponse{Seeds: seeds}, nil
}

func (h *GameApiHandler) GetPlanets(ctx context.Context, request *gameapi.GetPlanetsRequest) (*gameapi.GetPlanetsResponse, error) {
	entities := h.game.State().Planet().Entities()

	planets := []*protocolEntity.Planet{}

	for _, planetEntity := range entities {
		planetComponent, err := h.game.State().Planet().Get(planetEntity)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get planet")
		}

		areaComponent, err := h.game.State().Area().GetArea(planetEntity)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get area")
		}

		planets = append(planets, &protocolEntity.Planet{
			Entity: uint64(planetEntity),
			Planet: planetComponent.Protobuf(),
			Area:   areaComponent.Protobuf(),
		})
	}

	return &gameapi.GetPlanetsResponse{
		Planets: planets,
	}, nil
}
