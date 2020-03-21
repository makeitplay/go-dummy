package main

import (
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/the-dummies-go/coach"
	"math/rand"
	"time"

	"github.com/makeitplay/arena"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/dummy"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
)

func main() {
	rand.Seed(time.Now().Unix())
	serverConfig := new(client.Configuration)
	serverConfig.ParseFromFlags()
	serverConfig.LogLevel = logrus.DebugLevel
	dummy.GameConfig = serverConfig
	gamer := &client.Gamer{}

	dummy.PlayerNumber = serverConfig.PlayerNumber
	dummy.TeamPlace = serverConfig.TeamPlace
	dummy.MyRule = strategy.DefinePlayerRule(serverConfig.PlayerNumber)
	dummy.TeamBallPossession = dummy.TeamPlace
	dummy.ClientResponder = gamer

	//gamer.OnAnnouncement = reactToNewState

	gamerCtx, err := gamer.Play(dummy.GetInitialRegion().Center(serverConfig.TeamPlace), serverConfig)
	if err != nil {
		log.Fatal(err)
	}
	dummy.DS = coach.NewImageBasedDataSever("supporting", dummy.TeamPlace)
	gamer.OnMessage = func(msg client.GameMessage) {
		go func() {
			logrus.Warnf("SOME MSG MSG! %v (%v)", msg.Type, msg.Data)
			switch msg.Type {
			case orders.ANSWER:
				logrus.Warnf("DEBUG MSG! %v (%v)", dummy.WaitingAnswer, msg.Data)
				if dummy.WaitingAnswer {
					dummy.TunnelMsg <- msg
				} else {
					logrus.Warnf("Not for me")
				}
			case orders.ANNOUNCEMENT:
				reactToNewState(gamerCtx.CreateTurnContext(msg))
			case orders.RIP:
				gamer.StopToPlay(true)
			}
		}()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-signalChan:
		logrus.Print("*********** INTERRUPTION SIGNAL ****************")
		gamer.StopToPlay(true)
	case <-gamerCtx.Done():
		logrus.Print("*********** Game stopped ****************")
	}

}

func reactToNewState(ctx client.TurnContext) {

	switch ctx.GameMsg().GameInfo.State {
	case arena.Listening:
		if ctx.GameMsg().Ball().Holder != nil {
			dummy.TeamBallPossession = ctx.GameMsg().Ball().Holder.TeamPlace
		}

		ctx.Player().Velocity.Add()

		player := &dummy.Dummy{
			GameMsg:     ctx.GameMsg(),
			Player:      ctx.Player(),
			PlayerState: strategy.DetermineMyState(ctx),
			TeamState:   strategy.DetermineMyTeamState(ctx, dummy.TeamBallPossession),
			Logger:      ctx.Logger(),
		}

		ctx.Logger().Infof("my state: %s", player.PlayerState)
		player.React()
		dummy.LastHolderFrom = ctx.GameMsg().Ball().Holder
	}
}
