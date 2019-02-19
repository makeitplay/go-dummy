package brain

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/BasicTypes"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"math"
)

// PlayerState defines states specific for players
type PlayerState BasicTypes.State

const (
	// Supporting identifies the player supporting the team mate
	Supporting PlayerState = "supporting"
	// HoldingTheBall identifies the player holding	the ball
	HoldingTheBall PlayerState = "holding"
	// Defending identifies the player defending against the opponent team
	Defending PlayerState = "defending"
	// DisputingTheBall identifies the player disputing the ball
	DisputingTheBall PlayerState = "disputing"
)

// PerfectPassDistance stores the constant distance where the ball reach in max speed after 1 frame
const PerfectPassDistance = float64(units.BallMaxSpeed - (units.BallDeceleration / 2))

// orderForDisputingTheBall returns a debug msg and a list of order for the DisputingTheBall state
func (b *Brain) orderForDisputingTheBall() (msg string, orders []BasicTypes.Order) {
	if b.ShouldIDisputeForTheBall() {
		msg = "Disputing for the ball"
		//orders = []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(b.LastMsg.GameInfo.Ball.Coords)}
		speed, target := b.FindBestPointInterceptBall()
		orders = []BasicTypes.Order{b.CreateMoveOrder(target, speed)}
		return msg, orders
	} else {
		if b.myCurrentRegion() != b.GetActiveRegion(TeamState) {
			return "Moving to my region", []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(b.GetActiveRegionCenter(TeamState))}
		} else {
			return "Holding position", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		}
	}
}

// orderForSupporting returns a debug msg and a list of order for the Supporting state
func (b *Brain) orderForSupporting() (msg string, orders []BasicTypes.Order) {
	if b.ShouldIAssist() { // middle players will give support
		return b.orderForActiveSupport()
	}
	return b.orderForPassiveSupport()
}

