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

func (r RegionCode) Center(place Units.TeamPlace) Physics.Point  {
	center := Physics.Point{
		PosX: (r.X * RegionWidth) + (RegionWidth/2),
		PosY: (r.Y * RegionHeight) + (RegionHeight/2),
	}
	if place == Units.AwayTeam {
		center = MirrorCoordsToAway(center)
	}
	return center
}

func (r RegionCode) ForwardRightCorner(place Units.TeamPlace) Physics.Point  {
	fr := Physics.Point{
		PosX: (r.X + 1) * RegionWidth,
		PosY: r.Y * RegionHeight,
	}
	if place == Units.AwayTeam {
		fr = MirrorCoordsToAway(fr)
	}
	return fr
}

func (r RegionCode) ForwardLeftCorner(place Units.TeamPlace) Physics.Point  {
	fl := Physics.Point{
		PosX: (r.X + 1) * RegionWidth,
		PosY: (r.Y + 1) * RegionHeight,
	}
	if place == Units.AwayTeam {
		fl = MirrorCoordsToAway(fl)
	}
	return fl
}

func (r RegionCode) Forwards() RegionCode  {
	if r.X == 7 {
		return r
	}
	return RegionCode{
		X: r.X + 1,
		Y: r.Y,
	}
}

func (r RegionCode) Backwards() RegionCode  {
	if r.X == 0 {
		return r
	}
	return RegionCode{
		X: r.X - 1,
		Y: r.Y,
	}
}

func (r RegionCode) Left() RegionCode  {
	if r.Y == 3 {
		return r
	}
	return RegionCode{
		X: r.X,
		Y: r.Y + 1,
	}
}

func (r RegionCode) Right() RegionCode  {
	if r.Y == 0 {
		return r
	}
	return RegionCode{
		X: r.X,
		Y: r.Y - 1,
	}
}

func (r RegionCode) ChessDistanceTo(b RegionCode) int {
	return int(math.Max(
		math.Abs(float64(r.X - b.X)),
		math.Abs(float64(r.Y - b.Y)),
	))
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


// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field is based on the bottom left corner!
func MirrorCoordsToAway(coords Physics.Point) Physics.Point {
	return Physics.Point{
		PosX: Units.CourtWidth - coords.PosX,
		PosY: Units.CourtHeight - coords.PosY,
	}
}