package strategy

import (
	"github.com/lugobots/arena"
	"github.com/lugobots/arena/physics"
	"github.com/lugobots/arena/units"
	"math"
)

// RegionHeight defines the height of a region
const RegionHeight = units.FieldHeight / 4

// RegionWidth defines the width of a region
const RegionWidth = units.FieldWidth / 8

// PlayerRegion defines a region based on the left-bottom and top-right coordinates of the region
type PlayerRegion struct {
	CornerA physics.Point
	CornerB physics.Point
}

// RegionCode is a cartesian coordinate based on the region division of the court
// The region code is always based on the home team perspective, so the coordinates must be mirrored to apply the
// same location in the field to the away team
type RegionCode struct {
	X int
	Y int
}

// Center finds the central coordinates of the region
func (r RegionCode) Center(place arena.TeamPlace) physics.Point {
	center := physics.Point{
		PosX: (r.X * RegionWidth) + (RegionWidth / 2),
		PosY: (r.Y * RegionHeight) + (RegionHeight / 2),
	}
	if place == arena.AwayTeam {
		center = MirrorCoordsToAway(center)
	}
	return center
}

// ForwardRightCorner finds the point of the region at the right edge that is closer to the attack field
func (r RegionCode) ForwardRightCorner(place arena.TeamPlace) physics.Point {
	fr := physics.Point{
		PosX: (r.X + 1) * RegionWidth,
		PosY: r.Y * RegionHeight,
	}
	if place == arena.AwayTeam {
		fr = MirrorCoordsToAway(fr)
	}
	return fr
}

// ForwardLeftCorner finds the point of the region at the left edge that is closer to the attack field
func (r RegionCode) ForwardLeftCorner(place arena.TeamPlace) physics.Point {
	fl := physics.Point{
		PosX: (r.X + 1) * RegionWidth,
		PosY: (r.Y + 1) * RegionHeight,
	}
	if place == arena.AwayTeam {
		fl = MirrorCoordsToAway(fl)
	}
	return fl
}

// Forwards finds the next region towards to the attack field, or return itself when there is no region in front of it
func (r RegionCode) Forwards() RegionCode {
	if r.X == 7 {
		return r
	}
	return RegionCode{
		X: r.X + 1,
		Y: r.Y,
	}
}

// Backwards finds the next region towards to the defense field, or return itself when there is no region in behind of it
func (r RegionCode) Backwards() RegionCode {
	if r.X == 0 {
		return r
	}
	return RegionCode{
		X: r.X - 1,
		Y: r.Y,
	}
}

// Left finds the region in the left side of this region, or return itself when there is no region there
func (r RegionCode) Left() RegionCode {
	if r.Y == 3 {
		return r
	}
	return RegionCode{
		X: r.X,
		Y: r.Y + 1,
	}
}

// Right finds the region in the right side of this region, or return itself when there is no region there
func (r RegionCode) Right() RegionCode {
	if r.Y == 0 {
		return r
	}
	return RegionCode{
		X: r.X,
		Y: r.Y - 1,
	}
}

// ChessDistanceTo calculates what is the chess distance (steps towards any direction, even diagonal) between these regions
func (r RegionCode) ChessDistanceTo(b RegionCode) int {
	return int(math.Max(
		math.Abs(float64(r.X-b.X)),
		math.Abs(float64(r.Y-b.Y)),
	))
}

func GetRegionCode(a physics.Point, place arena.TeamPlace) RegionCode {
	if place == arena.AwayTeam {
		a = MirrorCoordsToAway(a)
	}
	cx := float64(a.PosX / RegionWidth)
	cy := float64(a.PosY / RegionHeight)
	return RegionCode{
		X: int(math.Min(cx, 7)),
		Y: int(math.Min(cy, 3)),
	}
}

// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field is based on the bottom left corner!
func MirrorCoordsToAway(coords physics.Point) physics.Point {
	return physics.Point{
		PosX: units.FieldWidth - coords.PosX,
		PosY: units.FieldHeight - coords.PosY,
	}
}
