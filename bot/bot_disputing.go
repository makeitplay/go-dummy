package bot

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"math"
)

func (b Bot) OnDisputing(ctx context.Context, turn coach.TurnData) (err error) {
	debugMsg := ""
	orders := make([]proto.PlayerOrder, 0, 2)
	var moveOrder proto.PlayerOrder

	//expectedRegion := GetMyRegion(OnAttack, b.Positioner, turn.Me.Number)
	//moveOrder, err = field.MakeOrderMoveMaxSpeed(*turn.Me.Position, expectedRegion.Center())
	//if err != nil {
	//	return fmt.Errorf("was not able to move this turn: %s", err)
	//}
	//orders = []proto.PlayerOrder{moveOrder}
	//return send(ctx, turn, orders, debugMsg)

	if ShouldIDisputeForTheBall(turn.Me, turn.Snapshot) {
		speed, target := FindBestPointInterceptBall(*turn.Snapshot.Ball, turn.Me)

		moveOrder, err = field.MakeOrderMove(*turn.Me.Position, target, speed)
		if err != nil {
			return fmt.Errorf("was not able to move this turn: %s", err)
		}
		orders = []proto.PlayerOrder{moveOrder, field.MakeOrderCatch()}
		debugMsg = "trying to catch the ball"
	} else {
		// @see Readme- ignored errors
		ballRegion, _ := b.Positioner.GetPointRegion(*turn.Snapshot.Ball.Position)

		// @see Readme- ignored errors
		teamState, _ := DetermineTeamState(ballRegion, turn.Me.TeamSide, b.BallPossessionTeam)

		// @see Readme- ignored errors
		currentReg, _ := b.Positioner.GetPointRegion(*turn.Me.Position)

		expectedRegion := GetMyRegion(teamState, b.Positioner, turn.Me.Number)

		b.log.Debugf("C %s, Ex %s", currentReg, expectedRegion)

		if currentReg.String() != expectedRegion.String() {
			moveOrder, err = field.MakeOrderMoveMaxSpeed(*turn.Me.Position, expectedRegion.Center())
			if err != nil {
				return fmt.Errorf("was not able to move this turn: %s", err)
			}
			debugMsg = fmt.Sprintf("moving to my region: %v", expectedRegion)
		} else {
			moveOrder, err = field.MakeOrderMove(*turn.Me.Position, *turn.Snapshot.Ball.Position, 0)
			if err != nil {
				return fmt.Errorf("was not able to move this turn: %s", err)
			}
			debugMsg = fmt.Sprintf("holding position (%s %s)", expectedRegion, currentReg)
		}
		orders = []proto.PlayerOrder{moveOrder}
	}

	return send(ctx, turn, orders, debugMsg)
}

func ShouldIDisputeForTheBall(me *proto.Player, snapshot *proto.GameSnapshot) bool {
	// suggestion: The bot should considering if the ball is coming to his position, or going away from him
	distanceToBall := snapshot.Ball.Position.DistanceTo(*me.Position)

	playerCloser := 0
	for _, p := range field.GetTeam(snapshot, me.TeamSide).Players {
		ddd := p.Position.DistanceTo(*snapshot.Ball.Position)
		if p.Number != me.Number && distanceToBall > ddd {
			playerCloser++
			if playerCloser > 1 {
				return false
			}
		}
	}
	return true

	//
	//player := d.Player
	//if d.ShouldIDisputeForTheBall() {
	//msg = "Disputing for the ball"
	////orders = []orders.Order{d.CreateMoveOrderMaxSpeed(d.LastMsg.GameInfo.Ball.Coords)}
	//speed, target := strategy.FindBestPointInterceptBall(d.GameMsg.Ball(), player)
	//movOrder, err := player.CreateMoveOrder(target, speed)
	//if err != nil {
	//d.Logger.Errorf("error creating move order: %s ", err)
	//msg = "sorry, I won't play this turn"
	//} else {
	//ordersSet = []orders.Order{movOrder}
	//}
	//} else {
	//if d.myCurrentRegion() != d.GetActiveRegion() {
	//movOrder, err := player.CreateMoveOrderMaxSpeed(d.GetActiveRegionCenter())
	//if err != nil {
	//d.Logger.Errorf("error creating move order: %s ", err)
	//msg = "sorry, I won't play this turn"
	//} else {
	//msg = "Moving to my region"
	//ordersSet = []orders.Order{movOrder}
	//}
	//} else {
	//msg = "Holding position"
	//ordersSet = []orders.Order{player.CreateStopOrder(*player.Velocity.Direction)}
	//}
}

func FindBestPointInterceptBall(ball proto.Ball, player *proto.Player) (speed float64, target proto.Point) {
	if ball.Velocity.Speed == 0 {
		return field.PlayerMaxSpeed, *ball.Position
	} else {
		calcBallPos := func(frame int) *proto.Point {
			//S = So + VT + (aT^2)/2
			V := ball.Velocity.Speed
			T := float64(frame)
			a := -field.BallDeceleration
			distance := V*T + (a*math.Pow(T, 2))/2
			if distance <= 0 {
				return nil
			}
			vectorToBal, _ := ball.Velocity.Direction.Copy().SetLength(distance)
			ballTarget := vectorToBal.TargetFrom(*ball.Position)
			return &ballTarget
		}
		frames := 1
		lastBallPosition := ball.Position
		for {
			ballLocation := calcBallPos(frames)
			if ballLocation == nil {
				break
			}
			minDistanceToTouch := ballLocation.DistanceTo(*player.Position) - ((field.BallSize + field.PlayerSize) / 2)

			if minDistanceToTouch <= float64(field.PlayerMaxSpeed*frames) {
				if frames > 1 {
					return field.PlayerMaxSpeed, *ballLocation
				} else {
					return player.Position.DistanceTo(*ballLocation), *ballLocation
				}
			}
			lastBallPosition = ballLocation
			frames++
		}
		return field.PlayerMaxSpeed, *lastBallPosition
	}

}
