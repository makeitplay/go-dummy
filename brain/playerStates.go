package brain

import (
	"math"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/client-player-go/Game"
	"reflect"
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
	obstacles := watchOpponentOnMyRoute(b.Player, b.OpponentGoal().Center)

	if len(obstacles) == 0 {
		return "I am free yet", []BasicTypes.Order{b.orderAdvance()}
	} else if len(obstacles) == 1 {
		num := reflect.ValueOf(obstacles).MapKeys()[0].String()
		if b.calcDistanceScale(b.GetOpponentPlayer(b.LastMsg.GameInfo, num).Coords) != DISTANCE_SCALE_FAR {
			return "Dribble this guys (not yet)", b.orderPassTheBall()
		} else {
			return "Advance watching", []BasicTypes.Order{b.orderAdvance()}
		}
	} else {
		nearstObstacle := float64(Units.CourtWidth) //just initializing with a high value
		//num := int(reflect.ValueOf(obstacles).MapKeys()[0].Int())
		for opponentId := range obstacles {
			oppCoord := b.GetOpponentPlayer(b.LastMsg.GameInfo, opponentId).Coords
			oppDist := b.Coords.DistanceTo(oppCoord)
			if oppDist < nearstObstacle {
				nearstObstacle = oppDist
			}
		}
		if nearstObstacle < Units.PlayerMaxSpeed*2 {
			return "I need help guys!", b.orderPassTheBall()
		} else {
			return "Advance watching", []BasicTypes.Order{b.orderAdvance()}
		}
	}
}

func (b *Brain) orderForAtckHoldFrg() (msg string, orders []BasicTypes.Order) {
	goalCoords := b.OpponentGoal().Center
	goalDistance := b.Coords.DistanceTo(goalCoords)
	if math.Abs(goalDistance) < BallMaxSafePassDistance(Units.BallMaxSpeed) {
		return "Shoot!", []BasicTypes.Order{b.CreateKickOrder(goalCoords, Units.BallMaxSpeed)}
	} else {
		obstacles := watchOpponentOnMyRoute(b.Player, b.OpponentGoal().Center)

		if len(obstacles) == 0 {
			return "I am still free", []BasicTypes.Order{b.orderAdvance()}
		} else if len(obstacles) == 1 {
			num := reflect.ValueOf(obstacles).MapKeys()[0].String()
			if b.calcDistanceScale(b.GetOpponentPlayer(b.LastMsg.GameInfo, num).Coords) != DISTANCE_SCALE_FAR {
				return "Dribble this guys (not yet)", b.orderPassTheBall()
			} else {
				return "Advace watching", []BasicTypes.Order{b.orderAdvance()}
			}
		} else {
			return "I need help", b.orderPassTheBall()
		}

	}
}

//func (b *Brain) orderForAtckHelpHse() (msg string, orders []BasicTypes.Order) {
//	if b.isItInMyActiveRegion(b.Coords) {
//		switch b.calcDistanceScale(b.LastMsg.GameInfo.Ball.Coords) {
//		case DISTANCE_SCALE_FAR:
//			msg = "Let's attack!"
//			orders = []BasicTypes.Order{b.CreateMoveOrder(b.LastMsg.GameInfo.Ball.Coords)}
//		case DISTANCE_SCALE_NEAR:
//			msg = "Given space"
//			opositPoint := Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords).Invert().TargetFrom(b.Coords)
//			vectorToOpositPoint := Physics.NewVector(b.Coords, b.OpponentGoal().Center)
//			vectorToOpositPoint.Add(Physics.NewVector(b.Coords, opositPoint))
//			orders = []BasicTypes.Order{b.CreateMoveOrder(vectorToOpositPoint.TargetFrom(b.Coords))}
//		case DISTANCE_SCALE_GOOD:
//			msg = "Give me the ball!"
//			orders = []BasicTypes.Order{b.CreateMoveOrder(b.LastMsg.GameInfo.Ball.Coords)}
//		}
//	} else {
//		msg = "I'll be right here"
//		myRegionVector := Physics.NewVector(b.Coords, b.GetActiveRegionCenter()).Invert().TargetFrom(b.Coords)
//		offensivePosition := Physics.NewVector(b.Coords, b.OpponentGoal().Center)
//		offensivePosition.Add(Physics.NewVector(b.Coords, myRegionVector))
//		orders = []BasicTypes.Order{b.CreateMoveOrder(offensivePosition.TargetFrom(b.Coords))}
//	}
//	return msg, orders
//}
//
//func (b *Brain) orderForAtckHelpFrg() (msg string, orders []BasicTypes.Order) {
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


		obstaclesFromMe := watchOpponentOnMyRoute(b.Player, playerMate.Coords)
		obstaclesToGoal := watchOpponentOnMyRoute(playerMate, b.OpponentGoal().Center)
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
		} else if obstaclesToGoal[0] > 3.0 * Units.PlayerMaxSpeed {
			//commons.LogDebug("obstaclesToGoal are further than 3 frames")
			score += 30
		} else if obstaclesToGoal[0] > 1.0 * Units.PlayerMaxSpeed {
			//commons.LogDebug("obstaclesToGoal are further than 1 frame")
			score += 10
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