// orderForPassiveSupport returns a debug msg and a list of order for the Support state when the player is only holding position
func (b *Brain) orderForPassiveSupport() (msg string, orders []BasicTypes.Order) {
	var region strategy.RegionCode
	if b.ShouldIAssist() {
		region = FindSpotToAssist(
			b.LastMsg,
			b.LastMsg.GameInfo.Ball.Holder,
			b,
			false,
		)
	} else {
		region = b.GetActiveRegion(TeamState)
	}
	target := region.Center(b.TeamPlace)
	if b.Coords.DistanceTo(target) < units.PlayerMaxSpeed {
		if b.Velocity.Speed > 0 {
			orders = []BasicTypes.Order{b.CreateStopOrder(*physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
		}
	} else {
		orders = []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(target)}
	}
	return msg, orders
}

// orderForActiveSupport returns a debug msg and a list of order for the Support state when the player is assisting the ball holder
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
	if b.Coords.DistanceTo(target) < units.PlayerMaxSpeed {
		if b.Velocity.Speed > 0 {
			orders = []BasicTypes.Order{b.CreateStopOrder(*physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
		}
	} else {
		orders = []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(target)}
	}
	return "", orders
}

// orderForHoldingTheBall returns a debug msg and a list of order for the HoldingTheBall state
func (b *Brain) orderForHoldingTheBall() (msg string, orders []BasicTypes.Order) {
	goalCoords := b.OpponentGoal().Center
	goalDistance := b.Coords.DistanceTo(goalCoords)
	if goalDistance < strategy.RegionWidth*1.5 {
		nextSteps := physics.NewVector(b.Player.Coords, b.OpponentGoal().Center).SetLength(units.PlayerMaxSpeed * 2)
		obstacles := watchOpponentOnMyRoute(b.LastMsg.GameInfo, b.Player, nextSteps.TargetFrom(b.Player.Coords))
		if len(obstacles) == 0 && goalDistance < units.GoalZoneRange+units.PlayerMaxSpeed/2 {
			return "I am free yet", []BasicTypes.Order{b.orderAdvance()}
		}
		speed, target := b.FindBestPointShootTheBall()
		return "Shoot!", []BasicTypes.Order{b.CreateKickOrder(target, speed)}
	} else {
		nextSteps := physics.NewVector(b.Player.Coords, b.OpponentGoal().Center).SetLength(units.PlayerMaxSpeed * 5)
		obstacles := watchOpponentOnMyRoute(b.LastMsg.GameInfo, b.Player, nextSteps.TargetFrom(b.Player.Coords))
		if len(obstacles) == 0 {
			if MyRule == strategy.DefensePlayer && (TeamState == strategy.Neutral || TeamState == strategy.Offensive) {
				return "Let's pass", b.orderPassTheBall()
			}
			return "I am free yet", []BasicTypes.Order{b.orderAdvance()}
		} else {
			return "I need help guys!", b.orderPassTheBall()
		}
	}
}

// orderForDefending returns a debug msg and a list of order for the Defending state
func (b *Brain) orderForDefending() (msg string, orders []BasicTypes.Order) {
	if b.ShouldIDisputeForTheBall() {
		speed, target := b.FindBestPointInterceptBall()
		orders = []BasicTypes.Order{b.CreateMoveOrder(target, speed)}
	} else {
		target := b.GetActiveRegion(TeamState).Center(b.TeamPlace)
		if b.Coords.DistanceTo(target) < units.PlayerMaxSpeed {
			if b.Velocity.Speed > 0 {
				orders = []BasicTypes.Order{b.CreateStopOrder(*physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
			}
		} else {
			orders = []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(target)}
		}
	}
	//nothing more smart than that so far. stay stopped
	return msg, orders
}

// orderForGoalkeeper returns a debug msg and a list of order for the Goalkeeper state
func (b *Brain) orderForGoalkeeper() (msg string, orders []BasicTypes.Order) {
	//V = Vo + at -> t = Vo/a
	//framesToStop := units.BallMaxSpeed/units.BallDeceleration
	// (a*t^2)/2 + v*t - s
	//ballLongestShot := units.BallMaxSpeed*framesToStop + (-units.BallDeceleration/2) * math.Pow(framesToStop, 2)

	myGoal := b.DefenseGoal()
	longestDistance := units.GoalWidth - units.GoalKeeperJumpSpeed
	//s = so + vt
	t := float64(longestDistance/units.PlayerMaxSpeed) + 1 //11

	distanceWatchBall := units.BallMaxSpeed*t + float64(-units.BallDeceleration/2)*math.Pow(t, 2)

	if b.LastMsg.GameInfo.Ball.Coords.DistanceTo(myGoal.Center) <= distanceWatchBall {
		distanceToTopPole := b.LastMsg.GameInfo.Ball.Coords.DistanceTo(myGoal.TopPole)
		distanceToBottomPole := b.LastMsg.GameInfo.Ball.Coords.DistanceTo(myGoal.BottomPole)
		//find how many frames it would take from the closest place
		//(a*t^2)/2 + v*t - s
		t1, t2 := QuadraticResults(-units.BallDeceleration/2, units.BallMaxSpeed, -distanceToTopPole)
		framesToTop := int(math.Ceil(math.Min(t1, t2)))

		t1, t2 = QuadraticResults(-units.BallDeceleration/2, units.BallMaxSpeed, -distanceToBottomPole)
		framesToBottom := int(math.Ceil(math.Min(t1, t2)))

		var poleInRisk physics.Point
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
		maxDistanceICanRun := float64(units.PlayerMaxSpeed*frameToReact) + units.GoalKeeperJumpSpeed
		safePoint := physics.NewVector(poleInRisk, myGoal.Center).SetLength(maxDistanceICanRun).TargetFrom(poleInRisk)
		distanceToSafePoint := safePoint.DistanceTo(b.Coords)
		if distanceToSafePoint > units.PlayerMaxSpeed {
			return "Run to best spot!", []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(safePoint)}
		} else if distanceToSafePoint < 5 { //just a tolerance
			return "Be focused!!", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		} else {
			return "To center", []BasicTypes.Order{b.CreateMoveOrder(safePoint, distanceToSafePoint)}
		}

	} else {
		distanceFromMiddle := b.Coords.DistanceTo(myGoal.Center)
		if distanceFromMiddle > units.PlayerMaxSpeed {
			return "Back to position!", []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(myGoal.Center)}
		} else if distanceFromMiddle < 5 { //just a tolerance
			return "Just watch the game!", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		} else {
			return "To center", []BasicTypes.Order{b.CreateMoveOrder(myGoal.Center, distanceFromMiddle)}
		}
	}
}

