package dummy

import (
	"fmt"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"math/rand"
	"time"
)

func (d *Dummy) orderForSupporting() (msg string, orders []orders.Order) {
	if d.ShouldIAssist() { // middle players will give support
		if rand.Int()%100 < 1 {
			question := client.TrainingQuestion{
				Question:   "Where should I go?",
				QuestionId: fmt.Sprintf("%s-%s", d.Player.Id, time.Now()),
				PlayerId:   d.Player.ID(),
				Alternatives: []string{
					"front",
					"right",
					"left",
					"back",
					"stay",
					"ignore",
				},
			}
			if err := client.AskQuestion(question, *GameConfig); err == nil {
				d.Logger.Warn("question sent")
				TunnelMsg = make(chan client.GameMessage)
				WaitingAnswer = true
				ds, err := DS.SaveSample(d.GameMsg.GameInfo)
				if err != nil {
					d.Logger.Errorf("did not create the state: %s", err )
				}
				d.Logger.Warnf("Bora esperar! ")
				var answer string
				for WaitingAnswer {
					select {
					case debugMsg := <-TunnelMsg:
						var ok bool
						if answer, ok = debugMsg.Data[question.QuestionId].(string); ok {
							d.Logger.Warnf("Got the answer! %s", answer)
							d.Logger.Warnf("Recebeu")
							WaitingAnswer = false
						} else {
							d.Logger.Warnf("Not yet")
						}
					}
				}

				if answer != "ignore" {
					ds.Save(answer)
				}
			}
		}
		return d.orderForActiveSupport()
	}


	return d.orderForPassiveSupport()
}

func (d *Dummy) ShouldIAssist() bool {
	myDistance := d.Player.Coords.DistanceTo(d.GameMsg.GameInfo.Ball.Holder.Coords)
	holderId := d.GameMsg.GameInfo.Ball.Holder.ID()
	playerCloser := 0

	shouldAssist := true
	// should have at least 3 supporters in the perimeters around the ball holder
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
	references := []physics.Point{ballPosition}

	d.GameMsg.ForEachPlayByTeam(TeamPlace, func(index int, player *client.Player) {
		if player.ID() != holderId && // the holder cannot help himself
			player.Number != PlayerNumber && // I wont count to myself
			//player.Coords.DistanceTo(d.GameMsg.GameInfo.Ball.Coords) <= DistanceNear && //the guys who are already supporting
			player.Coords.DistanceTo(d.Player.Coords) < DistanceBeside {
			references = append(references, player.Coords)
		}
	})

	target := FindSpotToAssist(references, d.Player)

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


	if target == nil {
		advance, _ := d.Player.CreateMoveOrderMaxSpeed(d.Player.OpponentGoal().Center)
		return msg, []orders.Order{advance}
	}

	stopOrder := d.Player.CreateStopOrder(*d.Player.Velocity.Direction)
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

	target := d.GetActiveRegionCenter()
	if target.DistanceTo(d.GameMsg.Ball().Coords) > DistanceDistant {
		d.Logger.Errorf("---- I AM TOOOO FAR")
		pathToGoal, err := physics.NewVector(target, d.GameMsg.Ball().Coords)
		if err == nil {
			pathToGoal.SetLength(DistanceFar)
			target = pathToGoal.TargetFrom(target)
		}
	}
	//target = target.MiddlePointTo(d.GameMsg.Ball().Coords)
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
