package brain

import (
	"math"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/client-player-go/Game"
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

		return "Holding position for support", []BasicTypes.Order{b.CreateMoveOrder(b.myActiveRegionCenter(TeamState))}
	}
}
//
//func (b *Brain) orderForDsptNfblFrg() (msg string, orders []BasicTypes.Order) {
//	return b.orderForDsptNfblHse()
//}

func (b *Brain) orderForDsptFrblHse() (msg string, orders []BasicTypes.Order) {
	return b.orderForDsptNfblHse()
	//msg = "Try to catch the ball"
	//
	//if b.isItInMyActiveRegion(b.LastMsg.GameInfo.Ball.Coords) {
	//	backOffDir := Physics.NewVector(b.Coords, b.DefenseGoal().Center)
	//	backOffDir.Add(Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))
	//	orders = []BasicTypes.Order{b.CreateMoveOrder(backOffDir.TargetFrom(b.Coords))}
	//} else {
	//	orders = []BasicTypes.Order{b.CreateMoveOrder(b.myActiveRegionCenter())}
	//}
	//return msg, orders
}

//func (b *Brain) orderForDsptFrblFrg() (msg string, orders []BasicTypes.Order) {
//	msg = "Watch out the ball"
//	backOffDir := Physics.NewVector(b.Coords, b.myActiveRegionCenter())
//	backOffDir.Add(Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords))
//	orders = []BasicTypes.Order{b.CreateMoveOrder(backOffDir.TargetFrom(b.Coords))}
//	return msg, orders
//}
//endregion Disputing states

//region Attack states

//func (b *Brain) orderForAtckHoldHse() (msg string, orders []BasicTypes.Order) {
//	obstacles := watchOpponentOnMyRoute(b.Player, b.OpponentGoal().Center, ERROR_MARGIN_RUNNING)
//
//	if len(obstacles) == 0 {
//		return "I am free yet", []BasicTypes.Order{b.orderAdvance()}
//	} else if len(obstacles) == 1 {
//		num := int(reflect.ValueOf(obstacles).MapKeys()[0].Int())
//		if b.calcDistanceScale(b.GetOpponentTeam(b.LastMsg.GameInfo).Players[num].Coords) != DISTANCE_SCALE_FAR {
//			return "Dribble this guys (not yet)", []BasicTypes.Order{b.orderPassTheBall()}
//		} else {
//			return "Advance watching", []BasicTypes.Order{b.orderAdvance()}
//		}
//	} else {
//		nearstObstacle := float64(Units.CourtWidth) //just initializing with a high value
//		//num := int(reflect.ValueOf(obstacles).MapKeys()[0].Int())
//		for opponentId := range obstacles {
//			oppCoord := b.GetOpponentTeam(b.LastMsg.GameInfo).Players[opponentId].Coords
//			oppDist := b.Coords.DistanceTo(oppCoord)
//			if oppDist < nearstObstacle {
//				nearstObstacle = oppDist
//			}
//		}
//		if nearstObstacle < Units.PlayerMaxSpeed*2 {
//			return "I need help", []BasicTypes.Order{b.orderPassTheBall(), b.orderAdvance()}
//		} else {
//			return "Advance watching", []BasicTypes.Order{b.orderAdvance()}
//		}
//	}
//}
//
//func (b *Brain) orderForAtckHoldFrg() (msg string, orders []BasicTypes.Order) {
//	goalCoords := b.OpponentGoal().Center
//	goalDistance := b.Coords.DistanceTo(goalCoords)
//	if int(math.Abs(goalDistance)) < BallMaxDistance() {
//		return "Shoot!", []BasicTypes.Order{b.CreateKickOrder(goalCoords)}
//	} else {
//		obstacles := watchOpponentOnMyRoute(b.Player, b.OpponentGoal().Center, ERROR_MARGIN_RUNNING)
//
//		if len(obstacles) == 0 {
//			return "I am still free", []BasicTypes.Order{b.orderAdvance()}
//		} else if len(obstacles) == 1 {
//			num := int(reflect.ValueOf(obstacles).MapKeys()[0].Int())
//			if b.calcDistanceScale(b.GetOpponentTeam(b.LastMsg.GameInfo).Players[num].Coords) != DISTANCE_SCALE_FAR {
//				return "Dribble this guys (not yet)", []BasicTypes.Order{b.orderPassTheBall()}
//			} else {
//				return "Advace watching", []BasicTypes.Order{b.orderAdvance()}
//			}
//		} else {
//			return "I need help", []BasicTypes.Order{b.orderPassTheBall(), b.orderAdvance()}
//		}
//
//	}
//}
//
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
//		myRegionVector := Physics.NewVector(b.Coords, b.myActiveRegionCenter()).Invert().TargetFrom(b.Coords)
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
//			offensiveZone := Physics.NewVector(b.Coords, b.myActiveRegionCenter())
//			offensiveZone.Add(Physics.NewVector(b.Coords, b.OpponentGoal().Center))
//			orders = []BasicTypes.Order{b.CreateMoveOrder(offensiveZone.TargetFrom(b.Coords))}
//		case DISTANCE_SCALE_GOOD:
//			msg = "Holding positiong for attack"
//			offensiveZone := Physics.NewVector(b.Coords, b.LastMsg.GameInfo.Ball.Coords)
//			offensiveZone.Add(Physics.NewVector(b.Coords, b.OpponentGoal().Center))
//			orders = []BasicTypes.Order{b.CreateMoveOrder(offensiveZone.TargetFrom(b.Coords))}
//		}
//	} else {
//		regionCenter := b.myActiveRegionCenter()
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
//		backOffDir.Add(Physics.NewVector(b.Coords, b.myActiveRegionCenter()))
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
//		backOffDir.Add(Physics.NewVector(b.Coords, b.myActiveRegionCenter()))
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

