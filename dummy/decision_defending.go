package dummy

import (
	"github.com/makeitplay/arena/orders"
)

func (d *Dummy) orderForDefending() (msg string, ordersSet []orders.Order) {
	return d.orderForDisputingTheBall()
}
