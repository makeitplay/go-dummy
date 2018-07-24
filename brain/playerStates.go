package brain

import (
	"math"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/go-dummy/strategy"
	"sort"
)

type PlayerState BasicTypes.State

const (
	// attacking, holding the ball, home field
	AtckHoldHse PlayerState = "atk-hld-hs"
	// attacking, holding the ball, foreign field
	AtckHoldFrg PlayerState = "atk-hld-fr"
	// attacking, helping the team, home field
	AtckHelpHse PlayerState = "atk-hlp-hs"
	// attacking, helping the team, foreign field
	AtckHelpFrg PlayerState = "atk-hlp-fr"

	// defading, on my region, home field
	DefdMyrgHse PlayerState = "dfd-mrg-hs"
	// defading, on my region, foreign field
	DefdMyrgFrg PlayerState = "dfd-mrg-fr"
	// defading, on other region, home field
	DefdOtrgHse PlayerState = "dfd-org-hs"
	// defading, on other region, foreign field
	DefdOtrgFrg PlayerState = "dfd-org-fr"

	// disputing, near to the ball, home field
	DsptNfblHse PlayerState = "dsp-nbl-hs"
	// disputing, near to the ball, foreign field
	DsptNfblFrg PlayerState = "dsp-nbl-fr"
	// disputing, far to the ball, home field
	DsptFrblHse PlayerState = "dsp-fbl-hs"
	// disputing, far to the ball, foreign field
	DsptFrblFrg PlayerState = "dsp-fbl-fr"
)

type DistanceScale string

const (
	DISTANCE_SCALE_NEAR DistanceScale = "near"
	DISTANCE_SCALE_FAR  DistanceScale = "far"
	DISTANCE_SCALE_GOOD DistanceScale = "good"
)

//region Disputing states
func (b *Brain) orderForDsptNfblHse() (msg string, orders []BasicTypes.Order) {
	if b.ShouldIDisputeForTheBall() {
		msg = "Disputing for the ball"
		orders = []BasicTypes.Order{b.CreateMoveOrder(b.LastMsg.GameInfo.Ball.Coords)}
		return msg, orders
	} else {
		if b.myCurrentRegion() != b.GetActiveRegion(TeamState) {
			return "Moving to my region", []BasicTypes.Order{b.CreateMoveOrder(b.GetActiveRegionCenter(TeamState))}
		} else {
			return "Holding position", []BasicTypes.Order{b.CreateStopOrder(*b.Velocity.Direction)}
		}
	}
}

func (b *Brain) orderForDsptNfblFrg() (msg string, orders []BasicTypes.Order) {
	return b.orderForDsptNfblHse()
}

func (b *Brain) orderForDsptFrblHse() (msg string, orders []BasicTypes.Order) {
	return b.orderForDsptNfblHse()
}

func (b *Brain) orderForDsptFrblFrg() (msg string, orders []BasicTypes.Order) {
	return b.orderForDsptNfblHse()
}
//endregion Disputing states

//region Attack states

func (b *Brain) orderForAtckHoldHse() (msg string, orders []BasicTypes.Order) {
	nextSteps := Physics.NewVector(b.Player.Coords, b.OpponentGoal().Center).SetLength(Units.PlayerMaxSpeed * 5)
	obstacles := watchOpponentOnMyRoute(b.LastMsg.GameInfo, b.Player, nextSteps.TargetFrom(b.Player.Coords))

	if len(obstacles) == 0 {
		return "I am free yet", []BasicTypes.Order{b.orderAdvance()}
	} else {
		return "I need help guys!", b.orderPassTheBall()
	}
}

func (b *Brain) orderForAtckHoldFrg() (msg string, orders []BasicTypes.Order) {
	goalCoords := b.OpponentGoal().Center
	goalDistance := b.Coords.DistanceTo(goalCoords)
	if math.Abs(goalDistance) < BallMaxSafePassDistance(Units.BallMaxSpeed) {
		return "Shoot!", []BasicTypes.Order{b.CreateKickOrder(goalCoords, Units.BallMaxSpeed)}
	} else {
		nextSteps := Physics.NewVector(b.Player.Coords, b.OpponentGoal().Center).SetLength(Units.PlayerMaxSpeed * 5)
		obstacles := watchOpponentOnMyRoute(b.LastMsg.GameInfo, b.Player, nextSteps.TargetFrom(b.Player.Coords))

		if len(obstacles) == 0 {
			return "I am free yet", []BasicTypes.Order{b.orderAdvance()}
		} else {
			return "I need help guys!", b.orderPassTheBall()
		}
	}
}

