package blockchain

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type BlockTicker struct {
	settings NetworkSettings
}

func NewBlockTicker(settings NetworkSettings) *BlockTicker {
	return &BlockTicker{
		settings: settings,
	}
}

func (t *BlockTicker) WaitForNext(ctx context.Context, previousBlockTimestamp BlockTimestamp) (*BlockTimestamp, error) {
	nextBlockTimestamp := t.getNext(previousBlockTimestamp).
		Add(time.Millisecond * 50)

	waitDuration := nextBlockTimestamp.Sub(time.Now())

	select {
	case <-ctx.Done():
		return nil, ErrCanceledBlockTimestampWait
	case <-time.NewTimer(waitDuration).C:
	}

	return &nextBlockTimestamp, nil
}

func (t *BlockTicker) getNext(previousBlockTimestamp BlockTimestamp) BlockTimestamp {
	nextBlockTimestamp := previousBlockTimestamp.Add(t.settings.BlockInterval)

	if nextBlockTimestamp.Before(time.Now()) {
		return CreateBlockTimestampFromNow()
	}

	return nextBlockTimestamp
}

var (
	ErrCanceledBlockTimestampWait = errors.New("canceled block timestamp wait")
)
