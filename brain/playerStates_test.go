package brain

import (
	"testing"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/go-dummy/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/commons/Units"
)

func TestBrain_orderForAtckHelpFrg(t *testing.T) {
	player := new(Brain)
	player.Player = new(Game.Player)
	player.Number = BasicTypes.PlayerNumber("8")
	player.Coords = strategy.RegionCode{4,2}.Center(Units.HomeTeam)
	MyRule = strategy.MiddlePlayer

	holder := new(Brain)
	holder.Player = new(Game.Player)
	holder.Number = BasicTypes.PlayerNumber("7")
	holder.Coords = strategy.RegionCode{4,1}.Center(Units.HomeTeam)

	lastMsg := MountGameInfo()
	lastMsg.GameInfo.HomeTeam.Players = append(lastMsg.GameInfo.HomeTeam.Players, player.Player, holder.Player)

	lastMsg.GameInfo.Ball.Holder = holder.Player
	lastMsg.GameInfo.Ball.Coords = holder.Coords

	TeamState = player.DetermineMyTeamState(lastMsg)

	player.LastMsg = lastMsg
	holder.LastMsg = lastMsg

	msg, order := player.orderForAtckHelpFrg()
	assert.Equal(t, "", msg)
	assert.Equal(t, string(strategy.Offensive), string(TeamState))
	assert.Len(t, order, 1)
	assert.Len(t, order[0], 1)
}