func (b *Brain) orderForAtckHelpHse() (msg string, orders []BasicTypes.Order) {
	var region strategy.RegionCode
	if b.ShouldIAssist() {
		region = FindSpotToAssist(
			b.LastMsg,
			b.LastMsg.GameInfo.Ball.Holder,
			b,
			false,
		)
	}else {
		region = b.GetActiveRegion(TeamState)
	}
	target := region.Center(b.TeamPlace)
	if b.Coords.DistanceTo(target) < Units.PlayerMaxSpeed {
		if b.Velocity.Speed > 0 {
			orders = []BasicTypes.Order{b.CreateStopOrder(*Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
		}
	} else {
		orders = []BasicTypes.Order{b.CreateMoveOrder(target)}
	}
	return msg, orders
}

func (b *Brain) orderForAtckHelpFrg() (msg string, orders []BasicTypes.Order) {
	if MyRule == strategy.DefensePlayer || !b.ShouldIAssist() { // middle players will give support
		return b.orderForAtckHelpHse()
	} else {
		//var bestCandidatePoint Physics.Point
		bestCandidateRegion := FindSpotToAssist(
			b.LastMsg,
			b.LastMsg.GameInfo.Ball.Holder,
			b,
			true,
		)
		//target := bestCandidateRegion.Center(b.TeamPlace)
		target := FindBestPointInRegionToAssist(
			b.LastMsg,
			bestCandidateRegion,
			b.LastMsg.GameInfo.Ball.Holder,
			)
		//obstacles := watchOpponentOnMyRoute(b.LastMsg.GameInfo.Ball.Holder, bestCandidatePoint)
		if b.Coords.DistanceTo(target) < Units.PlayerMaxSpeed {
			if b.Velocity.Speed > 0 {
				orders = []BasicTypes.Order{b.CreateStopOrder(*Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))}
			}
		} else {
			orders = []BasicTypes.Order{b.CreateMoveOrder(target)}
		}
	}
	return msg, orders
}
func FindBestPointInRegionToAssist(gameMessage Game.GameMessage, region strategy.RegionCode, assisted *Game.Player, ) (target Physics.Point) {
	centerPoint := region.Center(assisted.TeamPlace)
	vctToCenter := Physics.NewVector(assisted.Coords, centerPoint).SetLength(strategy.RegionWidth)
	obstacles := watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, vctToCenter.TargetFrom(assisted.Coords))
	if len(obstacles) == 0 {
		return vctToCenter.TargetFrom(assisted.Coords)
	} else {
		initialVector := vctToCenter
		avoidObstacles := func(ang float64) bool  {
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

		if !avoidObstacles(10) && !avoidObstacles(-10){
			target = initialVector.TargetFrom(assisted.Coords)
		}
	}

	return

}

func FindSpotToAssist(gameMessage Game.GameMessage, assisted *Game.Player, assistant *Brain, offensively bool) strategy.RegionCode {
	var availableSpots []strategy.RegionCode
	var spotList []strategy.RegionCode
	if offensively {
		spotList = ListSpotsCandidatesToOffensiveAssistance(assisted, assistant)
	} else {
		spotList = ListSpotsCandidatesToDefensiveAssistance(assisted, assistant)
	}
	for _, region := range spotList{
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
		APoints += float64(availableSpots[a].X - availableSpots[b].X) * 2.5
		return APoints >= 0
	})

	if len(availableSpots) > 0 {
		return availableSpots[0]
	}
	return assistant.GetActiveRegion(TeamState)
}
func ListSpotsCandidatesToOffensiveAssistance(assisted *Game.Player, assistant *Brain) []strategy.RegionCode {
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
func ListSpotsCandidatesToDefensiveAssistance(assisted *Game.Player, assistant *Brain) []strategy.RegionCode {
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

func isPerfectPlace(coords Physics.Point, gameMessage Game.GameMessage, assisted *Game.Player, assistant *Brain) bool {
	obstacles := watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, coords)
	bestPlaceRegion := strategy.GetRegionCode(coords, assistant.TeamPlace)

	thereIsOpponents := len(obstacles)
	thereIsNoMate := len(assistant.GetPlayersInRegion(bestPlaceRegion, assistant.FindMyTeamStatus(gameMessage.GameInfo))) == 0
	return thereIsOpponents == 0 && thereIsNoMate
}


//	if b.isItInMyActiveRegion(b.Coords) {
//		switch b.calcDistanceScale(b.LastMsg.GameInfo.Ball.Coords) {
//		case DISTANCE_SCALE_FAR:
//			msg = "Supporting on attack"
//			orders = []BasicTypes.Order{b.CreateMoveOrder(b.LastMsg.GameInfo.Ball.Coords)}
//		case DISTANCE_SCALE_NEAR:
//			msg = "Helping on attack"
//
//			offensiveZone := Physics.NewVector(b.Coords, b.GetActiveRegionCenter())
//			offensiveZone.Add(Physics.NewVector(b.Coords, b.OpponentGoal().Center))
//			orders = []BasicTypes.Order{b.CreateMoveOrder(offensiveZone.TargetFrom(b.Coords))}
//		case DISTANCE_SCALE_GOOD:
//			msg = "Holding positiong for attack"
//			offensiveZone := Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords)
//			offensiveZone.Add(Physics.NewVector(b.Coords, b.OpponentGoal().Center))
//			orders = []BasicTypes.Order{b.CreateMoveOrder(offensiveZone.TargetFrom(b.Coords))}
//		}
//	} else {
//		regionCenter := b.GetActiveRegionCenter()
//		return "Backing to my position", []BasicTypes.Order{b.CreateMoveOrder(regionCenter)}
//	}
//	return msg, orders
//}

//endregion Attack states

//region Defending states

//func (b *Brain) orderForDefdMyrgHse() (msg string, orders []BasicTypes.Order) {
//	orders = []BasicTypes.Order{b.CreateMoveOrder(b.LastMsg.GameInfo.Ball.Coords)}
//	return "Running towards the ball", orders
//}
//
//func (b *Brain) orderForDefdMyrgFrg() (msg string, orders []BasicTypes.Order) {
//	switch b.calcDistanceScale(b.LastMsg.GameInfo.Ball.Coords) {
//	case DISTANCE_SCALE_NEAR:
//		// too close
//		msg = "Pressing the player"
//		orders = []BasicTypes.Order{b.CreateMoveOrder(b.LastMsg.GameInfo.Ball.Coords)}
//	case DISTANCE_SCALE_FAR:
//		//get closer
//		msg = "Back to my position!"
//		var backOffPos Physics.Point
//		region := b.GetActiveRegion()
//		backOffPos = region.CentralDefense()
//		orders = []BasicTypes.Order{b.CreateMoveOrder(backOffPos)}
//	case DISTANCE_SCALE_GOOD:
//		msg = "Holding positiong"
//	}
//	//nothing more smart than that so far. stay stopped
//	return msg, orders
//}
//
//func (b *Brain) orderForDefdOtrgHse() (msg string, orders []BasicTypes.Order) {
//
//	if b.calcDistanceScale(b.LastMsg.GameInfo.Ball.Coords) == DISTANCE_SCALE_NEAR {
//		msg = "Defensing while back off"
//		backOffDir := Physics.NewVector(b.Coords, b.DefenseGoal().Center)
//		backOffDir.Add(Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))
//		orders = []BasicTypes.Order{b.CreateMoveOrder(backOffDir.TargetFrom(b.Coords))}
//	} else {
//		msg = "Back off!"
//		backOffDir := Physics.NewVector(b.Coords, b.DefenseGoal().Center)
//		backOffDir.Add(Physics.NewVector(b.Coords, b.GetActiveRegionCenter()))
//		orders = []BasicTypes.Order{b.CreateMoveOrder(backOffDir.TargetFrom(b.Coords))}
//	}
//	//nothing more smart than that so far. stay stopped
//	return msg, orders
//}
//
//func (b *Brain) orderForDefdOtrgFrg() (msg string, orders []BasicTypes.Order) {
//	if b.calcDistanceScale(b.LastMsg.GameInfo.Ball.Coords) == DISTANCE_SCALE_NEAR {
//		msg = "Defensing while back off"
//		backOffDir := Physics.NewVector(b.Coords, b.DefenseGoal().Center)
//		backOffDir.Add(Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))
//		orders = []BasicTypes.Order{b.CreateMoveOrder(backOffDir.TargetFrom(b.Coords))}
//	} else {
//		msg = "Back off!"
//		backOffDir := Physics.NewVector(b.Coords, b.DefenseGoal().Center)
//		backOffDir.Add(Physics.NewVector(b.Coords, b.GetActiveRegionCenter()))
//		orders = []BasicTypes.Order{b.CreateMoveOrder(backOffDir.TargetFrom(b.Coords))}
//	}
//	//nothing more smart than that so far. stay stopped
//	return msg, orders
//}

//endregion Defending states


//region helpers
func (b *Brain) orderAdvance() BasicTypes.Order {
	return b.CreateMoveOrder(b.OpponentGoal().Center)
}

func (b *Brain) orderPassTheBall() []BasicTypes.Order {
	bestCandidate := new(Game.Player)
	bestScore := 0
	for _, playerMate := range b.FindMyTeamStatus(b.LastMsg.GameInfo).Players {
		if playerMate.ID() == b.ID() {
			continue
		}
		//commons.LogWarning("Evaluating %s", playerMate.Number)


		obstaclesFromMe := watchOpponentOnMyRoute(b.LastMsg.GameInfo, b.Player, playerMate.Coords)
		obstaclesToGoal := watchOpponentOnMyRoute(b.LastMsg.GameInfo, playerMate, b.OpponentGoal().Center)
		distanceFromMe := b.Coords.DistanceTo(playerMate.Coords)
		distanceToGoal := playerMate.Coords.DistanceTo(b.OpponentGoal().Center)

		//commons.LogDebug("distanceFromMe %f", distanceFromMe)
		//commons.LogDebug("distanceToGoal %f", distanceToGoal)

		score := 0
		score -= len(obstaclesFromMe) * 10
		//commons.LogDebug("obstaclesFromMe %d", len(obstaclesFromMe) )
		if len(obstaclesToGoal) == 0 {
			//commons.LogDebug("obstaclesToGoal %d", len(obstaclesToGoal))
			score += 40
		//} else if obstaclesToGoal[0] > 3.0 * Units.PlayerMaxSpeed {
		//	commons.LogDebug("obstaclesToGoal are further than 3 frames")
			//score += 30
		//} else if obstaclesToGoal[0] > 1.0 * Units.PlayerMaxSpeed {
		//	commons.LogDebug("obstaclesToGoal are further than 1 frame")
			//score += 10
		}

		if distanceFromMe <= Units.BallMaxSpeed / 2 {
			//commons.LogDebug("too close")
			score -= 5
		} else if distanceFromMe > Units.BallMaxSpeed {
			//commons.LogDebug("too far")
			score -= 10
		}else {
			//commons.LogDebug("great distance")
			score += 20
		}
		if distanceToGoal < Units.BallMaxSpeed {
			//commons.LogDebug("awesome location")
			score += 20
		}
		//App.Log("=Player %s | %d obs, %d obs2, %d DfomMe, %d DfomGoal  = Total %d",
		//	playerMate.Number,
		//	len(obstaclesFromMe),
		//	len(obstaclesToGoal),
		//	int(distanceFromMe),
		//	int(distanceToGoal),
		//	score,
		//	)
		//commons.LogWarning("score candidate %v\n\n------------", score)
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
	s := b.LastMsg.GameInfo.Ball.Coords.DistanceTo(target)
	v := Units.BallMaxSpeed
	a := float64(-Units.BallDeceleration)
	//(a/2)t^2 + vt - s

	// delta: b^2 -4.a.c
	delta := math.Pow(v, 2) - 4 * -s * a

	// quadratic formula: -b +/- sqrt(delta)/2a
	t1 := (- v + math.Sqrt(delta)) / a
	if math.IsNaN(t1) {// target too far
		return Units.BallMaxSpeed
	}

	//t2 := (- v - math.Sqrt(delta)) / a

	t := math.Ceil(t1)

	return (s - (a/2)*math.Pow(t,2)) / t
}

// calc a distance scale where the player could target
func (b *Brain) calcDistanceScale(target Physics.Point) DistanceScale {
	distance := math.Abs(b.Coords.DistanceTo(target))
	// try to be closer the player
	toFar := Units.PlayerMaxSpeed * 4
	toNear :=Units.PlayerMaxSpeed * 2

	if distance >= toFar {
		return DISTANCE_SCALE_FAR
	} else if distance < toNear {
		return DISTANCE_SCALE_NEAR
	} else {
		return DISTANCE_SCALE_GOOD
	}
}


//endregion
