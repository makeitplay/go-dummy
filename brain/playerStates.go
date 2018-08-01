package brain

import (
	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"
	"math"
	"sort"

	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
)

type PlayerState BasicTypes.State

const (
	Supporting PlayerState = "supporting"

	HoldingTheBall PlayerState = "holding"

	Defending PlayerState = "defending"

	// disputing the ball
	DisputingTheBall PlayerState = "disputing"
)

const PerfectPassDistance = float64(Units.BallMaxSpeed - (Units.BallDeceleration / 2))

//region Disputing states
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

//endregion Disputing states

//region Attack states

func (b *Brain) orderForAtckHoldFrg() (msg string, orders []BasicTypes.Order) {
	goalCoords := b.OpponentGoal().Center
	goalDistance := b.Coords.DistanceTo(goalCoords)
	if goalDistance < strategy.RegionWidth*1.5 {
		nextSteps := Physics.NewVector(b.Player.Coords, b.OpponentGoal().Center).SetLength(Units.PlayerMaxSpeed * 2)
		obstacles := watchOpponentOnMyRoute(b.LastMsg.GameInfo, b.Player, nextSteps.TargetFrom(b.Player.Coords))
		if len(obstacles) == 0 && goalDistance > Units.GoalZoneRange {
			return "I am free yet", []BasicTypes.Order{b.orderAdvance()}
		}
		speed, target := b.FindBestPointShootTheBall()
		return "Shoot!", []BasicTypes.Order{b.CreateKickOrder(target, speed)}
	} else {
		nextSteps := Physics.NewVector(b.Player.Coords, b.OpponentGoal().Center).SetLength(Units.PlayerMaxSpeed * 5)
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
	if b.Coords.DistanceTo(target) < Units.PlayerMaxSpeed {
		if b.Velocity.Speed > 0 {
			orders = []BasicTypes.Order{b.CreateStopOrder(*Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
		}
	} else {
		orders = []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(target)}
	}
	return msg, orders
}

func (b *Brain) orderForSupporting() (msg string, orders []BasicTypes.Order) {
	//if MyRule == strategy.DefensePlayer || !b.ShouldIAssist() { // middle players will give support
	//	return b.orderForPassiveSupport()
	//} else {
	//	return b.orderForActiveSupport()
	//}
	if b.ShouldIAssist() { // middle players will give support
		return b.orderForActiveSupport()
	}
	return b.orderForPassiveSupport()
}
func FindBestPointInRegionToAssist(gameMessage client.GameMessage, region strategy.RegionCode, assisted *client.Player) (target Physics.Point) {
	centerPoint := region.Center(assisted.TeamPlace)
	vctToCenter := Physics.NewVector(assisted.Coords, centerPoint).SetLength(strategy.RegionWidth)
	obstacles := watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, vctToCenter.TargetFrom(assisted.Coords))
	if len(obstacles) == 0 {
		return vctToCenter.TargetFrom(assisted.Coords)
	} else {
		initialVector := vctToCenter
		avoidObstacles := func(ang float64) bool {
			tries := 3
			for tries > 0 {
				vctToCenter.AddAngleDegree(ang)
				target = vctToCenter.TargetFrom(assisted.Coords)
				if region != strategy.GetRegionCode(target, assisted.TeamPlace) {
					//too far
					tries = 0
				}
				obstacles = watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, target)
				tries--
				if len(obstacles) <= 0 {
					return true
				}
			}
			return false
		}

		if !avoidObstacles(10) && !avoidObstacles(-10) {
			target = initialVector.TargetFrom(assisted.Coords)
		}
	}

	return

}

