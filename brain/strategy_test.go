package brain

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/makeitplay/go-dummy/strategy"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
)

func TestBrain_GetActiveRegion(t *testing.T) {
	aldair := Brain{}
	aldair.Player = new(Game.Player)
	aldair.Number = "3"
	assert.Equal(t, strategy.RegionCode{0,1}, aldair.GetActiveRegion(strategy.UnderPressure))
	assert.Equal(t, strategy.RegionCode{1,1}, aldair.GetActiveRegion(strategy.Defensive))

	aldair.Number = "8"
	assert.Equal(t, strategy.RegionCode{1,2}, aldair.GetActiveRegion(strategy.UnderPressure))
	assert.Equal(t, strategy.RegionCode{5,2}, aldair.GetActiveRegion(strategy.OnAttack))
}

func TestBrain_myCurrentRegion(t *testing.T) {
	aldair := Brain{}
	aldair.Player = new(Game.Player)
	aldair.TeamPlace = Units.HomeTeam
	aldair.Number = "3"

	aldair.Coords = Physics.Point{
		PosX: strategy.RegionWidth * 3.4,//3
		PosY: strategy.RegionHeight * 2.1,//2
	}
	assert.Equal(t, strategy.RegionCode{3,2}, aldair.myCurrentRegion())

	aldair.Coords = Physics.Point{
		PosX: strategy.RegionWidth * 4.2,//4
		PosY: strategy.RegionHeight * 0.1,//0
	}
	assert.Equal(t, strategy.RegionCode{4,0}, aldair.myCurrentRegion())

	aldair.TeamPlace = Units.AwayTeam


	aldair.Coords = Physics.Point{
		PosX: strategy.RegionWidth * 4.2,//3
		PosY: strategy.RegionHeight * 0.1,//3
	}
	assert.Equal(t, strategy.RegionCode{3,3}, aldair.myCurrentRegion())

}

func TestBrain_isItInMyActiveRegion(t *testing.T) {
	aldair := Brain{}
	aldair.Player = new(Game.Player)
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