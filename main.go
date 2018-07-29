package main

import (
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/commons"
	"github.com/makeitplay/the-dummies-go/brain"
	"github.com/makeitplay/the-dummies-go/strategy"
)

func main() {
	rand.Seed(time.Now().Unix())
	watchInterruptions()
	defer commons.Cleanup(false)
	serverConfig := new(client.Configuration)
	serverConfig.LoadCmdArg()

	brain.MyRule = strategy.DefinePlayerRule(serverConfig.PlayerNumber)
	brain.TeamBallPossession = serverConfig.TeamPlace

	player := &client.Player{}
	playerBrain := &brain.Brain{Player: player}
	playerBrain.TeamPlace = serverConfig.TeamPlace
	playerBrain.Number = serverConfig.PlayerNumber
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
