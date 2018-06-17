package brain

import (
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/GameState"
	"math"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/BasicTypes"
)
// distance considered "near" for a player to the ball
const DistanceNearBall = Units.CourtHeight / 2 // units float
const ERROR_MARGIN_RUNNING = 20.0
const ERROR_MARGIN_PASSING = 20.0


type Brain struct {
	*Game.Player
	State          PlayerState
}

func (b *Brain) ResetPosition() {
	region := b.myRegion()
	b.Coords = region.InitialPosition()
}

func (b *Brain) ProcessAnn(msg Game.GameMessage) {
	b.UpdatePosition(msg.GameInfo)
	commons.LogBroadcast("ANN %s", string(msg.State))
	switch GameState.State(msg.State) {
	case GameState.GETREADY:
	case GameState.LISTENING:
		b.State = b.DetermineMyState()
		commons.LogDebug("State: %s", b.State)
		b.TakeAnAction()
	}
}

func (b *Brain) DetermineMyState() PlayerState {
	var isOnMyField bool
	var subState string
	var ballPossess string

	if b.LastMsg.GameInfo.Ball.Holder == nil {
		ballPossess = "dsp" //disputing
		subState = "fbl"    //far
		if int(math.Abs(b.Coords.DistanceTo(b.LastMsg.GameInfo.Ball.Coords))) <= DistanceNearBall {
			subState = "nbl" //near
		}
	} else if b.LastMsg.GameInfo.Ball.Holder.TeamPlace == b.TeamPlace {
		ballPossess = "atk" //attacking
		subState = "hlp"    //helping
		if b.LastMsg.GameInfo.Ball.Holder.Id == b.Id {
			subState = "hld" //holdin
		}
	} else {
		ballPossess = "dfd"
		subState = "org"
		if b.isItInMyRegion(b.LastMsg.GameInfo.Ball.Coords) {
			subState = "mrg"
		}
	}

	if b.TeamPlace == Units.HomeTeam {
		isOnMyField = b.LastMsg.GameInfo.Ball.Coords.PosX <= Units.CourtWidth/2
	} else {
		isOnMyField = b.LastMsg.GameInfo.Ball.Coords.PosX >= Units.CourtWidth/2
	}
	fieldState := "fr"
	if isOnMyField {
		fieldState = "hs"
	}
	return PlayerState(ballPossess + "-" + subState + "-" + fieldState)
}

func (b *Brain)  TakeAnAction() {
	var orders []BasicTypes.Order
	var msg string

	switch b.State {
	case AtckHoldHse:
		msg, orders = b.orderForAtckHoldHse()
	case AtckHoldFrg:
		msg, orders = b.orderForAtckHoldFrg()
	case AtckHelpHse:
		msg, orders = b.orderForAtckHelpHse()
	case AtckHelpFrg:
		msg, orders = b.orderForAtckHelpFrg()
	case DefdMyrgHse:
		msg, orders = b.orderForDefdMyrgHse()
		orders = append(orders, b.CreateCatchOrder())
	case DefdMyrgFrg:
		msg, orders = b.orderForDefdMyrgFrg()
		orders = append(orders, b.CreateCatchOrder())
	case DefdOtrgHse:
		msg, orders = b.orderForDefdOtrgHse()
		orders = append(orders, b.CreateCatchOrder())
	case DefdOtrgFrg:
		msg, orders = b.orderForDefdOtrgFrg()
		orders = append(orders, b.CreateCatchOrder())
	case DsptNfblHse:
		msg, orders = b.orderForDsptNfblHse()
		orders = append(orders, b.CreateCatchOrder())
	case DsptNfblFrg:
		msg, orders = b.orderForDsptNfblFrg()
		orders = append(orders, b.CreateCatchOrder())
	case DsptFrblHse:
		msg, orders = b.orderForDsptFrblHse()
		orders = append(orders, b.CreateCatchOrder())
	case DsptFrblFrg:
		msg, orders = b.orderForDsptFrblFrg()
		orders = append(orders, b.CreateCatchOrder())
	}
	commons.LogDebug("Sending order")
	b.SendOrders(msg, orders...)

}