func FindSpotToAssist(gameMessage client.GameMessage, assisted *client.Player, assistant *Brain, offensively bool) strategy.RegionCode {
	var availableSpots []strategy.RegionCode
	var spotList []strategy.RegionCode
	if offensively {
		spotList = ListSpotsCandidatesToOffensiveAssistance(assisted, assistant)
	} else {
		spotList = ListSpotsCandidatesToDefensiveAssistance(assisted, assistant)
	}
	for _, region := range spotList {
		mateInTheRegion := assistant.GetPlayersInRegion(region, assistant.FindMyTeamStatus(gameMessage.GameInfo))
		if len(mateInTheRegion) == 0 {
			availableSpots = append(availableSpots, region)
		} else if region == assistant.GetActiveRegion(TeamState) {
			// eu to no meu canto, me deixe em paz
			availableSpots = append(availableSpots, region)
		} else {
			frankenstein := Brain{Player: mateInTheRegion[0]}
			isHimTheOwner := region == frankenstein.GetActiveRegion(TeamState)
			if !isHimTheOwner && assistant.myCurrentRegion() == region {
				// two invasors disputing
				myDistanceToTheBall := assistant.Coords.DistanceTo(assisted.Coords)
				invasorDistanceToTheBall := assistant.Coords.DistanceTo(mateInTheRegion[0].Coords)
				if myDistanceToTheBall < invasorDistanceToTheBall {
					availableSpots = append(availableSpots, region)
				}
			}
		}
	}
	sort.Slice(availableSpots, func(a, b int) bool {
		teamStatus := assistant.GetOpponentTeam(gameMessage.GameInfo)
		opponentsInA := len(assistant.GetPlayersInRegion(availableSpots[a], teamStatus))
		opponentsInB := len(assistant.GetPlayersInRegion(availableSpots[b], teamStatus))

		distanceToA := math.Round(assistant.Coords.DistanceTo(availableSpots[a].Center(assistant.TeamPlace)) / strategy.RegionWidth)
		distanceToB := math.Round(assistant.Coords.DistanceTo(availableSpots[b].Center(assistant.TeamPlace)) / strategy.RegionWidth)

		distanceAToAssistant := math.Round(assisted.Coords.DistanceTo(availableSpots[a].Center(assistant.TeamPlace)) / strategy.RegionWidth)
		distanceBToAssistant := math.Round(assisted.Coords.DistanceTo(availableSpots[b].Center(assistant.TeamPlace)) / strategy.RegionWidth)

		APoints := distanceToB - distanceToA
		APoints += float64(opponentsInB - opponentsInA)
		APoints += distanceBToAssistant - distanceAToAssistant
		APoints += float64(availableSpots[a].X-availableSpots[b].X) * 2.5
		return APoints >= 0
	})

	if len(availableSpots) > 0 {
		return availableSpots[0]
	}
	return assistant.GetActiveRegion(TeamState)
}
func ListSpotsCandidatesToOffensiveAssistance(assisted *client.Player, assistant *Brain) []strategy.RegionCode {
	spotCollection := []strategy.RegionCode{}
	currentRegion := strategy.GetRegionCode(assisted.Coords, assistant.TeamPlace)

	front := currentRegion.Forwards()
	if front != currentRegion {
		spotCollection = append(spotCollection, front)
	}

	assistantActiveRegion := assistant.GetActiveRegion(TeamState)

	goodRegionA := front.Left()
	if currentRegion != front && goodRegionA.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, goodRegionA)
	}
	goodRegionB := front.Right()
	if currentRegion != front && goodRegionB.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, goodRegionB)
	}

	fairRegionA := currentRegion.Left()
	if currentRegion != fairRegionA && fairRegionA.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, fairRegionA)
	}
	fairRegionB := currentRegion.Right()
	if currentRegion != fairRegionB && fairRegionB.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, fairRegionB)
	}
	return spotCollection
}
func ListSpotsCandidatesToDefensiveAssistance(assisted *client.Player, assistant *Brain) []strategy.RegionCode {
	spotCollection := []strategy.RegionCode{}
	currentRegion := strategy.GetRegionCode(assisted.Coords, assistant.TeamPlace)

	back := currentRegion.Backwards()
	if back != currentRegion {
		spotCollection = append(spotCollection, back)
	}

	goodRegionA := back.Left()
	if currentRegion != back {
		spotCollection = append(spotCollection, goodRegionA)
	}
	goodRegionB := back.Right()
	if currentRegion != back {
		spotCollection = append(spotCollection, goodRegionB)
	}

	return spotCollection
}

