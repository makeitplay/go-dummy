package bot

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/geo"
	"github.com/lugobots/lugo4go/v2/proto"
	"math"
	"sort"
)


func (b Bot) OnHolding(ctx context.Context, data coach.TurnData) error {
	b.BallPossessionTeam = data.Me.TeamSide
	b.LastBallHolder = 0

	shootDecision := ShootingEvaluation(data.Me, data.Snapshot)
	playerCandidates := GetPassingCandidates(data.Me, data.Snapshot, b.LastBallHolder)
	passingDecision := PassingEvaluation(data.Me, playerCandidates, data.Snapshot)

	b.log.Debugf("Shoot: %d, passing: %d", shootDecision, passingDecision)
	// @explain:
	// 				Passing:	MustNot			shouldNot		may			Should			Must
	// Shooting:
	// MustNot					Advance			Advance			Advance		Pass			Pass
	// ShouldNot				Advance			Advance			Advance		Pass			Pass
	// May						Advance			Advance			CheckScore	Pass			Pass
	// Should					Shoot			Shoot			Shoot		Shoot			Shoot
	// Must						Shoot			Shoot			Shoot		Shoot			Shoot
	goalTarget := field.GetOpponentGoal(data.Me.TeamSide).Center
	if shootDecision >= Should {
		kickOrder, err := field.MakeOrderKick(*data.Snapshot.Ball, goalTarget, field.BallMaxSpeed)
		if err != nil {
			return fmt.Errorf("was not able to shoot: %s", err)
		}
		return send(ctx, data, []proto.PlayerOrder{kickOrder}, "shooting!")
	}
	// @speed-buster: we should check if there is a player, and avoid panics. However, "passingDecision" would not
	// be "positive" if there was not candidates.
	playerTarget := playerCandidates[0]
	if passingDecision >= Should {//|| (shootDecision == May && passingDecision == May && playerTarget.ShootingEvaluation > shootDecision) {
		kickOrder, err := field.MakeOrderKick(*data.Snapshot.Ball, *playerTarget.Player.Position, field.BallMaxSpeed)
		if err != nil {
			return fmt.Errorf("was not able to pass: %s", err)
		}
		return send(ctx, data, []proto.PlayerOrder{kickOrder}, fmt.Sprintf("passing! [%d] %v (vs %v)", passingDecision, playerTarget.Score, playerCandidates[1].Score))
	}

	//advance
	moveOrder, err := field.MakeOrderMoveMaxSpeed(*data.Me.Position, goalTarget)
	if err != nil {
		return fmt.Errorf("was not able to move forward: %s", err)
	}
	return send(ctx, data, []proto.PlayerOrder{moveOrder}, "advancing!")
}

type evaluator interface {
	ShootingEvaluation(me *proto.Player, snapshot *proto.GameSnapshot) FuzzyScale
	PassingEvaluation(me *proto.Player, snapshot *proto.GameSnapshot, candidates []PassingScore) FuzzyScale
	GetPassingCandidates(me *proto.Player, snapshot *proto.GameSnapshot, lastHolder uint32) []PassingScore
	IsObstacleForPassing(me *proto.Player, target proto.Point, opponents []*proto.Player) bool
	CountCloseOpponents(teamMate *proto.Player, opponents []*proto.Player) int
}


func ShootingEvaluation(me *proto.Player, snapshot *proto.GameSnapshot) FuzzyScale {
	opponentSide := field.GetOpponentSide(me.TeamSide)
	goalCenter := field.GetTeamsGoal(opponentSide).Center
	distanceToShoot := snapshot.Ball.Position.DistanceTo(goalCenter)
	if distanceToShoot >= DistanceFar {
		return MustNot
	}

	countObstacles := 0
	kickDirection, _ := proto.NewVector(*snapshot.Ball.Position, goalCenter)
	for _, p := range field.GetTeam(snapshot, opponentSide).Players {
		if p.Number != field.GoalkeeperNumber {
			angle := geo.AngleWithRoute(*kickDirection, goalCenter, *p.Position)
			if angle < 20 {
				countObstacles++
				if countObstacles > 1 {
					return ShouldNot
				}
			}
		}
	}
	// not that far, we could take the risk
	if distanceToShoot >= DistanceNear {
		return May
	} else if countObstacles == 1 { // close, but with obstacles
		return Should
	}
	// no obstacles, and pretty close!
	return Must
}

