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
