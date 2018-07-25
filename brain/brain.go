package brain

import (
	"math"
	"fmt"

	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/GameState"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Physics"

	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/go-dummy/strategy"
)

// distance considered "near" for a player to the ball
const DistanceNearBall = strategy.RegionWidth // units float
const ERROR_MARGIN_RUNNING = 20.0
const ERROR_MARGIN_PASSING = 20.0

var TeamState = strategy.Defensive
var TeamBallPossession Units.TeamPlace
var MyRule strategy.PlayerRule

type Brain struct {
	*Game.Player
	State PlayerState
}

func (b *Brain) ResetPosition() {
	region := b.GetActiveRegion(strategy.Defensive)
	if MyRule == strategy.AttackPlayer {
		region = b.GetActiveRegion(strategy.UnderPressure)
	}
	b.Coords = region.Center(b.TeamPlace)
}

func (b *Brain) ProcessAnn(msg Game.GameMessage) {
	b.UpdatePosition(msg.GameInfo)
	commons.LogBroadcast(string(msg.State))
	switch GameState.State(msg.State) {
	case GameState.GETREADY:
	case GameState.LISTENING:
		if msg.GameInfo.Ball.Holder != nil {
			TeamBallPossession = msg.GameInfo.Ball.Holder.TeamPlace
		}
		TeamState = b.DetermineMyTeamState(msg)
		b.State = b.DetermineMyState()
		b.TakeAnAction()
	}
}

func (b *Brain) DetermineMyState() PlayerState {
	var isOnMyField bool
	var subState string
	var ballPossess string

	if b.LastMsg.GameInfo.Ball.Holder == nil {
		ballPossess = "dsp" //disputing
		subState = "fbl"    //far
		if int(math.Abs(b.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Coords))) <= DistanceNearBall {
			subState = "nbl" //near
		}
	} else if b.LastMsg.GameInfo.Ball.Holder.TeamPlace == b.TeamPlace {
		ballPossess = "atk" //attacking
		subState = "hlp"    //helping
		if b.LastMsg.GameInfo.Ball.Holder.ID() == b.ID() {
			subState = "hld" //holding
		}
	} else {
		ballPossess = "dfd" //defending
		subState = "org"
		if b.isItInMyActiveRegion(b.LastMsg.GameInfo.Ball.Coords, strategy.Defensive) {
			subState = "mrg"
		}
	}

	if b.TeamPlace == Units.HomeTeam {
		isOnMyField = b.LastMsg.GameInfo.Ball.Coords.PosX <= Units.CourtWidth/2
	} else {
		isOnMyField = b.LastMsg.GameInfo.Ball.Coords.PosX >= Units.CourtWidth/2
	}
	fieldState := "fr"
	if isOnMyField {
		fieldState = "hs"
	}
	return PlayerState(ballPossess + "-" + subState + "-" + fieldState)
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
	case DsptNfblHse:
		msg, orders = b.orderForDsptNfblHse()
		orders = append(orders, b.CreateCatchOrder())
	case DsptNfblFrg:
		msg, orders = b.orderForDsptNfblFrg()
		orders = append(orders, b.CreateCatchOrder())
	case DsptFrblHse:
		msg, orders = b.orderForDsptFrblHse()
		orders = append(orders, b.CreateCatchOrder())
	case DsptFrblFrg:
		msg, orders = b.orderForDsptFrblFrg()
		orders = append(orders, b.CreateCatchOrder())

	case AtckHoldHse:
		msg, orders = b.orderForAtckHoldHse()
	case AtckHoldFrg:
		msg, orders = b.orderForAtckHoldFrg()
	case AtckHelpHse:
		msg, orders = b.orderForAtckHelpHse()
	case AtckHelpFrg:
		msg, orders = b.orderForAtckHelpFrg()
	case DefdMyrgHse, DefdMyrgFrg, DefdOtrgHse, DefdOtrgFrg:
		msg, orders = b.orderForDefdOtrgFrg()
		orders = append(orders, b.CreateCatchOrder())
	default:
		msg = "Freeze position"
		orders = []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}

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
		if
		player.ID() != holderId && // the holder cannot help himself
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
	if goalkeeper.Coords.PosY > Units.CourtHeight / 2 {
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
	t := float64(longestDistance / Units.PlayerMaxSpeed) + 1 //11

	distanceWatchBall := Units.BallMaxSpeed*t + float64(-Units.BallDeceleration/2) * math.Pow(t, 2)

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
		maxDistanceICanRun := float64 (Units.PlayerMaxSpeed * frameToReact) + Units.GoalKeeperJumpLength
		safePoint := Physics.NewVector(poleInRisk, myGoal.Center).SetLength(maxDistanceICanRun).TargetFrom(poleInRisk)
		distanceToSafePoint := safePoint.DistanceTo(b.Coords)
		if distanceToSafePoint > Units.PlayerMaxSpeed {
			return "Run to best spot!", []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(safePoint)}
		} else if distanceToSafePoint < 5 {//just a tolerance
			return "Be focused!!", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		} else {
			return "To center", []BasicTypes.Order{b.CreateMoveOrder(safePoint, distanceToSafePoint)}
		}

	} else {
		distanceFromMiddle := b.Coords.DistanceTo(myGoal.Center)
		if distanceFromMiddle > Units.PlayerMaxSpeed {
			return "Back to position!", []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(myGoal.Center)}
		} else if distanceFromMiddle < 5 {//just a tolerance
			return "Just watch the game!", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		} else {
			return "To center", []BasicTypes.Order{b.CreateMoveOrder(myGoal.Center, distanceFromMiddle)}
		}
	}
}
