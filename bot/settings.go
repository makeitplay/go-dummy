package bot

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
)

// IMPORTANT: all this constant sets below may be changed (see each set instructions). However, any change will
// affect the tactic defined in tactic.go file. So you must go there and adapt your tactics to your new settings.

type TeamState string

type Role string

// Do not remove, rename, or add constants here.
// You however, may increase or decrease their values to change the precision of the Positioner.
// These values defines how the field will be divided by the Positioner to create a field map.
const (
	RegionCols = 8
	RegionRows = 4
)

// please update the tests if you include more states, or exclude some of them.
const (
	Initial       TeamState = "initial"
	UnderPressure TeamState = "under-pressure"
	Defensive     TeamState = "defensive"
	Neutral       TeamState = "neutral"
	Offensive     TeamState = "offensive"
	OnAttack      TeamState = "on-attack"
)
const (
	Defense Role = "defense"
	Middle  Role = "middle"
	Attack  Role = "attack"
)

type FuzzyScale int

const (
	MustNot FuzzyScale = iota
	ShouldNot
	May
	Should
	Must
)

const (
	DistanceNear = field.FieldWidth / 8
	DistanceFar  = DistanceNear * 3
)

type RegionCode struct {
	Col uint8
	Row uint8
}

type RegionMap map[TeamState]RegionCode

func DefineRegionMap(number uint32) RegionMap {
	return roleMap[number]
}

func send(ctx context.Context, turn coach.TurnData, orders []proto.PlayerOrder, debugMsg string) error {
	r, err := turn.Sender.Send(ctx, orders, debugMsg)
	if err != nil {
		return fmt.Errorf("error sending the orders: %s", err)
	}
	if r.Code != proto.OrderResponse_SUCCESS {
		return fmt.Errorf("game server returned an error on our order: %s", err)
	}
	return nil
}

//func GetBallRegion(positioner coach.Positioner, ball proto.Ball, logger lugo4go.Logger) coach.Region  {
//	reg, err := positioner.GetPointRegion(*ball.Position)
//	if err != nil {
//		logger.Errorf("could not find the ball region: %s", err)
//		return nil
//	}
//	return reg
//}

func GetMyRegion(teamState TeamState, positioner coach.Positioner, number uint32) coach.Region {
	regCode := DefineRegionMap(number)[teamState]
	// @see Readme- ignored errors
	r, _ := positioner.GetRegion(regCode.Col, regCode.Row)
	return r
}
