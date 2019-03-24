package dummy

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"math"
	"sort"
)

// Shoot/Pass		MustNot			shouldNot		may			Should			Must
// MustNot			Advance			Advance			Advance		Pass			Pass
// ShouldNot		Advance			Advance			Advance		Pass			Pass
// May				*				*				*			Pass			Pass
// Should			Shoot			Shoot			Shoot		Shoot			Shoot
// Must				Shoot			Shoot			Shoot		Shoot			Shoot
func (d *Dummy) orderForHoldingTheBall() (msg string, ordersSet []orders.Order) {
	player := d.Player

	shouldIShoot, target := ShouldShoot(player, d.GameMsg)
	if shouldIShoot >= Should {
		order, _ := d.Player.CreateKickOrder(d.GameMsg.Ball(), *target, units.BallMaxSpeed)
		return "Shoot!", []orders.Order{order}
	}

	straightForwards := physics.Point{
		PosX: player.OpponentGoal().Center.PosX,
		PosY: player.Coords.PosY,
	}
	if math.Abs(float64(player.Coords.PosY-player.OpponentGoal().Center.PosY)) < float64(DistanceFar) {
		straightForwards = player.OpponentGoal().Center
	}

	shouldIPass, candidatePlayers := fuzzyDecisionPass(player, d.GameMsg)
	orderToAdvance, err := d.Player.CreateMoveOrderMaxSpeed(straightForwards)
	if err != nil {
		return "something it wrong", []orders.Order{d.Player.CreateStopOrder(*player.Velocity.Direction)}
	}
	if shouldIPass >= Should {
		bastCandidate := electBestCandidate(candidatePlayers, d.GameMsg)
		order, _ := d.Player.CreateKickOrder(d.GameMsg.Ball(), bastCandidate.Coords, units.BallMaxSpeed)
		return "Found a well positioned mate", []orders.Order{order}
	}

	if shouldIShoot <= ShouldNot {
		return "No good options, keeping moving", []orders.Order{orderToAdvance}
	}

	//from this point we know we should not pass, and we may shoot but it is probably a bad idea
	return "just keep swimming", []orders.Order{orderToAdvance}
}

func ShouldShoot(player *client.Player, gameMsg *client.GameMessage) (FuzzyScale, *physics.Point) {
	distanceToShoot := DistanceForShooting(gameMsg.Ball(), player.OpponentGoal())
	if distanceToShoot >= DistanceFar {
		return MustNot, nil
	}
	betterTargetShoot := FindBestPointShootTheBall(player, gameMsg.GameInfo)
	shootVector, err := physics.NewVector(player.Coords, betterTargetShoot)
	targetAlternative := player.OpponentGoal().Center //we hope it never happens, but we need to return a value error if it does
	if err != nil {
		return MustNot, &targetAlternative
	}
	// @todo needs enhancement: if an opponent player stays in our way inside the goal zone, the player won't kick neither advance
	obstaclesToTarget, err := strategy.WatchOpponentOnMyRoute(gameMsg.Ball().Coords, shootVector.TargetFrom(player.Coords), units.BallSize, player.GetOpponentTeam(gameMsg.GameInfo))
	if err != nil {
		return MustNot, &targetAlternative
	}

	if len(obstaclesToTarget) == 0 {
		if distanceToShoot <= DistanceBeside {
			return Must, &betterTargetShoot
		}
		if distanceToShoot <= DistanceNear {
			return Should, &betterTargetShoot
		}
		if distanceToShoot <= DistanceFar {
			return May, &betterTargetShoot
		}
		return ShouldNot, &betterTargetShoot
	}
	return MustNot, &betterTargetShoot
}

