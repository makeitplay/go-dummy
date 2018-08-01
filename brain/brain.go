package brain

import (
	"fmt"
	"math"

	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/GameState"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"

	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
)

var TeamState = strategy.Defensive
var TeamBallPossession Units.TeamPlace
var MyRule strategy.PlayerRule

type Brain struct {
	*client.Player
	State PlayerState
}

func (b *Brain) ResetPosition() {
	region := b.GetActiveRegion(strategy.Defensive)
	if MyRule == strategy.AttackPlayer {
		region = b.GetActiveRegion(strategy.UnderPressure)
	}
	b.Coords = region.Center(b.TeamPlace)
}

func (b *Brain) ProcessAnn(msg client.GameMessage) {
	b.UpdatePosition(msg.GameInfo)
	commons.LogBroadcast(string(msg.State))
	switch GameState.State(msg.State) {
	case GameState.GetReady:
	case GameState.Listening:
		if msg.GameInfo.Ball.Holder != nil {
			TeamBallPossession = msg.GameInfo.Ball.Holder.TeamPlace
		}
		TeamState = b.DetermineMyTeamState(msg)
		b.State = b.DetermineMyState()
		b.TakeAnAction()
	}
}

func (b *Brain) DetermineMyState() PlayerState {
	if b.LastMsg.GameInfo.Ball.Holder == nil {
		return DisputingTheBall
	} else if b.LastMsg.GameInfo.Ball.Holder.TeamPlace == b.TeamPlace {
		if b.LastMsg.GameInfo.Ball.Holder.ID() == b.ID() {
			return HoldingTheBall
		}
		return Supporting
	}
	return Defending
}

func (b *Brain) TakeAnAction() {
	var orders []BasicTypes.Order
	var msg string

	if b.IsGoalkeeper() {
		msg, orders = b.orderForGoalkeeper()
		b.SendOrders(fmt.Sprintf("[%s-%s] %s", b.State, TeamState, msg), orders...)
		return
	}
	switch b.State {
	case DisputingTheBall:
		msg, orders = b.orderForDisputingTheBall()
		orders = append(orders, b.CreateCatchOrder())
	case HoldingTheBall:
		msg, orders = b.orderForAtckHoldFrg()
	case Defending:
		msg, orders = b.orderForDefending()
		orders = append(orders, b.CreateCatchOrder())
	case Supporting:
		msg, orders = b.orderForSupporting()
	}
	b.SendOrders(fmt.Sprintf("[%s-%s] %s", b.State, TeamState, msg), orders...)
}

func (b *Brain) ShouldIDisputeForTheBall() bool {
	if strategy.GetRegionCode(b.LastMsg.GameInfo.Ball.Coords, b.TeamPlace).ChessDistanceTo(b.GetActiveRegion(TeamState)) < 2 {
		return true
	}
	myDistance := b.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Coords)
	playerCloser := 0
	for _, teamMate := range b.FindMyTeamStatus(b.LastMsg.GameInfo).Players {
		if teamMate.Number != b.Number && teamMate.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Coords) < myDistance {
			playerCloser++
			if playerCloser > 1 { // are there more than on player closer to the ball than me?
				return false
			}
		}
	}
	return true
}

func (b *Brain) ShouldIAssist() bool {
	holderRule := strategy.DefinePlayerRule(b.LastMsg.GameInfo.Ball.Holder.Number)
	if strategy.DefinePlayerRule(b.LastMsg.GameInfo.Ball.Holder.Number) == MyRule {
		return true
	}
	myDistance := b.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Holder.Coords)
	holderId := b.LastMsg.GameInfo.Ball.Holder.ID()
	playerCloser := 0
	for _, player := range b.FindMyTeamStatus(b.LastMsg.GameInfo).Players {
		if player.ID() != holderId && // the holder cannot help himself
			player.Number != b.Number && // I wont count to myself
			strategy.DefinePlayerRule(player.Number) != holderRule && // I wont count with the players rule mates because they should ALWAYS help
			player.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Coords) < myDistance {
			playerCloser++
			if playerCloser > 1 { // are there more than two player closer to the ball than me?
				return false
			}
		}
	}
	return true
}

func (b *Brain) FindBestPointInterceptBall() (speed float64, target Physics.Point) {
	if b.LastMsg.GameInfo.Ball.Velocity.Speed == 0 {
		return Units.PlayerMaxSpeed, b.LastMsg.GameInfo.Ball.Coords
	} else {
		calcBallPos := func(frame int) *Physics.Point {
			//S = So + VT + (aT^2)/2
			V := b.LastMsg.GameInfo.Ball.Velocity.Speed
			T := float64(frame)
			a := -Units.BallDeceleration
			distance := V*T + (a*math.Pow(T, 2))/2
			if distance <= 0 {
				return nil
			}
			ballTarget := b.LastMsg.GameInfo.Ball.Velocity.Direction.Copy().
				SetLength(distance).
				TargetFrom(b.LastMsg.GameInfo.Ball.Coords)
			return &ballTarget
		}
		frames := 1
		lastBallPosition := b.LastMsg.GameInfo.Ball.Coords
		for {
			ballLocation := calcBallPos(frames)
			if ballLocation == nil {
				break
			}
			minDistanceToTouch := ballLocation.DistanceTo(b.Coords) - ((Units.BallSize + Units.PlayerSize) / 2)

			if minDistanceToTouch <= float64(Units.PlayerMaxSpeed*frames) {
				if frames > 1 {
					return Units.PlayerMaxSpeed, *ballLocation
				} else {
					return b.Coords.DistanceTo(*ballLocation), *ballLocation
				}
			}
			lastBallPosition = *ballLocation
			frames++
		}
		return Units.PlayerMaxSpeed, lastBallPosition
	}
}

