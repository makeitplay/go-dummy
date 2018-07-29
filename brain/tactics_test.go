package brain

import (
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBrain_GetActiveRegion(t *testing.T) {
	aldair := Brain{}
	aldair.Player = new(client.Player)
	aldair.Number = "3"
	assert.Equal(t, strategy.RegionCode{0, 1}, aldair.GetActiveRegion(strategy.UnderPressure))
	assert.Equal(t, strategy.RegionCode{1, 1}, aldair.GetActiveRegion(strategy.Defensive))

	aldair.Number = "8"
	assert.Equal(t, strategy.RegionCode{1, 2}, aldair.GetActiveRegion(strategy.UnderPressure))
	assert.Equal(t, strategy.RegionCode{5, 2}, aldair.GetActiveRegion(strategy.OnAttack))
}

func TestBrain_myCurrentRegion(t *testing.T) {
	aldair := Brain{}
	aldair.Player = new(client.Player)
	aldair.TeamPlace = Units.HomeTeam
	aldair.Number = "3"

	aldair.Coords = Physics.Point{
		PosX: strategy.RegionWidth * 3.4,  //3
		PosY: strategy.RegionHeight * 2.1, //2
	}
	assert.Equal(t, strategy.RegionCode{3, 2}, aldair.myCurrentRegion())

	aldair.Coords = Physics.Point{
		PosX: strategy.RegionWidth * 4.2,  //4
		PosY: strategy.RegionHeight * 0.1, //0
	}
	assert.Equal(t, strategy.RegionCode{4, 0}, aldair.myCurrentRegion())

	aldair.TeamPlace = Units.AwayTeam

	aldair.Coords = Physics.Point{
		PosX: strategy.RegionWidth * 4.2,  //3
		PosY: strategy.RegionHeight * 0.1, //3
	}
	assert.Equal(t, strategy.RegionCode{3, 3}, aldair.myCurrentRegion())

}

func TestBrain_isItInMyActiveRegion(t *testing.T) {
	aldair := Brain{}
	aldair.Player = new(client.Player)
	aldair.TeamPlace = Units.HomeTeam
	aldair.Number = "5"

	assert.True(t, aldair.isItInMyActiveRegion(Physics.Point{
		PosX: strategy.RegionWidth * 4.2,
		PosY: strategy.RegionHeight * 3.1,
	}, strategy.OnAttack))

	assert.False(t, aldair.isItInMyActiveRegion(Physics.Point{
		PosX: strategy.RegionWidth * 3.2,
		PosY: strategy.RegionHeight * 3.1,
	}, strategy.OnAttack))

	assert.False(t, aldair.isItInMyActiveRegion(Physics.Point{
		PosX: strategy.RegionWidth * 4.2,
		PosY: strategy.RegionHeight * 2.1,
	}, strategy.OnAttack))

}

func TestDetermineMyTeamState_NoBall(t *testing.T) {
	msg := client.GameMessage{}
	msg.GameInfo = client.GameInfo{}
	msg.GameInfo.Ball = client.Ball{}

	homePlayer := new(Brain)
	homePlayer.Player = new(client.Player)
	homePlayer.TeamPlace = Units.HomeTeam

	awayPlayer := new(Brain)
	awayPlayer.Player = new(client.Player)
	awayPlayer.TeamPlace = Units.AwayTeam

	TeamBallPossession = Units.AwayTeam
	msg.GameInfo.Ball.Coords = strategy.RegionCode{0, 2}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.UnderPressure, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.OnAttack, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{1, 1}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.UnderPressure, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.OnAttack, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{2, 1}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.Defensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Offensive, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{3, 1}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.Defensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Offensive, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{4, 1}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.Neutral, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{5, 1}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.Neutral, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{6, 1}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.Offensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{7, 1}.Center(Units.HomeTeam)
	assert.Equal(t, strategy.Offensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(msg))

	TeamBallPossession = Units.HomeTeam
	msg.GameInfo.Ball.Coords = strategy.RegionCode{0, 2}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.OnAttack, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.UnderPressure, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{1, 1}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.OnAttack, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.UnderPressure, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{2, 1}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.Offensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{3, 1}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.Offensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{4, 1}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.Neutral, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{5, 1}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.Neutral, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{6, 1}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.Defensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Offensive, awayPlayer.DetermineMyTeamState(msg))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{7, 1}.Center(Units.AwayTeam)
	assert.Equal(t, strategy.Defensive, homePlayer.DetermineMyTeamState(msg))
	assert.Equal(t, strategy.Offensive, awayPlayer.DetermineMyTeamState(msg))

}
