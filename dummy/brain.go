package dummy

import (
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/sirupsen/logrus"
)

var ClientResponder client.Responder
var TeamPlace arena.TeamPlace
var PlayerNumber arena.PlayerNumber
var MyRule strategy.PlayerRule
var LastHolderFrom *client.Player

// TeamBallPossession stores the team's name that has touched on the ball for the last time
var TeamBallPossession arena.TeamPlace

type Dummy struct {
	GameMsg     *client.GameMessage
	Player      *client.Player
	TeamState   strategy.TeamState
	PlayerState strategy.PlayerState
	Logger      *logrus.Entry
}

func (d *Dummy) React() {
	var ordersSet []orders.Order
	var msg string
	if d.Player.IsGoalkeeper() {
		msg, ordersSet = d.orderForGoalkeeper()
	} else {
		switch d.PlayerState {
		case strategy.DisputingTheBall:
			msg, ordersSet = d.orderForDisputingTheBall()
			if !Passing {
				ordersSet = append(ordersSet, d.Player.CreateCatchOrder())
			}
		case strategy.Supporting:
			Passing = false
			msg, ordersSet = d.orderForSupporting()
		case strategy.HoldingTheBall:
			Passing = false
			msg, ordersSet = d.orderForHoldingTheBall()
		case strategy.Defending:
			Passing = false
			msg, ordersSet = d.orderForDefending()
			ordersSet = append(ordersSet, d.Player.CreateCatchOrder())
		}
	}

	ClientResponder.SendOrders(fmt.Sprintf("%s %s", d.Player.ID(), msg), ordersSet...)
}

// @todo Needs enhancement: the player does not consider the position of the other supporters, so if two players are behind the opponent it does not try to help
func (d *Dummy) ShouldIDisputeForTheBall() bool {
	if strategy.GetRegionCode(d.GameMsg.GameInfo.Ball.Coords, TeamPlace).ChessDistanceTo(d.GetActiveRegion()) < 2 {
		return true
	}
	myDistance := d.Player.Coords.DistanceTo(d.GameMsg.GameInfo.Ball.Coords)
	playerCloser := 0
	for _, teamMate := range d.Player.GetMyTeamStatus(d.GameMsg.GameInfo).Players {
		if teamMate.Number != PlayerNumber && teamMate.Coords.DistanceTo(d.GameMsg.GameInfo.Ball.Coords) < myDistance {
			playerCloser++
			if playerCloser > 1 { // are there more than on player closer to the ball than me?
				return false
			}
		}
	}
	return true
}
