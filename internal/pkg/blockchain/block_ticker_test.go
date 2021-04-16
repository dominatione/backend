package blockchain

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlockTicker_WaitForNext(t *testing.T) {
	var err error
	var blockTimestamp *BlockTimestamp

	blockTicker := NewBlockTicker(NetworkSettings{BlockInterval: time.Second * 2})

	timeoutCtx, _ := context.WithTimeout(context.TODO(), time.Millisecond)

	blockTimestamp, err = blockTicker.WaitForNext(timeoutCtx, CreateBlockTimestampFromNow())
	assert.ErrorIs(t, err, ErrCanceledBlockTimestampWait)
	assert.Nil(t, blockTimestamp)
}
