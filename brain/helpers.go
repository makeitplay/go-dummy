package brain

import (
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/client-player-go/Game"
)


func BallMaxSafePassDistance(Speed float64) float64 {
	return Speed + (Units.BallDeceleration)/2
}

// Opponent id and angle between it and the target
func watchOpponentOnMyRoute(player *Game.Player, target Physics.Point) map[int]float64 {
	opponentTeam := player.GetOpponentTeam(player.LastMsg.GameInfo)
	opponents := make(map[int]float64)
	vectorExpected := Physics.NewVector(player.Coords, target)
	for _, opponent := range opponentTeam.Players {
		collisionPoint := opponent.VectorCollides(*vectorExpected, player.Coords, float64(player.Size))

		if collisionPoint != nil {
			opponents[opponent.Id] = collisionPoint.DistanceTo(player.Coords)
		}
	}
	return opponents
}