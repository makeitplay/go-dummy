package strategy

import (
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"
	"math"
)
const RegionHeight = Units.CourtHeight / 4
const RegionWidth = Units.CourtWidth / 8

type PlayerRegion struct {
	CornerA Physics.Point
	CornerB Physics.Point
}

type RegionCode struct {
	X int
	Y int
}

func GetRegionCode(a Physics.Point, place Units.TeamPlace) RegionCode {
	if place == Units.AwayTeam {
		a = MirrorCoordsToAway(a)
	}
	cx := float64(a.PosX / RegionWidth)
	cy := float64(a.PosY / RegionHeight)
	return RegionCode{
		X: int(math.Min(cx, 7)),
		Y: int(math.Min(cy, 3)),
	}
}

func GetRegionCenter(code RegionCode, place Units.TeamPlace) Physics.Point {
	if code.X > 7 || code.Y > 3 {
		panic("invalid region code")
	}
	center := Physics.Point{
		PosX: (code.X * RegionWidth) + (RegionWidth/2),
		PosY: (code.Y * RegionHeight) + (RegionHeight/2),
	}
	if place == Units.AwayTeam {
		center = MirrorCoordsToAway(center)
	}
	return center
}

// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field is based on the bottom left corner!
func MirrorCoordsToAway(coords Physics.Point) Physics.Point {
	return Physics.Point{
		PosX: Units.CourtWidth - coords.PosX,
		PosY: Units.CourtHeight - coords.PosY,
	}
}