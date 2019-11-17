package dummy

import (
	"github.com/lugobots/arena/orders"
)

func (d *Dummy) orderForDefending() (msg string, ordersSet []orders.Order) {
	return d.orderForDisputingTheBall()
}