func (b *Brain) FindBestPointShootTheBall() (speed float64, target Physics.Point) {
	goalkeeper := b.GetOpponentPlayer(b.LastMsg.GameInfo, BasicTypes.PlayerNumber("1"))
	if goalkeeper.Coords.PosY > Units.CourtHeight/2 {
		return Units.BallMaxSpeed, Physics.Point{
			PosX: b.OpponentGoal().BottomPole.PosX,
			PosY: b.OpponentGoal().BottomPole.PosY + Units.BallSize,
		}
	} else {
		return Units.BallMaxSpeed, Physics.Point{
			PosX: b.OpponentGoal().TopPole.PosX,
			PosY: b.OpponentGoal().TopPole.PosY - Units.BallSize,
		}
	}
}

func (b *Brain) orderForGoalkeeper() (msg string, orders []BasicTypes.Order) {
	//V = Vo + at -> t = Vo/a
	//framesToStop := Units.BallMaxSpeed/Units.BallDeceleration
	// (a*t^2)/2 + v*t - s
	//ballLongestShot := Units.BallMaxSpeed*framesToStop + (-Units.BallDeceleration/2) * math.Pow(framesToStop, 2)

	myGoal := b.DefenseGoal()
	longestDistance := Units.GoalWidth - Units.GoalKeeperJumpLength
	//s = so + vt
	t := float64(longestDistance/Units.PlayerMaxSpeed) + 1 //11

	distanceWatchBall := Units.BallMaxSpeed*t + float64(-Units.BallDeceleration/2)*math.Pow(t, 2)

	if b.LastMsg.GameInfo.Ball.Coords.DistanceTo(myGoal.Center) <= distanceWatchBall {
		distanceToTopPole := b.LastMsg.GameInfo.Ball.Coords.DistanceTo(myGoal.TopPole)
		distanceToBottomPole := b.LastMsg.GameInfo.Ball.Coords.DistanceTo(myGoal.BottomPole)
		//find how many frames it would take from the closest place
		//(a*t^2)/2 + v*t - s
		t1, t2 := QuadraticResults(-Units.BallDeceleration/2, Units.BallMaxSpeed, -distanceToTopPole)
		framesToTop := int(math.Ceil(math.Min(t1, t2)))

		t1, t2 = QuadraticResults(-Units.BallDeceleration/2, Units.BallMaxSpeed, -distanceToBottomPole)
		framesToBottom := int(math.Ceil(math.Min(t1, t2)))

		var poleInRisk Physics.Point
		var frameToReact int
		if framesToTop < framesToBottom {
			poleInRisk = myGoal.TopPole
			frameToReact = framesToTop
		} else {
			poleInRisk = myGoal.BottomPole
			frameToReact = framesToBottom
		}
		//the furthest safe place from the most risk side
		//S = so + vt
		maxDistanceICanRun := float64(Units.PlayerMaxSpeed*frameToReact) + Units.GoalKeeperJumpLength
		safePoint := Physics.NewVector(poleInRisk, myGoal.Center).SetLength(maxDistanceICanRun).TargetFrom(poleInRisk)
		distanceToSafePoint := safePoint.DistanceTo(b.Coords)
		if distanceToSafePoint > Units.PlayerMaxSpeed {
			return "Run to best spot!", []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(safePoint)}
		} else if distanceToSafePoint < 5 { //just a tolerance
			return "Be focused!!", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		} else {
			return "To center", []BasicTypes.Order{b.CreateMoveOrder(safePoint, distanceToSafePoint)}
		}

	} else {
		distanceFromMiddle := b.Coords.DistanceTo(myGoal.Center)
		if distanceFromMiddle > Units.PlayerMaxSpeed {
			return "Back to position!", []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(myGoal.Center)}
		} else if distanceFromMiddle < 5 { //just a tolerance
			return "Just watch the game!", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		} else {
			return "To center", []BasicTypes.Order{b.CreateMoveOrder(myGoal.Center, distanceFromMiddle)}
		}
	}
}

func (b *Brain) orderForActiveSupport() (msg string, orders []BasicTypes.Order) {
	bestCandidateRegion := FindSpotToAssist(
		b.LastMsg,
		b.LastMsg.GameInfo.Ball.Holder,
		b,
		true,
	)
	target := FindBestPointInRegionToAssist(
		b.LastMsg,
		bestCandidateRegion,
		b.LastMsg.GameInfo.Ball.Holder,
	)
	if b.Coords.DistanceTo(target) < Units.PlayerMaxSpeed {
		if b.Velocity.Speed > 0 {
			orders = []BasicTypes.Order{b.CreateStopOrder(*Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
		}
	} else {
		orders = []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(target)}
	}
	return "", orders
}