func (b *Brain) orderPassTheBall() BasicTypes.Order {
	bestCandidate := new(Game.Player)
	bestScore := 0
	for _, playerMate := range b.GetMyTeam(b.LastMsg.GameInfo).Players {
		if playerMate.Id == b.Id {
			continue
		}
		obstaclesFromMe := watchOpponentOnMyRoute(b.Player, playerMate.Coords, ERROR_MARGIN_PASSING)
		obstaclesToGoal := watchOpponentOnMyRoute(playerMate, b.OpponentGoal().Center, ERROR_MARGIN_RUNNING)
		distanceFromMe := b.Coords.DistanceTo(playerMate.Coords)
		distanceToGoal := playerMate.Coords.DistanceTo(b.OpponentGoal().Center)
		score := 1000
		score -= len(obstaclesFromMe) * 10
		score -= len(obstaclesToGoal) * 5
		score -= int(distanceFromMe * 0.5)
		score -= int(distanceToGoal * 0.5)

		//App.Log("=Player %s | %d obs, %d obs2, %d DfomMe, %d DfomGoal  = Total %d",
		//	playerMate.Number,
		//	len(obstaclesFromMe),
		//	len(obstaclesToGoal),
		//	int(distanceFromMe),
		//	int(distanceToGoal),
		//	score,
		//	)

		if score > bestScore {
			bestScore = score
			bestCandidate = playerMate
		}
	}

	//App.Log("\n=Best candidate %d ", bestCandidate.Number)
	return b.CreateKickOrder(bestCandidate.Coords)
}

// calc a distance scale where the player could target
func (b *Brain) calcDistanceScale(target Physics.Point) DistanceScale {
	distance := math.Abs(b.Coords.DistanceTo(target))
	// try to be closer the player
	fieldDiagonal := math.Hypot(float64(Units.CourtHeight), float64(Units.CourtWidth))
	toFar := fieldDiagonal / 3
	toNear := fieldDiagonal / 5

	if distance >= toFar {
		return DISTANCE_SCALE_FAR
	} else if distance < toNear {
		return DISTANCE_SCALE_NEAR
	} else {
		return DISTANCE_SCALE_GOOD
	}
}

// Opponent id and angle between it and the target
func watchOpponentOnMyRoute(player *Game.Player, target Physics.Point, errMarginDegree float64) map[int]float64 {
	opponentTeam := player.GetOpponentTeam(player.LastMsg.GameInfo)
	opponents := make(map[int]float64)
	for _, opponent := range opponentTeam.Players {
		angle, isObstacle := player.IsObstacle(target, opponent.Coords, errMarginDegree)
		if isObstacle {
			opponents[opponent.Id] = angle
		}
	}
	return opponents
}
//endregion
