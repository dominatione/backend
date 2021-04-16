package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFromBlockTimestampUnixMilliseconds(t *testing.T) {
	blockTimestamp := CreateBlockTimestampFromUnixMilliseconds(1)
	unixMilliseconds := blockTimestamp.UnixMilliseconds()
	assert.EqualValues(t, 1, unixMilliseconds)
}
