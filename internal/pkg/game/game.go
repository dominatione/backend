package game

import (
	"context"
	"github.com/dominati-one/backend/internal/pkg/game/event"
	"github.com/dominati-one/backend/internal/pkg/game/world"
	"github.com/dominati-one/backend/internal/pkg/security"
	blockchainProtocol "github.com/dominati-one/backend/pkg/protocol/blockchain"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type Game struct {
	log        zerolog.Logger
	state      *world.State
	worldClock *world.WorldClock
}

func NewGame() *Game {
	worldClock := world.NewWorldClock(100)

	return &Game{
		log:        log.With().Str("applicationComponent", "game").Logger(),
		state:      world.NewState(),
		worldClock: worldClock,
	}
}

func (g *Game) SetCurrentTimestamp(timestamp uint64) error {
	delta, err := g.worldClock.SetCurrentTimestamp(timestamp)
	if err != nil {
		return errors.Wrap(err, "unable to set current timestamp on world clock")
	}

	startTime := time.Now()

	g.log.Trace().Msg("Delta time processing started.")

	err = g.state.ApplyDeltaTime(delta)
	if err != nil {
		return errors.Wrap(err, "unable to apply delta time on world state")
	}

	processingDuration := time.Now().Sub(startTime)

	g.log.Trace().Dur("processingDuration", processingDuration).Msg("Delta time processing finished.")

	return nil
}

func (g *Game) ApplyEvent(blockchainEvent *blockchainProtocol.Event) error {
	signature, err := security.NewSignature(blockchainEvent.Signature)
	if err != nil {
		return errors.Wrap(err, "unable to create signature")
	}

	if createPlanetEvent := blockchainEvent.Body.GetCreatePlanet(); createPlanetEvent != nil {
		return event.NewCreatePlanetHandler(g.state).Handle(createPlanetEvent, signature)
	}

	if createPlayerEvent := blockchainEvent.Body.GetCreatePlayer(); createPlayerEvent != nil {
		return event.NewCreatePlayerHandler(g.state).Handle(createPlayerEvent, signature)
	}

	return nil
}

func (g *Game) VerifyEvent(blockchainEvent *blockchainProtocol.Event) error {
	return g.Clone().ApplyEvent(blockchainEvent)
}

func (g *Game) Start(ctx context.Context) error {
	return nil
}

func (g *Game) State() *world.State {
	return g.state
}

func (g *Game) WorldClock() *world.WorldClock {
	return g.worldClock
}

func (g *Game) Clone() *Game {
	return &Game{
		log:        zerolog.Nop(),
		state:      g.state.Clone(),
		worldClock: g.worldClock.Clone(),
	}
}
