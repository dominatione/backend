package world

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	SecondDeltaFactor    float32 = 1.0 / (1.0)
	MinuteDeltaFactor            = 1.0 / (60.0)
	HourDeltaFactor              = 1.0 / (60.0 * 60.0)
	HalfDayDeltaFactor           = 1.0 / (12.0 * 60.0 * 60.0)
	DayDeltaFactor               = 1.0 / (1.0 * 24.0 * 60.0 * 60.0)
	TwoDaysDeltaFactor           = 1.0 / (2.0 * 24.0 * 60.0 * 60.0)
	ThreeDaysDeltaFactor         = 1.0 / (3.0 * 24.0 * 60.0 * 60.0)
	FourDaysDeltaFactor          = 1.0 / (4.0 * 24.0 * 60.0 * 60.0)
	WeekDeltaFactor              = 1.0 / (7.0 * 24.0 * 60.0 * 60.0)
)

type WorldClock struct {
	log                zerolog.Logger
	firstTickTimestamp uint64
	lastTickTimestamp  uint64
	compression        uint64
}

func NewWorldClock(compression uint64) *WorldClock {
	return &WorldClock{
		log:                log.With().Str("applicationComponent", "game").Str("gameComponent", "worldClock").Logger(),
		firstTickTimestamp: 0,
		lastTickTimestamp:  0,
		compression:        compression,
	}
}

func (c *WorldClock) Clone() *WorldClock {
	return &WorldClock{
		log:                zerolog.Nop(),
		firstTickTimestamp: c.firstTickTimestamp,
		lastTickTimestamp:  c.lastTickTimestamp,
		compression:        c.compression,
	}
}

func (c *WorldClock) SetCurrentTimestamp(timestamp uint64) (uint64, error) {
	if c.firstTickTimestamp == 0 && c.lastTickTimestamp == 0 {
		c.firstTickTimestamp = timestamp
		c.lastTickTimestamp = timestamp
		c.log.Info().Time("currentTime", c.Time()).Msg("Received first timestamp.")
		return 0, nil
	}

	if c.lastTickTimestamp > timestamp {
		return 0, ErrCurrentTimeLessThanLastTimeEvent
	}

	delta := (timestamp - c.lastTickTimestamp) * c.compression

	c.lastTickTimestamp = timestamp

	c.log.Info().
		Time("currentTime", c.Time()).
		Uint64("deltaTime", delta).
		Msg("Clock tick.")

	return delta, nil
}

func (c *WorldClock) Time() time.Time {
	delta := time.Duration((c.lastTickTimestamp-c.firstTickTimestamp)/1000) * time.Second * time.Duration(c.compression)

	return time.Time{}.Add(delta)
}

var (
	ErrCurrentTimeLessThanLastTimeEvent = errors.New("current time less than last time event")
)
