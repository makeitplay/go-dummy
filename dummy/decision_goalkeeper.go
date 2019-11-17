package dummy

import (
	"github.com/lugobots/arena"
	"github.com/lugobots/arena/orders"
	"github.com/lugobots/arena/physics"
	"github.com/lugobots/arena/units"
	"github.com/lugobots/client-player-go"
	"github.com/lugobots/the-dummies-go/strategy"
	"math"
)

func (d *Dummy) orderForGoalkeeper() (msg string, ordersSet []orders.Order) {
	keeper := d.Player
	ball := d.GameMsg.Ball()
	myGoal := keeper.DefenseGoal()
	stopOrder := d.Player.CreateStopOrder(*keeper.Velocity.Direction)
	if ball.Coords.DistanceTo(keeper.Coords) > DistanceFar {
		return "no risk of goal", stayAtGoalCenter(keeper)
	}
	//1st - Based on the ball's speed in X axis, let find how long the ball would take to reach the coord X of the goal
	//2nd - find the nearest point the bal may reach at the goal
	//3rd - calculate the further point the keeper to be from that point to cover more area
	ballTarget, timeToReach, coming := findThreatenedSpot(ball, myGoal)
	if !coming {
		optimumPoint := optimumWatchingPosition(myGoal, ballTarget, timeToReach)
		if optimumPoint == keeper.Coords {
			return "I am ready", []orders.Order{stopOrder}
		}
		moveOrder, _ := d.Player.CreateMoveOrderMaxSpeed(optimumPoint)
		return "save position", []orders.Order{moveOrder}
	}

	return "Catch", catchingTheBall(keeper, ballTarget, timeToReach)
}
func stayAtGoalCenter(player *client.Player) []orders.Order {
	stopOrder := player.CreateStopOrder(*player.Velocity.Direction)
	center := player.DefenseGoal().Center
	if player.Coords == center {
		return []orders.Order{stopOrder}
	}
	speed := math.Min(units.PlayerMaxSpeed, center.DistanceTo(player.Coords))
	order, _ := player.CreateMoveOrder(center, speed)
	return []orders.Order{order}
}

