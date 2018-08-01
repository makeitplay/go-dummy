package brain

import (
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBrain_ShouldIDisputeForTheBall(t *testing.T) {
	msg := client.GameMessage{}
	msg.GameInfo = client.GameInfo{}
	msg.GameInfo.Ball = client.Ball{}
	msg.GameInfo.Ball.Coords = (strategy.RegionCode{0, 0}).Center(Units.HomeTeam)

	aldair := Brain{}
	aldair.Player = new(client.Player)
	aldair.TeamPlace = Units.HomeTeam
	aldair.Number = "7"
	aldair.Coords = (strategy.RegionCode{1, 1}).Center(Units.HomeTeam)

	bebeto := Brain{}
	bebeto.Player = new(client.Player)
	bebeto.Number = "8"
	bebeto.TeamPlace = Units.HomeTeam
	bebeto.Coords = (strategy.RegionCode{0, 2}).Center(Units.HomeTeam)

	ronaldo := Brain{}
	ronaldo.Player = new(client.Player)
	ronaldo.TeamPlace = Units.HomeTeam
	ronaldo.Number = "9"
	ronaldo.Coords = (strategy.RegionCode{2, 0}).Center(Units.HomeTeam)

	msg.GameInfo.HomeTeam.Players = []*client.Player{}
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, aldair.Player)
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, bebeto.Player)
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, ronaldo.Player)
	aldair.LastMsg = msg
	assert.True(t, aldair.ShouldIDisputeForTheBall())

	aldair.Coords = (strategy.RegionCode{2, 2}).Center(Units.HomeTeam)
	bebeto.Coords = (strategy.RegionCode{0, 1}).Center(Units.HomeTeam)
	ronaldo.Coords = (strategy.RegionCode{3, 0}).Center(Units.HomeTeam)
	assert.True(t, aldair.ShouldIDisputeForTheBall())

	aldair.Coords = (strategy.RegionCode{2, 2}).Center(Units.HomeTeam)
	bebeto.Coords = (strategy.RegionCode{0, 1}).Center(Units.HomeTeam)
	ronaldo.Coords = (strategy.RegionCode{1, 0}).Center(Units.HomeTeam)
	assert.False(t, aldair.ShouldIDisputeForTheBall())

	msg.GameInfo.HomeTeam.Players = []*client.Player{}
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, aldair.Player)
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, bebeto.Player)
	msg.GameInfo.AwayTeam.Players = []*client.Player{}
	msg.GameInfo.AwayTeam.Players = append(msg.GameInfo.AwayTeam.Players, ronaldo.Player)
	aldair.LastMsg = msg
	assert.True(t, aldair.ShouldIDisputeForTheBall())
}
