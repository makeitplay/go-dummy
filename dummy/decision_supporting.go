package dummy

import (
	"fmt"
	"github.com/lugobots/arena/orders"
	"github.com/lugobots/arena/physics"
	"github.com/lugobots/arena/units"
	"github.com/lugobots/client-player-go"
	"github.com/lugobots/the-dummies-go/strategy"
)

func (d *Dummy) orderForSupporting() (msg string, orders []orders.Order) {
	if d.ShouldIAssist() { // middle players will give support
		return d.orderForActiveSupport()
	}
	return d.orderForPassiveSupport()
}

func (d *Dummy) ShouldIAssist() bool {
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

func (d *Dummy) orderForActiveSupport() (msg string, ordersSet []orders.Order) {
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

func (d *Dummy) orderForPassiveSupport() (msg string, ordersSet []orders.Order) {
	player := d.Player
	var region strategy.RegionCode
	region = d.GetActiveRegion()

	target := region.Center(TeamPlace)
	if player.Coords.DistanceTo(target) < units.PlayerMaxSpeed {
		if player.Velocity.Speed > 0 {
			stopOrder := d.Player.CreateStopOrder(*d.Player.Velocity.Direction)
			ordersSet = []orders.Order{stopOrder}
		}
	} else {
		orderMove, err := d.Player.CreateMoveOrderMaxSpeed(target)
		if err != nil {
			msg = fmt.Sprintf("Sorry, I won't play this turn: %s", err)
			d.Logger.Errorf("error creating order to passive support: %s", err)
		} else {
			ordersSet = []orders.Order{orderMove}
		}
	}
	return msg, ordersSet
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
