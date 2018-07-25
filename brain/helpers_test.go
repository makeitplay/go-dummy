package brain

import (
	"testing"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/stretchr/testify/assert"
	"github.com/makeitplay/commons/BasicTypes"
)

func MountGameInfo() Game.GameMessage {
	fakeMsg := Game.GameMessage{}
	fakeMsg.GameInfo.Ball = Game.Ball{}
	fakeMsg.GameInfo.Turn = 1
	fakeMsg.GameInfo.AwayTeam = Game.Team{}
	fakeMsg.GameInfo.HomeTeam  = Game.Team{}
	return fakeMsg

}

func TestBrain_watchOpponentOnMyRoute(t *testing.T) {
	msg := Game.GameMessage{}
	msg.GameInfo = Game.GameInfo{}
	msg.GameInfo.Ball = Game.Ball{}
	msg.GameInfo.Ball.Coords = strategy.RegionCode{0, 0}.Center(Units.HomeTeam)

	A := Brain{}
	A.Player = new(Game.Player)
	A.Number = BasicTypes.PlayerNumber("1")
	A.TeamPlace = Units.HomeTeam
	A.Size = Units.PlayerSize

	B := Brain{}
	B.Player = new(Game.Player)
	B.Number = BasicTypes.PlayerNumber("1")
	B.TeamPlace = Units.AwayTeam
	B.Size = Units.PlayerSize

	C := Brain{}
	C.Player = new(Game.Player)
	C.Number = BasicTypes.PlayerNumber("2")
	C.TeamPlace = Units.AwayTeam
	C.Size = Units.PlayerSize

	D := Brain{}
	D.Player = new(Game.Player)
	D.Number = BasicTypes.PlayerNumber("3")
	D.TeamPlace = Units.AwayTeam
	D.Size = Units.PlayerSize

	msg.GameInfo.HomeTeam.Players = []*Game.Player{}
	msg.GameInfo.AwayTeam.Players = []*Game.Player{}

	msg.GameInfo.HomeTeam.Players = append(msg.GameInfo.HomeTeam.Players, A.Player)
	msg.GameInfo.AwayTeam.Players = append(msg.GameInfo.AwayTeam.Players, B.Player)
	msg.GameInfo.AwayTeam.Players = append(msg.GameInfo.AwayTeam.Players, C.Player)
	msg.GameInfo.AwayTeam.Players = append(msg.GameInfo.AwayTeam.Players, D.Player)
	A.LastMsg = msg

	A.Coords = strategy.RegionCode{1, 1}.Center(Units.HomeTeam)
	B.Coords = strategy.RegionCode{2, 1}.Center(Units.HomeTeam)
	C.Coords = strategy.RegionCode{3, 1}.Center(Units.HomeTeam)
	D.Coords = strategy.RegionCode{4, 1}.Center(Units.HomeTeam)

	target := strategy.RegionCode{5, 1}.Center(Units.HomeTeam)
	obstacles := watchOpponentOnMyRoute(msg.GameInfo, A.Player, target)
	assert.Len(t, obstacles, 3)
	assert.Equal(t, float64(strategy.RegionWidth-Units.PlayerSize), A.Player.Coords.DistanceTo(obstacles[0]))
	assert.Equal(t, float64(2*strategy.RegionWidth-Units.PlayerSize), A.Player.Coords.DistanceTo(obstacles[1]))
	assert.Equal(t, float64(3*strategy.RegionWidth-Units.PlayerSize), A.Player.Coords.DistanceTo(obstacles[2]))

	D.Coords = strategy.RegionCode{4, 2}.Center(Units.HomeTeam)
	obstacles = watchOpponentOnMyRoute(msg.GameInfo, A.Player, target)
	assert.Len(t, obstacles, 2)

	C.Coords = strategy.RegionCode{3, 2}.Center(Units.HomeTeam)
	obstacles = watchOpponentOnMyRoute(msg.GameInfo, A.Player, target)
	assert.Len(t, obstacles, 1)

	B.Coords = strategy.RegionCode{2, 2}.Center(Units.HomeTeam)
	obstacles = watchOpponentOnMyRoute(msg.GameInfo, A.Player, target)
	assert.Len(t, obstacles, 0)
}