//region helpers

//orderAdvance creates a move order towards the goal
func (b *Brain) orderAdvance() BasicTypes.Order {
	return b.CreateMoveOrderMaxSpeed(b.OpponentGoal().Center)
}

//orderPassTheBall estimates the best team mate for receiving a ball and creates a order to pass the ball to him
func (b *Brain) orderPassTheBall() []BasicTypes.Order {
	bestCandidate := new(client.Player)
	bestScore := 0
	for _, playerMate := range b.GetMyTeamStatus(b.LastMsg.GameInfo).Players {
		if playerMate.ID() == b.ID() {
			continue
		}
		goalCenter := b.OpponentGoal().Center
		obstaclesFromMe := watchOpponentOnMyRoute(b.LastMsg.GameInfo, b.Player, playerMate.Coords)
		obstaclesToGoal := watchOpponentOnMyRoute(b.LastMsg.GameInfo, playerMate, goalCenter)
		distanceFromMe := b.Coords.DistanceTo(playerMate.Coords)
		distanceToGoal := playerMate.Coords.DistanceTo(goalCenter)

		score := 100
		score -= len(obstaclesFromMe) * 10
		if len(obstaclesToGoal) == 0 && distanceToGoal < units.CourtWidth/4 {
			score += 40
		} else if len(obstaclesToGoal) > 0 {
			if obstaclesToGoal[0].DistanceTo(goalCenter) > 3.0*units.PlayerMaxSpeed {
				commons.LogDebug("obstaclesToGoal are further than 3 frames")
				score += 10
			} else if obstaclesToGoal[0].DistanceTo(goalCenter) > 1.0*units.PlayerMaxSpeed {
				commons.LogDebug("obstaclesToGoal are further than 1 frame")
				score += 5
			}
		}

		if distanceFromMe <= units.BallMaxSpeed/2 {
			score -= 10
		} else if math.Abs(distanceFromMe-PerfectPassDistance) < units.PlayerMaxSpeed {
			score += 20
		} else if distanceFromMe <= strategy.RegionWidth { // trocar pela largura da Ragion
			//commons.LogDebug("too far")
			score += 10
		} else {
			score += 10
		}
		if bestScore != 0 && distanceToGoal < bestCandidate.Coords.DistanceTo(goalCenter) {
			score += 10
		}
		if score > bestScore {
			bestScore = score
			bestCandidate = playerMate
		}
	}
	bastSpeed := b.BestSpeedToTarget(bestCandidate.Coords)

	return []BasicTypes.Order{
		b.CreateStopOrder(*physics.NewVector(b.LastMsg.GameInfo.Ball.Coords, bestCandidate.Coords).Normalize()),
		b.CreateKickOrder(bestCandidate.Coords, bastSpeed),
	}
}

//BestSpeedToTarget calculates the best speed to reach a specific point with the ball
func (b *Brain) BestSpeedToTarget(target physics.Point) float64 {
	distance := b.LastMsg.GameInfo.Ball.Coords.DistanceTo(target)
	ballSpeed := units.BallMaxSpeed
	deceleration := float64(units.BallDeceleration)

	//quadratic formula (-a/2)t^2 + vt - s
	A := -deceleration / 2
	B := ballSpeed
	C := -distance

	// delta: B^2 -4.A.C
	delta := math.Pow(B, 2) - 4*A*C

	// quadratic formula: -b +/- sqrt(delta)/2a
	t1 := (-B + math.Sqrt(delta)) / (2 * A)
	if math.IsNaN(t1) { // target too far
		return units.BallMaxSpeed
	}
	//t2 := (- B / math.Sqrt(delta)) / (2*A) //opposite side

	//S = So + Vt + (at^2)2
	//v =  ( s - (at^2)/2 ) / t
	s := distance
	ac := -deceleration
	t := math.Ceil(t1) // there is no half frame, so, 1.3 means more than one frame
	return (s - ((ac * math.Pow(t, 2)) / 2)) / t
}

//endregion
