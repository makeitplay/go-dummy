package dummie

import (
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/the-dummies-go/strategy"
)

func (d *Dummie) orderForDisputingTheBall() (msg string, ordersSet []orders.Order) {
	player := d.Player
	if d.ShouldIDisputeForTheBall() {
		msg = "Disputing for the ball"
		//orders = []orders.Order{d.CreateMoveOrderMaxSpeed(d.LastMsg.GameInfo.Ball.Coords)}
		speed, target := strategy.FindBestPointInterceptBall(d.GameMsg.Ball(), player)
		movOrder, err := player.CreateMoveOrder(target, speed)
		if err != nil {
			d.Logger.Errorf("error creating move order: %s ", err)
			msg = "sorry, I won't play this turn"
		} else {
			ordersSet = []orders.Order{movOrder}
		}
	} else {
		if d.myCurrentRegion() != d.GetActiveRegion() {
			movOrder, err := player.CreateMoveOrderMaxSpeed(d.GetActiveRegionCenter())
			if err != nil {
				d.Logger.Errorf("error creating move order: %s ", err)
				msg = "sorry, I won't play this turn"
			} else {
				msg = "Moving to my region"
				ordersSet = []orders.Order{movOrder}
			}
		} else {
			msg = "Holding position"
			ordersSet = []orders.Order{player.CreateStopOrder(*player.Velocity.Direction)}
		}
	}
	return msg, ordersSet
}

func (d *Dummie) ShouldIDisputeForTheBall() bool {
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
