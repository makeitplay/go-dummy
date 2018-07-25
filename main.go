package main

import (
	"time"
	"os"
	"os/signal"
	"math/rand"

	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/the-dummies-go/brain"
	"github.com/makeitplay/the-dummies-go/strategy"
)

func main() {
	rand.Seed(time.Now().Unix())
	watchInterruptions()
	defer commons.Cleanup(false)
	serverConfig := new(Game.Configuration)
	commons.Load(serverConfig)
	serverConfig.LoadCmdArg()
	/**********************************************/

	brain.MyRule = strategy.DefinePlayerRule(serverConfig.PlayerNumber)
	brain.TeamBallPossession = serverConfig.TeamPlace

	player := &Game.Player{}
	playerBrain := &brain.Brain{Player: player}
	playerBrain.TeamPlace = serverConfig.TeamPlace
	playerBrain.Number = serverConfig.PlayerNumber
	playerBrain.Size = Units.PlayerSize
	playerBrain.ResetPosition()
	playerBrain.Player.OnAnnouncement = playerBrain.ProcessAnn
	playerBrain.Start(serverConfig)
}

func watchInterruptions() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			commons.Log("*********** INTERRUPTION SIGNAL ****************")
			commons.Cleanup(true)
			os.Exit(0)
		}
	}()
}
