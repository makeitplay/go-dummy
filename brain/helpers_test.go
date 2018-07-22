package brain

import (
	"testing"
	"github.com/makeitplay/go-dummy/strategy"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/stretchr/testify/assert"
	"github.com/makeitplay/commons/BasicTypes"
)

func TestBrain_watchOpponentOnMyRoute(t *testing.T) {
	msg := Game.GameMessage{}
	msg.GameInfo = Game.GameInfo{}
	msg.GameInfo.Ball = Game.Ball{}
	msg.GameInfo.Ball.Coords = strategy.GetRegionCenter(strategy.RegionCode{0,0}, Units.HomeTeam)


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

	A.Coords = strategy.GetRegionCenter(strategy.RegionCode{1,1}, Units.HomeTeam)
	B.Coords = strategy.GetRegionCenter(strategy.RegionCode{2,1}, Units.HomeTeam)
	C.Coords = strategy.GetRegionCenter(strategy.RegionCode{3,1}, Units.HomeTeam)
	D.Coords = strategy.GetRegionCenter(strategy.RegionCode{4,1}, Units.HomeTeam)


	target := strategy.GetRegionCenter(strategy.RegionCode{5,1}, Units.HomeTeam)
	objstacles := watchOpponentOnMyRoute(A.Player, target)
	assert.Len(t, objstacles, 3)

	D.Coords = strategy.GetRegionCenter(strategy.RegionCode{4,2}, Units.HomeTeam)
	objstacles = watchOpponentOnMyRoute(A.Player, target)
	assert.Len(t, objstacles, 2)

	C.Coords = strategy.GetRegionCenter(strategy.RegionCode{3,2}, Units.HomeTeam)
	objstacles = watchOpponentOnMyRoute(A.Player, target)
	assert.Len(t, objstacles, 1)

	B.Coords = strategy.GetRegionCenter(strategy.RegionCode{2,2}, Units.HomeTeam)
	objstacles = watchOpponentOnMyRoute(A.Player, target)
	assert.Len(t, objstacles, 0)
}
