package coach

import (
	"github.com/lugobots/client-player-go/v2/proto"
)

type Positioner interface {
	GetRegion(col, row uint8) (Region, error)
	GetPointRegion(point proto.Point) (Region, error)
}

type Region interface {
	Col() uint8
	Row() uint8
	Center() proto.Point
}
