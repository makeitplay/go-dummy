package coach

import (
	"github.com/lugobots/client-player-go/v2/proto"
)

// Positioner Helps the bots to see the fields from their team perspective instead of using the cartesian plan provided
// by the game server. Instead of base your logic on the axes X and Y, the positioner create a region map based
// on the team side.
// The region coordinates uses the defensive field's right corner as its origin.
// This mechanism if specially useful to define players regions based on their roles, since you do not have to mirror
// the coordinate, neither do extra logic to define regions on the field where the player should be.
type Positioner interface {
	// GetRegion Returns a region based on the coordinates and on the current field division
	GetRegion(col, row uint8) (Region, error)
	// GetPointRegion returns the region where that point is in
	GetPointRegion(point proto.Point) (Region, error)
}

// Region represent a quadrant on the field. It is not always squared form because you may define how many cols/rows
// the field will be divided in. So, based on that division (e.g. 4 rows, 6 cols) there will be a fixed number of regions
// and their coordinates will be zero-index (e.g. from 0 to 3 rows when divided in 4 rows).
type Region interface {
	// The col coordinate based on the field division
	Col() uint8
	// The row coordinate based on the field division
	Row() uint8
	// Return the point at the center of the quadrant represented by this Region. It is not always precise.
	Center() proto.Point
}