func fuzzyDecisionPass(player *client.Player, gameMsg *client.GameMessage) (FuzzyScale, []*client.Player) {
	candidates := []*client.Player{}

	decisionPass := MustNot
	smallestDistance := float64(units.FieldWidth)
	gameMsg.ForEachPlayByTeam(player.TeamPlace, func(index int, playerMate *client.Player) {
		if playerMate.Id == player.Id {
			return
		}
		targetToPlayerMate, err := physics.NewVector(player.Coords, playerMate.Coords)
		if err != nil {
			return
		}
		distanceFromMe := player.Coords.DistanceTo(playerMate.Coords)
		obstaclesToPlayer, err := strategy.WatchOpponentOnMyRoute(player.Coords, targetToPlayerMate.TargetFrom(player.Coords), units.BallSize, player.GetOpponentTeam(gameMsg.GameInfo))
		if err != nil || len(obstaclesToPlayer) > 0 {
			return //no decision bases on this player mate
		}

		shouldShoot, _ := ShouldShoot(playerMate, gameMsg)
		if shouldShoot == Must || shouldShoot == Should {
			if distanceFromMe >= DistanceFar {
				return //he is at a nice position, but the ball
			}
			decisionPass = Should //we are able to pass to someone who may score, so we should pass
		}
		smallestDistance = math.Min(distanceFromMe, smallestDistance)
		candidates = append(candidates, playerMate)
	})

	if len(candidates) == 0 { //if I pass the ball I know I will lose it
		return MustNot, candidates
	}

	gameMsg.ForEachPlayByTeam(player.GetOpponentPlace(), func(index int, opponent *client.Player) {
		if opponent.Coords.DistanceTo(player.Coords) < units.PlayerSize*2 {
			decisionPass = Must
		}
	})

	if decisionPass < May { //we have candidates to receive the pass, but no one may score, so may pass if we decide to
		decisionPass = ShouldNot
		if smallestDistance <= DistanceNear {
			decisionPass = May
		}
	}

	return decisionPass, candidates
}

func DistanceForShooting(ball client.Ball, goal arena.Goal) float64 {
	ref := physics.Point{
		PosX: goal.Center.PosX,
		PosY: ball.Coords.PosY,
	}
	if ball.Coords.PosY < units.GoalMinY {
		ref = goal.BottomPole
	} else if ball.Coords.PosY > units.GoalMaxY {
		ref = goal.TopPole
	}
	return ball.Coords.DistanceTo(ref)
}

// FindBestPointShootTheBall find a good target in the opponent goal to shoot the ball at.
// @todo needs enhancement: the method is only choosing a side of the goal, but could consider the player position
func FindBestPointShootTheBall(player *client.Player, gameInfo client.GameInfo) (target physics.Point) {
	goalkeeper := player.FindOpponentPlayer(gameInfo, arena.GoalkeeperNumber)
	goal := player.OpponentGoal()
	if goalkeeper.Coords.PosY > goal.Center.PosY {
		return physics.Point{
			PosX: goal.BottomPole.PosX,
			PosY: goal.BottomPole.PosY + (units.BallSize / 2),
		}
	} else {
		return physics.Point{
			PosX: goal.TopPole.PosX,
			PosY: goal.TopPole.PosY - (units.BallSize / 2),
		}
	}
}

func electBestCandidate(players []*client.Player, gameMsg *client.GameMessage) *client.Player {
	sort.Slice(players, func(i, j int) bool {
		return passReceiverScore(players[i], gameMsg) > passReceiverScore(players[j], gameMsg)
	})
	return players[0]
}

func passReceiverScore(player *client.Player, gameMsg *client.GameMessage) int {
	ball := gameMsg.Ball()
	distanceFromMe := ball.Coords.DistanceTo(player.Coords)
	distanceFromGoal := DistanceForShooting(gameMsg.Ball(), player.OpponentGoal())
	nearOpponents := 0
	gameMsg.ForEachPlayByTeam(player.GetOpponentPlace(), func(index int, opponent *client.Player) {
		if opponent.Coords.DistanceTo(player.Coords) < DistanceBeside {
			nearOpponents++
		}
	})

	total := 100

	total -= int(distanceFromMe) / (units.PlayerSize * 2)
	total -= int(distanceFromGoal) / (units.PlayerSize * 4)
	total -= nearOpponents

	if LastHolderFrom != nil && LastHolderFrom.Id == player.Id {
		//we probably received the ball from this guy, so let try do not send it to him
		total = int(0.8 * float64(total))
	}

	return total

}
