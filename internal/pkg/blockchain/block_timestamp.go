package blockchain

import (
	"time"
)

type BlockTimestamp struct {
	time.Time
}

func CreateBlockTimestampFromNow() BlockTimestamp {
	return BlockTimestamp{
		Time: time.Now(),
	}
}

func CreateBlockTimestampFromUnixMilliseconds(unixMilliseconds uint64) BlockTimestamp {
	nanoseconds := unixMilliseconds * 1000000

	return BlockTimestamp{
		Time: time.Unix(0, int64(nanoseconds)),
	}
}

func (t BlockTimestamp) Add(d time.Duration) BlockTimestamp {
	return BlockTimestamp{
		Time: t.Time.Add(d),
	}
}

func (t BlockTimestamp) Raw() time.Time {
	return t.Time
}

func (t BlockTimestamp) UnixMilliseconds() uint64 {
	return uint64(t.UnixNano() / int64(time.Millisecond))
}

var (
	EmptyBlockTimestamp = BlockTimestamp{}
)