func PassingEvaluation(me *proto.Player, candidates []PassingScore, snapshot *proto.GameSnapshot) FuzzyScale {

	// @explain:
	// first, let's evaluate based on the risk of losing the ball possession
	// Is there a too close opponent?
	// Is there an opponent approaching me?

	// too close, pass NOW!
	closestDistance := float64(field.FieldWidth)
	goalTarget := field.GetOpponentGoal(me.TeamSide).Center
	goalDirection, _ := proto.NewVector(*snapshot.Ball.Position, goalTarget)
	for _, opponent := range field.GetTeam(snapshot, field.GetOpponentSide(me.TeamSide)).Players {
		if opponent.Number != field.GoalkeeperNumber {
			distance := opponent.Position.DistanceTo(*opponent.Position)

			obstacles := isObstacleForPassing(me, goalTarget, opponents []*proto.Player)

			if distance < field.PlayerSize*1.5 || angle < AngleOpponentObstacle {
				return Must
			}
			if distance < closestDistance {
				closestDistance = distance
			}
		}
	}
	// Getting close, better to pass!
	if closestDistance < field.PlayerSize*1.5 {
		return Should
	}

	// @explain:
	// Now we know we are safe to pass or not based on the team player positions
	// are there any player I can pass to?
	// Is there a player who may score?
	// Is there a player with no obstacle between we?
	if len(candidates) == 0 { //that's really hard to happen, but it is possible.
		return ShouldNot
	}

	//best player to get the pass
	bastCandidate := candidates[0]
	// a positive score means that the player is a great position and/or with a chance to score
	if bastCandidate.Score > 0 {
		return Must
	}

	if bastCandidate.Score > -ObstaclePenalty {
		return May
	}

	// Bad news: all players are in bad positions
	// Good news: we do not have to pass now.
	return MustNot
}

type PassingScore struct {
	Player             *proto.Player
	Score              int
	ShootingEvaluation FuzzyScale
}

const ObstaclePenalty = 12
const CloseOpponentsPenalty = 5
const LastHolderPenalty = 20
const DistancePenalty = 2
const GoalChanceBonus = 20
const LocationBonus = 2
const AngleOpponentObstacle = 10
const closeOpponentDistance = field.PlayerSize

func GetPassingCandidates(me *proto.Player, snapshot *proto.GameSnapshot, lastHolder uint32) []PassingScore {
	candidates := []PassingScore{}

	myDistanceToGoal := me.Position.DistanceTo(field.GetTeamsGoal(me.TeamSide).Center)
	opponents := field.GetTeam(snapshot, field.GetOpponentSide(me.TeamSide)).Players

	for _, p := range field.GetTeam(snapshot, me.TeamSide).Players {
		// don't pass to the goal keeper, neither to myself
		if p.Number != field.GoalkeeperNumber && p.Number != me.Number {

			distanceFromMe := p.Position.DistanceTo(*me.Position)
			if distanceFromMe >= DistanceFar {
				continue
			}

			punctuation := 0
			punctuation -= isObstacleForPassing(me, p, opponents) * ObstaclePenalty
			punctuation -= countCloseOpponents(p, opponents) * CloseOpponentsPenalty
			punctuation -= (int(distanceFromMe) / DistanceNear) * DistancePenalty
			if p.Number == lastHolder {
				punctuation -= LastHolderPenalty
			}
			goalEvaluation := ShootingEvaluation(p, snapshot)
			if goalEvaluation > May {
				punctuation += GoalChanceBonus
			}
			distanceToGoal := p.Position.DistanceTo(field.GetTeamsGoal(p.TeamSide).Center)
			if distanceToGoal > myDistanceToGoal {
				punctuation += (int(distanceToGoal) / DistanceNear) * LocationBonus
			}
			candidates = append(candidates, PassingScore{
				Player:             p,
				Score:              punctuation,
				ShootingEvaluation: goalEvaluation,
			})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return candidates[i].ShootingEvaluation > candidates[j].ShootingEvaluation
		}
		return candidates[i].Score > candidates[j].Score
	})
	return candidates
}

func isObstacleForPassing(me *proto.Player, target proto.Point, opponents []*proto.Player) int {
	obstacles := 0
	if passDirection, err := proto.NewVector(*me.Position, target); err == nil && passDirection != nil {
		distanceToTeamMate := me.Position.DistanceTo(target)
		for _, opponent := range opponents {
			distanceToOpponent := me.Position.DistanceTo(*opponent.Position)
			if distanceToTeamMate >= distanceToOpponent {
				angle := geo.AngleWithRoute(*passDirection, *me.Position, *opponent.Position)
				obstacle := math.Sin(angle * math.Pi / 180)
				if math.Abs(angle) < 90 && obstacle * distanceToOpponent < field.BallSize * 2 {
					obstacles++
				}
			}
		}
	}
	return obstacles
}

func countCloseOpponents(teamMate *proto.Player, opponents []*proto.Player) int {
	closeOpponents := 0
	for _, opponent := range opponents {
		if opponent.Position.DistanceTo(*teamMate.Position) < closeOpponentDistance {
			closeOpponents++
		}
	}
	return closeOpponents
}
