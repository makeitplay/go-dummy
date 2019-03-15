package dummie

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"math"
)

// Shoot/Pass		MustNot			shouldNot		may			Should			Must
// MustNot			Advance			Advance			Advance		Pass			Pass
// ShouldNot		Advance			Advance			Advance		Pass			Pass
// May				*				*				*			Pass			Pass
// Should			Shoot			Shoot			Shoot		Shoot			Shoot
// Must				Shoot			Shoot			Shoot		Shoot			Shoot
func (d *Dummie) orderForHoldingTheBall() (msg string, ordersSet []orders.Order) {
	player := d.Player

	shouldIShoot, target := ShouldShoot(player, d.GameMsg)
	if shouldIShoot >= Should {
		order, _ := d.Player.CreateKickOrder(d.GameMsg.Ball(), *target, units.BallMaxSpeed)
		return "Shoot!", []orders.Order{order}
	}

	orderToAdvance, _ := d.Player.CreateMoveOrderMaxSpeed(physics.Point{
		PosX: player.OpponentGoal().Center.PosX,
		PosY: player.Coords.PosY,
	})

	shouldIPass, candidatePlayers := fuzzyDecisionPass(player, d.GameMsg)
	if shouldIPass >= Should {
		bastCandidate := electBestCandidate(candidatePlayers)
		order, _ := d.Player.CreateKickOrder(d.GameMsg.Ball(), bastCandidate.Coords, units.BallMaxSpeed)
		return "Found a well positioned mate", []orders.Order{order}
	}

	if shouldIShoot <= ShouldNot {
		return "No good options, keeping moving", []orders.Order{orderToAdvance}
	}

	//from this point we know we should not pass, and we may shoot but it is probably a bad idea
	return "do not stop to swimming", []orders.Order{orderToAdvance}

	return
}

func ShouldShoot(player *client.Player, gameMsg *client.GameMessage) (FuzzyScale, *physics.Point) {
	distanceToShoot := DistanceForShooting(player)
	if distanceToShoot >= DistanceFar {
		return MustNot, nil
	}
	betterTargetShoot := FindBestPointShootTheBall(player, gameMsg.GameInfo)
	shootVector, err := physics.NewVector(player.Coords, betterTargetShoot)
	targetAlternative := player.OpponentGoal().Center //we hope it never happens, but we need to return a value error if it does
	if err != nil {
		return MustNot, &targetAlternative
	}
	obstaclesToTarget, err := strategy.WatchOpponentOnMyRoute(player.Coords, shootVector.TargetFrom(player.Coords), units.BallSize, player.GetOpponentTeam(gameMsg.GameInfo))
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

func DistanceForShooting(player *client.Player) float64 {
	goal := player.OpponentGoal()
	ref := physics.Point{
		PosX: goal.Center.PosX,
		PosY: player.Coords.PosY,
	}
	if player.Coords.PosY < units.GoalMinY {
		ref = goal.BottomPole
	} else if player.Coords.PosY > units.GoalMaxY {
		ref = goal.TopPole
	}
	return player.Coords.DistanceTo(ref)
}

// FindBestPointShootTheBall find a good target in the opponent goal to shoot the ball at.
// @todo needs enhancement: the method is only choosing a side of the goal, but could consider the player position
func FindBestPointShootTheBall(player *client.Player, gameInfo client.GameInfo) (target physics.Point) {
	goalkeeper := player.FindOpponentPlayer(gameInfo, arena.GoalkeeperNumber)
	goal := player.OpponentGoal()
	if goalkeeper.Coords.PosY > goal.Center.PosY {
		return physics.Point{
			PosX: goal.BottomPole.PosX,
			PosY: goal.BottomPole.PosY + units.BallSize,
		}
	} else {
		return physics.Point{
			PosX: goal.TopPole.PosX,
			PosY: goal.TopPole.PosY - units.BallSize,
		}
	}
}

func electBestCandidate(players []*client.Player) *client.Player {
	return players[0]
}
