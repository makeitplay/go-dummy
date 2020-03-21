package dummy

import (
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/coach"
	"github.com/makeitplay/the-dummies-go/strategy"
	"k8s.io/apimachinery/pkg/util/rand"
	"math"
	"sort"
	"time"
)

var GameConfig *client.Configuration

var WaitingAnswer bool
var TunnelMsg chan client.GameMessage

var LastSuggestion string
var Passing bool
var DS coach.DataSaver
// Shoot/Pass		MustNot			shouldNot		may			Should			Must
// MustNot			Advance			Advance			Advance		Pass			Pass
// ShouldNot		Advance			Advance			Advance		Pass			Pass
// May				*				*				*			Pass			Pass
// Should			Shoot			Shoot			Shoot		Shoot			Shoot
// Must				Shoot			Shoot			Shoot		Shoot			Shoot
func (d *Dummy) orderForHoldingTheBall() (msg string, ordersSet []orders.Order) {
	player := d.Player

	ballPoint := d.GameMsg.Ball().Coords
	ballDistanceGoal := ballPoint.DistanceTo(player.OpponentGoal().Center)
	if false && rand.Int()%100 < 5 && ballDistanceGoal <= DistanceDistant {
		question := client.TrainingQuestion{
			Question:   "What should I do now?",
			QuestionId: fmt.Sprintf("%s-%s", d.Player.Id, time.Now()),
			PlayerId:   d.Player.ID(),
			Alternatives: []string{
				"pass",
				"shoot",
				"run",
				"dribble",
				"ignore",
			},
		}
		if err := client.AskQuestion(question, *GameConfig); err == nil {
			d.Logger.Warn("question sent")
			TunnelMsg = make(chan client.GameMessage)
			WaitingAnswer = true
			ds, err := DS.SaveSample(d.GameMsg.GameInfo)
			if err != nil {
				d.Logger.Errorf("did not create the state: %s", err )
			}
			d.Logger.Warnf("Bora esperar! ")
			var answer string
			for WaitingAnswer {
				select {
				case debugMsg := <-TunnelMsg:
					var ok bool
					if answer, ok = debugMsg.Data[question.QuestionId].(string); ok {
						d.Logger.Warnf("Got the answer! %s", answer)
						d.Logger.Warnf("Recebeu")
						WaitingAnswer = false
					} else {
						d.Logger.Warnf("Not yet")
					}
				}
			}
			straightForwards := physics.Point{
				PosX: player.OpponentGoal().Center.PosX,
				PosY: player.Coords.PosY,
			}
			if math.Abs(float64(player.Coords.PosY-player.OpponentGoal().Center.PosY)) < float64(DistanceFar) {
				straightForwards = player.OpponentGoal().Center
			}
			orderToAdvance, _ := d.Player.CreateMoveOrderMaxSpeed(straightForwards)

			switch answer {
			case "pass":
				_, candidatePlayers := fuzzyDecisionPass(player, d.GameMsg)
				if len(candidatePlayers) > 0 {
					bastCandidate := electBestCandidate(candidatePlayers, d.GameMsg)
					order, _ := d.Player.CreateKickOrder(d.GameMsg.Ball(), bastCandidate.Coords, units.BallMaxSpeed)
					ds.Save(answer)
					return "Ok mestre! Passando", []orders.Order{order}
				}
				return "Sorry mestre, canot pass", []orders.Order{orderToAdvance}
			case "shoot":
				_, target := ShouldShoot(player, d.GameMsg)
				if target == nil {
					tg := player.OpponentGoal().Center
					target = &tg
				}
				order, _ := d.Player.CreateKickOrder(d.GameMsg.Ball(), *target, units.BallMaxSpeed)
				Passing = true
				ds.Save(answer)
				return "OK mestre shhoting hoot!", []orders.Order{order}
			case "run", "dribble":
				ds.Save(answer)
				return "Go go go, o mestre mandou", []orders.Order{orderToAdvance}
			}

		}
	}

	shouldIShoot, target := ShouldShoot(player, d.GameMsg)
	minCase := Should
	goalPoint := d.Player.OpponentGoal().Center
	if goalPoint.DistanceTo(d.GameMsg.Ball().Coords) < DistanceDistant {
		minCase = May
	}

	if shouldIShoot >= minCase {
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
	if shouldIPass >= Should && len(candidatePlayers) > 0  {
		bastCandidate := electBestCandidate(candidatePlayers, d.GameMsg)
		order, _ := d.Player.CreateKickOrder(d.GameMsg.Ball(), bastCandidate.Coords, units.BallMaxSpeed)

		turnTo, _ := d.Player.CreateMoveOrderMaxSpeed(bastCandidate.Coords)
		Passing = true
		return "Found a well positioned mate", []orders.Order{turnTo, order}
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
		//targetToPlayerMate, err := physics.NewVector(player.Coords, playerMate.Coords)
		//if err != nil {
		//	return
		//}
		distanceFromMe := player.Coords.DistanceTo(playerMate.Coords)
		obstaclesToPlayer, err := strategy.WatchOpponentOnMyRoute(gameMsg.Ball().Coords, playerMate.Coords, units.BallSize/2, player.GetOpponentTeam(gameMsg.GameInfo))
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

		frontOfHim := (GameConfig.TeamPlace == arena.HomeTeam && opponent.Coords.PosX > player.Coords.PosX) ||
			(GameConfig.TeamPlace == arena.AwayTeam && opponent.Coords.PosX < player.Coords.PosX)

		if opponent.Coords.DistanceTo(player.Coords) < units.PlayerSize * 2 {
			nearOpponents++
			if frontOfHim {
				nearOpponents++
			}
		}
	})

	total := 100

	total -= (int(distanceFromMe) / units.PlayerSize) * 2
	total -= (int(distanceFromGoal) / units.PlayerSize) * 4
	// change the calc to only penalise when the opponent it between the player and the goal.
	total -= nearOpponents * 3

	if LastHolderFrom != nil && LastHolderFrom.Id == player.Id {
		//we probably received the ball from this guy, so let try do not send it to him
		total = int(0.8 * float64(total))
	}

	return total

}
