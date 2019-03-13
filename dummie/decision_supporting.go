package dummie

import (
	"fmt"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
)

func (d *Dummie) orderForSupporting() (msg string, orders []orders.Order) {
	if d.ShouldIAssist() { // middle players will give support
		return d.orderForActiveSupport()
	}
	return //d.orderForActiveSupport()
}

func (d *Dummie) ShouldIAssist() bool {
	if strategy.DefinePlayerRule(d.GameMsg.GameInfo.Ball.Holder.Number) <= MyRule {
		return true
	}
	myDistance := d.Player.Coords.DistanceTo(d.GameMsg.GameInfo.Ball.Holder.Coords)
	holderId := d.GameMsg.GameInfo.Ball.Holder.ID()
	playerCloser := 0

	shouldAssist := true
	// shoul have at least 3 supporters in the perimeters around the ball holder
	d.GameMsg.ForEachPlayByTeam(TeamPlace, func(index int, player *client.Player) {
		if player.ID() != holderId && // the holder cannot help himself
			player.Number != PlayerNumber && // I wont count to myself
			player.Coords.DistanceTo(d.GameMsg.GameInfo.Ball.Coords) < myDistance {
			playerCloser++
			if playerCloser > 2 { // are there more than two player closer to the ball than me?
				shouldAssist = false
			}
		}
	})
	return shouldAssist
}

func (d *Dummie) orderForActiveSupport() (msg string, ordersSet []orders.Order) {
	ballPosition := d.GameMsg.GameInfo.Ball.Holder.Coords
	holderId := d.GameMsg.GameInfo.Ball.Holder.ID()
	referencies := []physics.Point{ballPosition}

	d.GameMsg.ForEachPlayByTeam(TeamPlace, func(index int, player *client.Player) {
		if player.ID() != holderId && // the holder cannot help himself
			player.Number != PlayerNumber && // I wont count to myself
			//player.Coords.DistanceTo(d.GameMsg.GameInfo.Ball.Coords) <= DistanceNear && //the guys who are already supporting
			player.Coords.DistanceTo(d.Player.Coords) < DistanceBeside {
			referencies = append(referencies, player.Coords)
		}
	})

	target := FindSpotToAssist(referencies, d.Player)

	if target == nil {
		msg = "I am already in a good position"
	} else {
		msg = "Getting positioned to help"
	}
	obstacles, _ := strategy.WatchOpponentOnMyRoute(
		d.Player.Coords,
		ballPosition,
		units.PlayerSize,
		d.Player.GetOpponentTeam(d.GameMsg.GameInfo))
	if len(obstacles) > 0 {
		msg = "I need avoid obstacles"
		if optimumSpot, err := physics.NewVector(d.Player.Coords, ballPosition); err == nil {
			detour := optimumSpot.Perpendicular()
			if target == nil {
				target = detour
			} else {
				target.Add(detour)
			}
		}
	}

	stopOrder := d.Player.CreateStopOrder(*d.Player.Velocity.Direction)
	if target == nil {
		return msg, []orders.Order{stopOrder}
	}

	finalPoint := target.TargetFrom(d.Player.Coords)
	finalDistance := finalPoint.DistanceTo(d.Player.Coords)
	if finalDistance < units.PlayerSize && len(obstacles) == 0 {
		return msg, []orders.Order{stopOrder}
	}
	orderMove, err := d.Player.CreateMoveOrderMaxSpeed(finalPoint)
	if err != nil {
		msg = fmt.Sprintf("Sorry, I won't play this turn: %s", err)
		d.Logger.Errorf("error creating order to support: %s", err)
	}

	return msg, []orders.Order{orderMove}
}

func FindSpotToAssist(referencesPositions []physics.Point, player *client.Player) *physics.Vector {
	allVectors := []*physics.Vector{}
	for _, reference := range referencesPositions {
		directionFromRef, err := physics.NewVector(reference, player.Coords)
		if err == nil { //an error means the player is already in the perfect spot
			directionFromRef.SetLength(DistanceBeside) //from all refs, we want to keep "DistanceBeside" distance
			targerFromRef := directionFromRef.TargetFrom(reference)

			vectorToGoodSpot, err := physics.NewVector(player.Coords, targerFromRef)
			if err == nil { //an error means the player is already in the perfect spot from the ball perspective
				allVectors = append(allVectors, vectorToGoodSpot)
			}
		}
	}
	if len(allVectors) == 0 {
		return nil
	}

	final := allVectors[0]
	if len(allVectors) == 1 {
		return final
	}
	for _, vec := range allVectors[1:] {
		final.Add(vec)
	}

	return final
}

