package dummie

import (
	"github.com/makeitplay/arena/orders"
)

func (d *Dummie) orderForDefending() (msg string, ordersSet []orders.Order) {
	return d.orderForDisputingTheBall()
}
