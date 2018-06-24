package strategy

import (
	"testing"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/stretchr/testify/assert"
	"github.com/makeitplay/commons/Units"
)

func TestDetermineTeamState_NoBall(t *testing.T) {
	msg := Game.GameMessage{}
	msg.GameInfo = Game.GameInfo{}
	msg.GameInfo.Ball = Game.Ball{}

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{0,2}, Units.HomeTeam)
	assert.Equal(t, UnderPressure, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{1,1}, Units.HomeTeam)
	assert.Equal(t, UnderPressure, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{2,1}, Units.HomeTeam)
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{3,0}, Units.HomeTeam)
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{4,2}, Units.HomeTeam)
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{5,2}, Units.HomeTeam)
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{6,0}, Units.HomeTeam)
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, UnderPressure, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{7,0}, Units.HomeTeam)
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, UnderPressure, DetermineTeamState(msg, Units.AwayTeam))

}
func TestDetermineTeamState_WithBall(t *testing.T) {
	msg := Game.GameMessage{}
	msg.GameInfo = Game.GameInfo{}
	msg.GameInfo.Ball = Game.Ball{}
	msg.GameInfo.Ball.Holder = &Game.Player{
		TeamPlace: Units.HomeTeam,
	}

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{0,2}, Units.HomeTeam)
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{1,1}, Units.HomeTeam)
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{2,1}, Units.HomeTeam)
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{3,0}, Units.HomeTeam)
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Neutral, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{4,2}, Units.HomeTeam)
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{5,2}, Units.HomeTeam)
	assert.Equal(t, Offensive, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, Defensive, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{6,0}, Units.HomeTeam)
	assert.Equal(t, OnAttack, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, UnderPressure, DetermineTeamState(msg, Units.AwayTeam))

	msg.GameInfo.Ball.Coords = GetRegionCenter(RegionCode{7,0}, Units.HomeTeam)
	assert.Equal(t, OnAttack, DetermineTeamState(msg, Units.HomeTeam))
	assert.Equal(t, UnderPressure, DetermineTeamState(msg, Units.AwayTeam))

}