func findThreatenedSpot(ball client.Ball, goal arena.Goal) (target physics.Point, framesToReach int, coming bool) {
	ballSpeed := ball.Velocity.Speed
	ballXSpeed := ball.Velocity.Direction.Cos() * ballSpeed
	ballYSpeed := ball.Velocity.Direction.Sin() * ballSpeed

	if ball.Holder != nil {
		//let think what could happen if the ball was kicked now
		ballSpeed = units.BallMaxSpeed
		// if the ball wasn't kicked yet, the nearest point is the threatened
		target = NearestGoalPoint(ball, goal)

		ballkick, _ := physics.NewVector(ball.Coords, target)
		ballXSpeed = ballkick.Cos() * ballSpeed
		ballYSpeed = ballkick.Sin() * ballSpeed
	}

	//S = So + V.T + (a/2).T^2
	//S: Goal X coord
	//So: Actual ball X coord
	//V: ballXSpeed
	//T: required
	//a: deceleration of the ball
	S := goal.Center.PosX
	So := ball.Coords.PosX
	a := -units.BallDeceleration / 2
	// (a/2).T^2 +  V.T + (So - S)
	t1, t2, err := strategy.QuadraticResults(a, ballXSpeed, float64(So-S))
	if err != nil {
		return
	}
	realTimeToReach := t1 // truncating as integer because our time is calculated on frames
	if t1 <= 0 || (t2 > 0 && t2 < t1) {
		realTimeToReach = t2
	}

	if realTimeToReach < 0 {
		return
	}
	framesToReach = int(math.Ceil(realTimeToReach))

	// if the ball was kicked, let find the target based on its velocity
	//S: required
	//So: Actual ball Y coord
	//V: ballYSpeed
	//T:  "realTimeToReach"
	//a: deceleration of the ball
	coordY := float64(ball.Coords.PosY) + (ballYSpeed * realTimeToReach) + (a/2)*math.Pow(realTimeToReach, 2)

	target = physics.Point{
		PosX: goal.Center.PosX,
		PosY: int(math.Round(coordY)),
	}

	coming = target.PosY < units.GoalMaxY && target.PosY > units.GoalMinY
	if ball.Holder != nil || ball.Velocity.Speed <= 0 {
		coming = false
	}
	return
}
func catchingTheBall(keeper *client.Player, ballGoalTarget physics.Point, timeAvailable int) []orders.Order {

	distanceFromTarget := math.Abs(float64(ballGoalTarget.PosY - keeper.Coords.PosY))

	if distanceFromTarget < units.BallSize+units.PlayerSize {
		//the keeper can already catch the ball!
		return []orders.Order{
			keeper.CreateStopOrder(*keeper.Velocity.Direction),
			keeper.CreateCatchOrder(),
		}
	}

	timeToCatch := int(distanceFromTarget / units.PlayerMaxSpeed)
	if timeAvailable <= units.GoalKeeperJumpDuration && timeToCatch > timeAvailable {
		idealSpeed := distanceFromTarget / units.GoalKeeperJumpDuration //we need ensure the jump won't be beyond the target
		jump, _ := keeper.CreateJumpOrder(ballGoalTarget, idealSpeed)
		return []orders.Order{
			jump,
			keeper.CreateCatchOrder(),
		}
	}

	//the keep has time to catch the ball
	keeperSpeed := units.PlayerMaxSpeed
	if distanceFromTarget < units.PlayerMaxSpeed {
		keeperSpeed = distanceFromTarget //we do not want to pass through the ball target
	}
	moveOrder, _ := keeper.CreateMoveOrder(ballGoalTarget, keeperSpeed)

	return []orders.Order{
		moveOrder,
		keeper.CreateCatchOrder(),
	}
}

// @todo it cam be enhanced: this function is not considering the player size, so the keeper could be further from the target sometimes
func optimumWatchingPosition(goal arena.Goal, threatenedPoint physics.Point, timeAvailable int) physics.Point {
	jumpDistance := units.GoalKeeperJumpDuration * units.GoalKeeperJumpSpeed

	distanceFromCenter := goal.Center.DistanceTo(threatenedPoint)
	if jumpDistance > distanceFromCenter {
		distanceFromCenter -= jumpDistance
		timeAvailable -= units.GoalKeeperJumpDuration
	}

	if timeAvailable <= 1 { // too late!
		return threatenedPoint
	}

	//this is the time the keep would take to reach the threatenedPoint if he started from the goal center
	timeNeededToReach := int(distanceFromCenter / units.PlayerMaxSpeed)
	if timeNeededToReach > timeAvailable {
		//the keeper is `lateTIme` frames late to reach the ball if the ball was kicked now
		lateTIme := timeNeededToReach - timeAvailable
		gapDistance := lateTIme * units.PlayerMaxSpeed
		//the gap is only in the Y axis, os it is easy to find the best point:

		optimumPoint := goal.Center
		if threatenedPoint.PosY > goal.Center.PosY { //above the center
			optimumPoint.PosY += gapDistance
		} else {
			optimumPoint.PosY -= gapDistance
		}
		return optimumPoint
	}

	//it's fine stay in the center, we have time enough to reach the threatenedPoint
	return goal.Center
}

func NearestGoalPoint(ball client.Ball, goalTarget arena.Goal) physics.Point {
	nearest := physics.Point{
		PosX: goalTarget.Center.PosX,
		PosY: ball.Coords.PosY,
	}
	if ball.Coords.PosY < units.GoalMinY {
		nearest.PosY = goalTarget.BottomPole.PosY + (units.BallSize / 2)
	} else if ball.Coords.PosY > units.GoalMaxY {
		nearest.PosY = goalTarget.TopPole.PosY - (units.BallSize / 2)
	}

	return nearest
}
