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

// TeamState stores the team state based on our strategy
var TeamState = strategy.Defensive

// TeamBallPossession stores the team's name that has touched on the ball for the last time
var TeamBallPossession Units.TeamPlace

// MyRule stores this player rule in the team
var MyRule strategy.PlayerRule

// Brain controls the player to have a behaviour during each state
type Brain struct {
	*client.Player
	State PlayerState
}

// ResetPosition set the player to it's initial position
func (b *Brain) ResetPosition() {
	region := b.GetActiveRegion(strategy.Defensive)
	if MyRule == strategy.AttackPlayer {
		region = b.GetActiveRegion(strategy.UnderPressure)
	}
	b.Coords = region.Center(b.TeamPlace)
}

// ProcessAnn is the callback function called when the player gets a new message from the game server
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

// DetermineMyState determine the player state bases on our strategy
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

// TakeAnAction sends orders to the game server based on the player state
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
	case Supporting:
		msg, orders = b.orderForSupporting()
	case HoldingTheBall:
		msg, orders = b.orderForHoldingTheBall()
	case Defending:
		msg, orders = b.orderForDefending()
		orders = append(orders, b.CreateCatchOrder())
	}
	b.SendOrders(fmt.Sprintf("[%s-%s] %s", b.State, TeamState, msg), orders...)
}

// ShouldIDisputeForTheBall returns true when the player should try to catch the ball
func (b *Brain) ShouldIDisputeForTheBall() bool {
	if strategy.GetRegionCode(b.LastMsg.GameInfo.Ball.Coords, b.TeamPlace).ChessDistanceTo(b.GetActiveRegion(TeamState)) < 2 {
		return true
	}
	myDistance := b.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Coords)
	playerCloser := 0
	for _, teamMate := range b.GetMyTeamStatus(b.LastMsg.GameInfo).Players {
		if teamMate.Number != b.Number && teamMate.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Coords) < myDistance {
			playerCloser++
			if playerCloser > 1 { // are there more than on player closer to the ball than me?
				return false
			}
		}
	}
	return true
}

// ShouldIAssist returns the ball when the player should support another team mate
func (b *Brain) ShouldIAssist() bool {
	holderRule := strategy.DefinePlayerRule(b.LastMsg.GameInfo.Ball.Holder.Number)
	if strategy.DefinePlayerRule(b.LastMsg.GameInfo.Ball.Holder.Number) == MyRule {
		return true
	}
	myDistance := b.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Holder.Coords)
	holderId := b.LastMsg.GameInfo.Ball.Holder.ID()
	playerCloser := 0
	for _, player := range b.GetMyTeamStatus(b.LastMsg.GameInfo).Players {
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

// FindBestPointInterceptBall finds a best spot around the ball holder to give support to him
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

// FindBestPointShootTheBall calculates the best point in the goal to shoot the ball
func (b *Brain) FindBestPointShootTheBall() (speed float64, target Physics.Point) {
	goalkeeper := b.FindOpponentPlayer(b.LastMsg.GameInfo, BasicTypes.PlayerNumber("1"))
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
