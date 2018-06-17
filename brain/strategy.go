package brain

import (
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/client-player-go/tatics"
	"github.com/makeitplay/commons/Units"
)

func (b *Brain) isItInMyRegion(coords Physics.Point) bool {
	myRagion := b.myRegion()
	isInX := coords.PosX >= myRagion.CornerA.PosX && coords.PosX <= myRagion.CornerB.PosX
	isInY := coords.PosY >= myRagion.CornerA.PosY && coords.PosY <= myRagion.CornerB.PosY
	return isInX && isInY
}

func (b *Brain) myRegionCenter() Physics.Point {
	myRegiao := b.myRegion()
	//regionDiagonal := math.Hypot(float64(myRegiao.CornerA.PosX), float64(myRegiao.CornerB.PosY))
	halfXDistance := (myRegiao.CornerB.PosX - myRegiao.CornerA.PosX) / 2
	halfYDistance := (myRegiao.CornerB.PosY - myRegiao.CornerA.PosY) / 2
	return Physics.Point{
		PosX: int(myRegiao.CornerA.PosX + halfXDistance),
		PosY: int(myRegiao.CornerA.PosY + halfYDistance),
	}
}

func (b *Brain) myRegion() tatics.PlayerRegion {
	myRagion := tatics.HomePlayersRegions[b.Number]
	if b.TeamPlace == Units.AwayTeam {
		myRagion = MirrorRegion(myRagion)
	}
	return myRagion
}
func MirrorRegion(region tatics.PlayerRegion) tatics.PlayerRegion {
	return tatics.PlayerRegion{
		CornerA: tatics.MirrorCoordToAway(region.CornerA), // have to switch the corner because the convention for Regions
		CornerB: tatics.MirrorCoordToAway(region.CornerB),
	}
}
