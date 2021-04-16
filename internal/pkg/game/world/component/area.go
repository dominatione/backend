package component

import (
	"github.com/dominati-one/backend/pkg/protocol/component"
	"github.com/rs/zerolog"
)

type Area struct {
	Width  uint32
	Height uint32
}

func (c Area) Protobuf() *component.Area {
	return &component.Area{
		Width:  c.Width,
		Height: c.Height,
	}
}

func (c Area) MarshalZerologObject(e *zerolog.Event) {
	e.Uint32("areaWidth", c.Width)
	e.Uint32("areaHeight", c.Height)
}
