package bot

import (
	"context"
	"github.com/lugobots/lugo4go/v2/coach"
)

// You probably DO NOT want to change the functions below

type BotMiddleware func(coach.Decider) coach.Decider

func NewMiddleare() BotMiddleware {
	return func(decider coach.Decider) coach.Decider {
		return middleware{
			next: decider,
		}
	}
}

type middleware struct {
	next coach.Decider
}

func (m middleware) OnDisputing(ctx context.Context, data coach.TurnData) error {
	return m.next.OnDisputing(ctx, data)
}

func (m middleware) OnDefending(ctx context.Context, data coach.TurnData) error {
	return m.next.OnDefending(ctx, data)
}

func (m middleware) OnHolding(ctx context.Context, data coach.TurnData) error {
	return m.next.OnHolding(ctx, data)
}

func (m middleware) OnSupporting(ctx context.Context, data coach.TurnData) error {
	return m.next.OnSupporting(ctx, data)
}

func (m middleware) AsGoalkeeper(ctx context.Context, data coach.TurnData) error {
	return m.next.AsGoalkeeper(ctx, data)
}
