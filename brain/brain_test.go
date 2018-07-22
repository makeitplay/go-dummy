package brain

import (
	"testing"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/go-dummy/strategy"
	"github.com/stretchr/testify/assert"
)

func TestBrain_ShouldIDisputeForTheBall(t *testing.T) {
	msg := Game.GameMessage{}
	msg.GameInfo = Game.GameInfo{}
	msg.GameInfo.Ball = Game.Ball{}
	msg.GameInfo.Ball.Coords = strategy.GetRegionCenter(strategy.RegionCode{0,0}, Units.HomeTeam)


	aldair := Brain{}
	aldair.Player = new(Game.Player)
	aldair.TeamPlace = Units.HomeTeam
	aldair.Number = "aldair"
	aldair.Coords = strategy.GetRegionCenter(strategy.RegionCode{1,1}, Units.HomeTeam)

	bebeto := Brain{}
	bebeto.Player = new(Game.Player)
	bebeto.Number = "bebeto"
	bebeto.TeamPlace = Units.HomeTeam
	bebeto.Coords = strategy.GetRegionCenter(strategy.RegionCode{0,2}, Units.HomeTeam)

	ronaldo := Brain{}
	ronaldo.Player = new(Game.Player)
	ronaldo.TeamPlace = Units.HomeTeam
	ronaldo.Number = "ronaldo"
	ronaldo.Coords = strategy.GetRegionCenter(strategy.RegionCode{2,0}, Units.HomeTeam)

	msg.GameInfo.HomeTeam.Players = []*Game.Player{}
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, aldair.Player)
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, bebeto.Player)
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, ronaldo.Player)
	aldair.LastMsg = msg
	assert.True(t, aldair.ShouldIDisputeForTheBall())

	aldair.Coords = strategy.GetRegionCenter(strategy.RegionCode{2,2}, Units.HomeTeam)
	bebeto.Coords = strategy.GetRegionCenter(strategy.RegionCode{0,1}, Units.HomeTeam)
	ronaldo.Coords = strategy.GetRegionCenter(strategy.RegionCode{3,0}, Units.HomeTeam)
	assert.True(t, aldair.ShouldIDisputeForTheBall())

	aldair.Coords = strategy.GetRegionCenter(strategy.RegionCode{2,2}, Units.HomeTeam)
	bebeto.Coords = strategy.GetRegionCenter(strategy.RegionCode{0,1}, Units.HomeTeam)
	ronaldo.Coords = strategy.GetRegionCenter(strategy.RegionCode{1,0}, Units.HomeTeam)
	assert.False(t, aldair.ShouldIDisputeForTheBall())

	msg.GameInfo.HomeTeam.Players = []*Game.Player{}
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, aldair.Player)
	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, bebeto.Player)
	msg.GameInfo.AwayTeam.Players = []*Game.Player{}
	msg.GameInfo.AwayTeam.Players =  append(msg.GameInfo.AwayTeam.Players, ronaldo.Player)
	aldair.LastMsg = msg
	assert.True(t, aldair.ShouldIDisputeForTheBall())
}
