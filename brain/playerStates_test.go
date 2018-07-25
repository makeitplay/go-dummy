package brain

import (
	"testing"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
	"math"
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

func TestBrain_BestSpeedToTarget(t *testing.T) {
	player := Brain{}
	player.Player = new(Game.Player)
	player.LastMsg = Game.GameMessage{GameInfo: Game.GameInfo{Ball: Game.Ball{}}}

	player.LastMsg.GameInfo.Ball.Coords = Physics.Point{}

	assert.Equal(t, Units.BallMaxSpeed, player.BestSpeedToTarget(Physics.Point{int(PerfectPassDistance),0}))

	// V = Vo + at
	// 0 = Units.BallMaxSpeed - Units.BallDeceleration * t
	// t = Units.BallMaxSpeed / Units.BallDeceleration
	timeToZero := math.Ceil(Units.BallMaxSpeed / Units.BallDeceleration)
	// S = V*t + (at^2)/2
	distanceToZero := Units.BallMaxSpeed * timeToZero + (-Units.BallDeceleration*math.Pow(timeToZero,2))/2

	assert.Equal(t, Units.BallMaxSpeed, player.BestSpeedToTarget(Physics.Point{int(distanceToZero),0}))

	//imprecise distance
	expectedInitialSpeed := Units.BallMaxSpeed * 0.7
	expectedFrames := 1.0
	impreciseDistance := expectedInitialSpeed * expectedFrames + (-Units.BallDeceleration*math.Pow(expectedFrames,2))/2
	assert.Equal(t, expectedInitialSpeed, player.BestSpeedToTarget(Physics.Point{int(impreciseDistance),0}))

	expectedInitialSpeed = Units.BallMaxSpeed * 0.9
	expectedFrames = 3.0
	impreciseDistance = expectedInitialSpeed * expectedFrames + (-Units.BallDeceleration*math.Pow(expectedFrames,2))/2
	assert.Equal(t, expectedInitialSpeed, player.BestSpeedToTarget(Physics.Point{int(impreciseDistance),0}))

	expectedInitialSpeed = Units.BallMaxSpeed * 0.95
	expectedFrames = 5.0
	impreciseDistance = expectedInitialSpeed * expectedFrames + (-Units.BallDeceleration*math.Pow(expectedFrames,2))/2
	assert.Equal(t, expectedInitialSpeed, player.BestSpeedToTarget(Physics.Point{int(impreciseDistance),0}))
}