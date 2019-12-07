package bot

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
)

type Bot struct {
	Positioner      coach.Positioner
	Role            Role
	regionMap       RegionMap
	InitialPosition proto.Point
	log             lugo4go.Logger
}

func NewBot(config lugo4go.Config, logger lugo4go.Logger) (*Bot, error) {
	var err error
	b := Bot{}
	b.log = logger
	b.Positioner, err = coach.NewPositioner(RegionCols, RegionRows, config.TeamSide)
	if err != nil {
		return nil, fmt.Errorf("could not create a positioner: %s", err)
	}

	if config.Number != field.GoalkeeperNumber {
		b.regionMap = DefineRegionMap(config)
		reg, err := b.Positioner.GetRegion(b.regionMap[Initial].Col, b.regionMap[Initial].Row)
		logger.Infof("My position: %d, %v (%v)", config.Number, reg, b.regionMap)
		if err != nil {
			return nil, fmt.Errorf("did not connected to the gRPC server at '%s': %s", config.GRPCAddress, err)
		}

		b.InitialPosition = reg.Center()
		b.Role = DefineRole(config.Number)

	} else {
		b.InitialPosition = field.GetTeamsGoal(config.TeamSide).Center
	}
	return &b, nil
}

func (b Bot) OnDisputing(ctx context.Context, data coach.TurnData) error {
	return myDecider(ctx, data)
}

func (b Bot) OnDefending(ctx context.Context, data coach.TurnData) error {
	return myDecider(ctx, data)
}

func (b Bot) OnHolding(ctx context.Context, data coach.TurnData) error {
	return myDecider(ctx, data)
}

func (b Bot) OnSupporting(ctx context.Context, data coach.TurnData) error {
	return myDecider(ctx, data)
}

func (b Bot) AsGoalkeeper(ctx context.Context, data coach.TurnData) error {
	return myDecider(ctx, data)
}

func myDecider(ctx context.Context, data coach.TurnData) error {
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	if field.IsBallHolder(data.Snapshot, data.Me) {
		orderToKick, err := field.MakeOrderKick(*data.Snapshot.Ball, field.GetOpponentGoal(data.Me.TeamSide).Center, field.BallMaxSpeed)
		if err != nil {
			return fmt.Errorf("could not create kick order during turn %d: %s", data.Snapshot.Turn, err)
		}
		orders = []proto.PlayerOrder{orderToKick}
	} else if data.Me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := field.MakeOrderMoveMaxSpeed(*data.Me.Position, *data.Snapshot.Ball.Position)
		if err != nil {
			return fmt.Errorf("could not create move order during turn %d: %s", data.Snapshot.Turn, err)
		}
		orders = []proto.PlayerOrder{orderToMove, field.MakeOrderCatch()}
	} else {
		orders = []proto.PlayerOrder{field.MakeOrderCatch()}
	}

	resp, err := data.Sender.Send(ctx, orders, "")
	if err != nil {
		return fmt.Errorf("could not send kick order during turn %d: %s", data.Snapshot.Turn, err)
	} else if resp.Code != proto.OrderResponse_SUCCESS {
		return fmt.Errorf("order sent not  order during turn %d: %s", data.Snapshot.Turn, err)
	}
	return nil
}