func (b *Brain) orderForDefending() (msg string, orders []BasicTypes.Order) {
	if b.ShouldIDisputeForTheBall() {
		speed, target := b.FindBestPointInterceptBall()
		orders = []BasicTypes.Order{b.CreateMoveOrder(target, speed)}
	} else {
		target := b.GetActiveRegion(TeamState).Center(b.TeamPlace)
		if b.Coords.DistanceTo(target) < Units.PlayerMaxSpeed {
			if b.Velocity.Speed > 0 {
				orders = []BasicTypes.Order{b.CreateStopOrder(*Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
			}
		} else {
			orders = []BasicTypes.Order{b.CreateMoveOrderMaxSpeed(target)}
		}
	}
	//nothing more smart than that so far. stay stopped
	return msg, orders
}

//endregion Defending states

//region helpers
func (b *Brain) orderAdvance() BasicTypes.Order {
	return b.CreateMoveOrderMaxSpeed(b.OpponentGoal().Center)
}

func (b *Brain) orderPassTheBall() []BasicTypes.Order {
	bestCandidate := new(client.Player)
	bestScore := 0
	for _, playerMate := range b.FindMyTeamStatus(b.LastMsg.GameInfo).Players {
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
		if len(obstaclesToGoal) == 0 && distanceToGoal < Units.CourtWidth/4 {
			score += 40
		} else if len(obstaclesToGoal) > 0 {
			if obstaclesToGoal[0].DistanceTo(goalCenter) > 3.0*Units.PlayerMaxSpeed {
				commons.LogDebug("obstaclesToGoal are further than 3 frames")
				score += 10
			} else if obstaclesToGoal[0].DistanceTo(goalCenter) > 1.0*Units.PlayerMaxSpeed {
				commons.LogDebug("obstaclesToGoal are further than 1 frame")
				score += 5
			}
		}

		if distanceFromMe <= Units.BallMaxSpeed/2 {
			//commons.LogDebug("too close")
			score -= 10
		} else if math.Abs(distanceFromMe-PerfectPassDistance) < Units.PlayerMaxSpeed {
			score += 20
		} else if distanceFromMe <= strategy.RegionWidth { // trocar pela largura da Ragion
			//commons.LogDebug("too far")
			score += 10
		} else {
			//commons.LogDebug("great distance")
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
	//commons.LogWarning("Best candidate %s", bestCandidate.Number)
	bastSpeed := b.BestSpeedToTarget(bestCandidate.Coords)

	return []BasicTypes.Order{
		b.CreateStopOrder(*Physics.NewVector(b.LastMsg.GameInfo.Ball.Coords, bestCandidate.Coords).Normalize()),
		b.CreateKickOrder(bestCandidate.Coords, bastSpeed),
	}
}

func (b *Brain) BestSpeedToTarget(target Physics.Point) float64 {
	distance := b.LastMsg.GameInfo.Ball.Coords.DistanceTo(target)
	ballSpeed := Units.BallMaxSpeed
	deceleration := float64(Units.BallDeceleration)

	//quadratic formula (-a/2)t^2 + vt - s
	A := -deceleration / 2
	B := ballSpeed
	C := -distance

	// delta: B^2 -4.A.C
	delta := math.Pow(B, 2) - 4*A*C

	// quadratic formula: -b +/- sqrt(delta)/2a
	t1 := (-B + math.Sqrt(delta)) / (2 * A)
	if math.IsNaN(t1) { // target too far
		return Units.BallMaxSpeed
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
