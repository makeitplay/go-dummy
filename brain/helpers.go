package brain

import (
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/client-player-go/Game"
)

func BallMaxSafePassDistance(Speed float64) float64 {
	return Speed + (Units.BallDeceleration)/2
}

type PointCollection []Physics.Point

func (s PointCollection) Len() int {
	return len(s)
}

func (s PointCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SortByDistance struct {
	PointCollection
	From    Physics.Point
}

func (s SortByDistance) Less(i, j int) bool {
	return s.From.DistanceTo(s.PointCollection[i]) < s.From.DistanceTo(s.PointCollection[j])
}

// watchOpponentOnMyRoute returns a list for obstacle between the player an it's target sorted by the distance to it
func watchOpponentOnMyRoute(player *Game.Player, target Physics.Point) PointCollection {
	opponentTeam := player.GetOpponentTeam(player.LastMsg.GameInfo)
	cllisionPoints := SortByDistance{From: player.Coords}

	vectorExpected := Physics.NewVector(player.Coords, target)
	for _, opponent := range opponentTeam.Players {
		collisionPoint := opponent.VectorCollides(*vectorExpected, player.Coords, float64(player.Size)/2)

		if collisionPoint != nil {
			cllisionPoints.PointCollection = append(cllisionPoints.PointCollection, *collisionPoint)
		}
	}
	return cllisionPoints.PointCollection
}
