package brain

import (
	"fmt"
	"math"

	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"

	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
)

// TeamState stores the team state based on our strategy
var TeamState = strategy.Defensive

// TeamBallPossession stores the team's name that has touched on the ball for the last time
var TeamBallPossession arena.TeamPlace

// MyRule stores this player rule in the team
var MyRule strategy.PlayerRule

// Brain controls the player to have a behaviour during each state
type Brain struct {
	TeamPlace arena.TeamPlace
	Number    arena.PlayerNumber
	State     PlayerState
	Responser client.Responder
}

// ProcessAnn is the callback function called when the player gets a new message from the game server
func (b *Brain) ProcessAnn(turn client.TurnContext) {

}

// DetermineMyState determine the player state bases on our strategy
func (b *Brain) DetermineMyState(turn client.TurnContext) PlayerState {

}

// TakeAnAction sends orders to the game server based on the player state
func (b *Brain) TakeAnAction(turn client.TurnContext) {

}

// ShouldIDisputeForTheBall returns true when the player should try to catch the ball
func (b *Brain) ShouldIDisputeForTheBall(turn client.TurnContext) bool {

}

// ShouldIAssist returns the ball when the player should support another team mate
func (b *Brain) ShouldIAssist(turn client.TurnContext) bool {
}

// FindBestPointShootTheBall calculates the best point in the goal to shoot the ball
func (b *Brain) FindBestPointShootTheBall(turn client.TurnContext) (speed float64, target physics.Point) {
	goalkeeper := b.FindOpponentPlayer(turn.GameMsg().GameInfo, BasicTypes.PlayerNumber("1"))
	if goalkeeper.Coords.PosY > units.FieldHeight/2 {
		return units.BallMaxSpeed, physics.Point{
			PosX: b.OpponentGoal().BottomPole.PosX,
			PosY: b.OpponentGoal().BottomPole.PosY + units.BallSize,
		}
	} else {
		return units.BallMaxSpeed, physics.Point{
			PosX: b.OpponentGoal().TopPole.PosX,
			PosY: b.OpponentGoal().TopPole.PosY - units.BallSize,
		}
	}
}
