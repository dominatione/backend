package component

import (
	"github.com/dominati-one/backend/pkg/protocol/component"
	"github.com/rs/zerolog"
)

type AreaPositionLayer uint8

const (
	AreaPositionLayerEmpty = iota
	AreaPositionLayerSurface
	AreaPositionLayerPlayer
)

type AreaPosition struct {
	Entity Entity
	Layer  AreaPositionLayer
	X      uint32
	Y      uint32
	Width  uint8
	Height uint8
}

func (c AreaPosition) Protobuf() *component.AreaPosition {
	return &component.AreaPosition{
		Entity: uint64(c.Entity),
		X:      uint32(c.X),
		Y:      uint32(c.Y),
		Width:  uint32(c.Width),
		Height: uint32(c.Height),
	}
}

func (s AreaPosition) MarshalZerologObject(e *zerolog.Event) {
	e.Str("areaPositionEntity", s.Entity.String())
	e.Uint32("areaPositionX", s.X)
	e.Uint32("areaPositionY", s.Y)
	e.Uint8("areaPositionWidth", s.Width)
	e.Uint8("areaPositionHeight", s.Height)
}
