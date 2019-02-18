package main

import (
	"math/rand"
	"time"

	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/brain"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/sirupsen/logrus"
)

func main() {
	rand.Seed(time.Now().Unix())
	serverConfig := new(client.Configuration)
	serverConfig.ParseFromFlags()

	brain.MyRule = strategy.DefinePlayerRule(serverConfig.PlayerNumber)
	brain.TeamBallPossession = serverConfig.TeamPlace

	player := &client.Player{}
	playerBrain := &brain.Brain{Player: player}
	playerBrain.TeamPlace = serverConfig.TeamPlace
	playerBrain.Number = serverConfig.PlayerNumber
	playerBrain.ResetPosition()
	playerBrain.Player.OnAnnouncement = playerBrain.ProcessAnn
	logger := logrus.New()

	playerBrain.Start(logger, serverConfig)
}
