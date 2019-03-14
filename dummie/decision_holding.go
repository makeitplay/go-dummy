package dummie

import (
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"math"
	"github.com/makeitplay/client-player-go"
)


func (d *Dummie) orderForHoldingTheBall() (msg string, orders []orders.Order) {
	player := d.Player

	if d.InShootZone() {
		goForward := player.OpponentGoal().Center
		vecToGoal, err := physics.NewVector(player.Coords, goForward)
		if err != nil {
			d.Logger.Errorf("could not find the vector to the goal: %s", err)
			return "I have a issue, somebody else hold the ball", d.orderPassTheBall()
		}
		vecToGoal.SetY(0) //we want the player go straight

		obstaclesOnPath, err := strategy.WatchOpponentOnMyRoute(player.Coords, vecToGoal.TargetFrom(player.Coords), units.PlayerSize, player.GetOpponentTeam(d.GameMsg.GameInfo))
		if len(obstaclesOnPath) > 0 {
			return "I am free yet", []orders.Order{d.orderAdvance()}
		}
		speed, target := d.FindBestPointShootTheBall()
		return "Shoot!", []orders.Order{d.CreateKickOrder(target, speed)}
	} else {
		nextSteps := physics.NewVector(d.Player.Coords, d.OpponentGoal().Center).SetLength(units.PlayerMaxSpeed * 5)
		obstacles := watchOpponentOnMyRoute(d.LastMsg.GameInfo, d.Player, nextSteps.TargetFrom(d.Player.Coords))
		if len(obstacles) == 0 {
			if MyRule == strategy.DefensePlayer && (TeamState == strategy.Neutral || TeamState == strategy.Offensive) {
				return "Let's pass", d.orderPassTheBall()
			}
			return "I am free yet", []orders.Order{d.orderAdvance()}
		} else {
			return "I need help guys!", d.orderPassTheBall()
		}
	}
}

func DistanceForShooting(player *client.Player) float64 {
	goal := player.OpponentGoal()
	distX := math.Abs(float64(goal.Center.PosX - player.Coords.PosX))

	goalCoord := physics.Point{
		PosX: goal.Center.PosX,
		PosY: player.Coords.PosY,
	}
	if player.Coords.PosY < units.GoalMinY {
		goalCoord.PosY = goal.BottomPole.PosY
	} else if player.Coords.PosY > units.GoalMaxY {
		goalCoord.PosY = goal.TopPole.PosY
	}
	return player.Coords.DistanceTo(goalCoord)
}


fazer methodo que recebe jogador e diz se ele pode chutar
