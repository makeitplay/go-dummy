package brain

import (
	"testing"
	"github.com/makeitplay/go-dummy/strategy"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/stretchr/testify/assert"
)

func TestBrain_watchOpponentOnMyRoute(t *testing.T) {
	msg := Game.GameMessage{}
	msg.GameInfo = Game.GameInfo{}
	msg.GameInfo.Ball = Game.Ball{}
	msg.GameInfo.Ball.Coords = strategy.GetRegionCenter(strategy.RegionCode{0,0}, Units.HomeTeam)


	A := Brain{}
	A.Player = new(Game.Player)
	A.TeamPlace = Units.HomeTeam
	A.Id = 1
	A.Size = Units.PlayerSize

	B := Brain{}
	B.Player = new(Game.Player)
	A.Id = 2
	B.Size = Units.PlayerSize
	B.TeamPlace = Units.AwayTeam

	C := Brain{}
	C.Player = new(Game.Player)
	C.Size = Units.PlayerSize
	C.TeamPlace = Units.AwayTeam
	C.Id = 5

	D := Brain{}
	D.Player = new(Game.Player)
	D.TeamPlace = Units.AwayTeam
	D.Id = 4
	D.Size = Units.PlayerSize

	msg.GameInfo.HomeTeam.Players = map[int]*Game.Player{}
	msg.GameInfo.AwayTeam.Players = map[int]*Game.Player{}

	msg.GameInfo.HomeTeam.Players[0] = A.Player
	msg.GameInfo.AwayTeam.Players[0] = B.Player
	msg.GameInfo.AwayTeam.Players[1] = C.Player
	msg.GameInfo.AwayTeam.Players[2] = D.Player
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
