package main

import (
	"math/rand"
	"time"

	"github.com/makeitplay/arena"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/dummie"
	"github.com/makeitplay/the-dummies-go/strategy"
)

func main() {
	rand.Seed(time.Now().Unix())
	serverConfig := new(client.Configuration)
	serverConfig.ParseFromFlags()

	gamer := &client.Gamer{}

	dummie.PlayerNumber = serverConfig.PlayerNumber
	dummie.TeamPlace = serverConfig.TeamPlace
	dummie.MyRule = strategy.DefinePlayerRule(serverConfig.PlayerNumber)
	dummie.TeamBallPossession = dummie.TeamPlace
	dummie.ClientResponder = gamer

	gamer.OnAnnouncement = reactToNewState
	gamer.Play(dummie.GetInitialRegion().Center(serverConfig.TeamPlace), serverConfig)
}

func reactToNewState(ctx client.TurnContext) {

	switch ctx.GameMsg().GameInfo.State {
	case arena.Listening:
		if ctx.GameMsg().Ball().Holder != nil {
			dummie.TeamBallPossession = ctx.GameMsg().Ball().Holder.TeamPlace
		}

		dummy := &dummie.Dummie{
			GameMsg:     ctx.GameMsg(),
			Player:      ctx.Player(),
			PlayerState: strategy.DetermineMyState(ctx),
			TeamState:   strategy.DetermineMyTeamState(ctx, dummie.TeamBallPossession),
			Logger:      ctx.Logger(),
		}

		ctx.Logger().Infof("my state: %s", dummy.PlayerState)
		dummy.React()
		dummie.LastHolderFrom = ctx.GameMsg().Ball().Holder
	}
}
