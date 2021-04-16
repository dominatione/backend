package backend

import (
	"context"
	"github.com/dominati-one/backend/internal/pkg/blockchain"
	"github.com/dominati-one/backend/internal/pkg/game"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type EventPump struct {
	log          zerolog.Logger
	eventEmitter *blockchain.EventEmitter
	game         *game.Game
	applyMutex   sync.Mutex
}

func NewEventPump(eventEmitter *blockchain.EventEmitter, game *game.Game) *EventPump {
	return &EventPump{
		log:          log.With().Str("applicationComponent", "eventPump").Logger(),
		eventEmitter: eventEmitter,
		game:         game,
	}
}

func (p *EventPump) Start(ctx context.Context) error {
	go p.applyCurrentTimeLoop(ctx)
	go p.applyEventLoop(ctx)

	p.log.Info().Msg("Started.")

	return nil
}

func (p *EventPump) applyEventLoop(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		event, err := p.eventEmitter.WaitForEvent(ctx)
		if err != nil {
			p.log.Error().Err(err).Msg("Waiting for event failed.")
			continue
		}

		p.applyMutex.Lock()

		if err = p.game.ApplyEvent(event); err != nil {
			p.log.Panic().Err(err).Msg("Unable to apply event to game.")
			return
		}

		p.applyMutex.Unlock()

	}
}

func (p *EventPump) applyCurrentTimeLoop(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		block, err := p.eventEmitter.WaitForBlock(ctx)
		if err != nil {
			p.log.Error().Err(err).Msg("Waiting for block failed.")
			continue
		}

		p.applyMutex.Lock()

		if err = p.game.SetCurrentTimestamp(block.Body.Timestamp); err != nil {
			p.log.Panic().Err(err).Msg("Unable to apply current time to game.")
			return
		}

		p.applyMutex.Unlock()

	}

}