// FindSpotToAssist finds a good region to support the ball holder
//func FindSpotToAssist(gameMessage client.GameMessage, assisted *client.Player, assistant *Brain, offensively bool) strategy.RegionCode {
//	var availableSpots []strategy.RegionCode
//	var spotList []strategy.RegionCode
//	if offensively {
//		spotList = ListSpotsCandidatesToOffensiveAssistance(assisted, assistant)
//	} else {
//		spotList = ListSpotsCandidatesToDefensiveAssistance(assisted, assistant)
//	}
//	for _, region := range spotList {
//		mateInTheRegion := assistant.GetPlayersInRegion(region, assistant.GetMyTeamStatus(gameMessage.GameInfo))
//		if len(mateInTheRegion) == 0 {
//			availableSpots = append(availableSpots, region)
//		} else if region == assistant.GetActiveRegion(TeamState) {
//			// eu to no meu canto, me deixe em paz
//			availableSpots = append(availableSpots, region)
//		} else {
//			frankenstein := Brain{Player: mateInTheRegion[0]}
//			isHimTheOwner := region == frankenstein.GetActiveRegion(TeamState)
//			if !isHimTheOwner && assistant.myCurrentRegion() == region {
//				// two invasors disputing
//				myDistanceToTheBall := assistant.Coords.DistanceTo(assisted.Coords)
//				invasorDistanceToTheBall := assistant.Coords.DistanceTo(mateInTheRegion[0].Coords)
//				if myDistanceToTheBall < invasorDistanceToTheBall {
//					availableSpots = append(availableSpots, region)
//				}
//			}
//		}
//	}
//	sort.Slice(availableSpots, func(a, b int) bool {
//		teamStatus := assistant.GetOpponentTeam(gameMessage.GameInfo)
//		opponentsInA := len(assistant.GetPlayersInRegion(availableSpots[a], teamStatus))
//		opponentsInB := len(assistant.GetPlayersInRegion(availableSpots[b], teamStatus))
//
//		distanceToA := math.Round(assistant.Coords.DistanceTo(availableSpots[a].Center(assistant.TeamPlace)) / strategy.RegionWidth)
//		distanceToB := math.Round(assistant.Coords.DistanceTo(availableSpots[b].Center(assistant.TeamPlace)) / strategy.RegionWidth)
//
//		distanceAToAssistant := math.Round(assisted.Coords.DistanceTo(availableSpots[a].Center(assistant.TeamPlace)) / strategy.RegionWidth)
//		distanceBToAssistant := math.Round(assisted.Coords.DistanceTo(availableSpots[b].Center(assistant.TeamPlace)) / strategy.RegionWidth)
//
//		APoints := distanceToB - distanceToA
//		APoints += float64(opponentsInB - opponentsInA)
//		APoints += distanceBToAssistant - distanceAToAssistant
//		APoints += float64(availableSpots[a].X-availableSpots[b].X) * 2.5
//		return APoints >= 0
//	})
//
//	if len(availableSpots) > 0 {
//		return availableSpots[0]
//	}
//	return assistant.GetActiveRegion(TeamState)
//}

// FindBestPointInRegionToAssist finds the best point to support the ball holder from within a region
//func FindBestPointInRegionToAssist(gameMessage client.GameMessage, region strategy.RegionCode, assisted *client.Player) (target physics.Point) {
//	centerPoint := region.Center(assisted.TeamPlace)
//	vctToCenter, err := physics.NewVector(assisted.Coords, centerPoint)
//	if err != nil {
//		SetLength(strategy.RegionWidth)
//	}
//
//	obstacles := watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, vctToCenter.TargetFrom(assisted.Coords))
//	if len(obstacles) == 0 {
//		return vctToCenter.TargetFrom(assisted.Coords)
//	} else {
//		initialVector := vctToCenter
//		avoidObstacles := func(ang float64) bool {
//			tries := 3
//			for tries > 0 {
//				vctToCenter.AddAngleDegree(ang)
//				target = vctToCenter.TargetFrom(assisted.Coords)
//				if region != strategy.GetRegionCode(target, assisted.TeamPlace) {
//					//too far
//					tries = 0
//				}
//				obstacles = watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, target)
//				tries--
//				if len(obstacles) <= 0 {
//					return true
//				}
//			}
//			return false
//		}
//
//		if !avoidObstacles(10) && !avoidObstacles(-10) {
//			target = initialVector.TargetFrom(assisted.Coords)
//		}
//	}
//	return
//}
