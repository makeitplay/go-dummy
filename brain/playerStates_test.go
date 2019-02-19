package brain

import (
	"github.com/makeitplay/arena/BasicTypes"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBrain_orderForAtckHelpFrg(t *testing.T) {
	player := new(Brain)
	player.Player = new(client.Player)
	player.Number = BasicTypes.PlayerNumber("8")
	player.Coords = strategy.RegionCode{4, 2}.Center(Units.HomeTeam)
	MyRule = strategy.MiddlePlayer

	holder := new(Brain)
	holder.Player = new(client.Player)
	holder.Number = BasicTypes.PlayerNumber("7")
	holder.Coords = strategy.RegionCode{4, 1}.Center(Units.HomeTeam)

	lastMsg := MountGameInfo()
	lastMsg.GameInfo.HomeTeam.Players = append(lastMsg.GameInfo.HomeTeam.Players, player.Player, holder.Player)

	lastMsg.GameInfo.Ball.Holder = holder.Player
	lastMsg.GameInfo.Ball.Coords = holder.Coords

	TeamState = player.DetermineMyTeamState(lastMsg)

	player.LastMsg = lastMsg
	holder.LastMsg = lastMsg

	msg, order := player.orderForSupporting()
	assert.Equal(t, "", msg)
	assert.Equal(t, string(strategy.Offensive), string(TeamState))
	assert.Len(t, order, 1)
	assert.Equal(t, order[0].Type, BasicTypes.MOVE)
}

func TestBrain_BestSpeedToTarget(t *testing.T) {
	player := Brain{}
	player.Player = new(client.Player)
	player.LastMsg = client.GameMessage{GameInfo: client.GameInfo{Ball: client.Ball{}}}

	player.LastMsg.GameInfo.Ball.Coords = physics.Point{}

	assert.Equal(t, Units.BallMaxSpeed, player.BestSpeedToTarget(physics.Point{int(PerfectPassDistance), 0}))

	// V = Vo + at
	// 0 = Units.BallMaxSpeed - Units.BallDeceleration * t
	// t = Units.BallMaxSpeed / Units.BallDeceleration
	timeToZero := math.Ceil(Units.BallMaxSpeed / Units.BallDeceleration)
	// S = V*t + (at^2)/2
	distanceToZero := Units.BallMaxSpeed*timeToZero + (-Units.BallDeceleration*math.Pow(timeToZero, 2))/2

	assert.Equal(t, Units.BallMaxSpeed, player.BestSpeedToTarget(physics.Point{int(distanceToZero), 0}))

	//imprecise distance
	expectedInitialSpeed := Units.BallMaxSpeed * 0.7
	expectedFrames := 1.0
	impreciseDistance := expectedInitialSpeed*expectedFrames + (-Units.BallDeceleration*math.Pow(expectedFrames, 2))/2
	assert.Equal(t, expectedInitialSpeed, player.BestSpeedToTarget(physics.Point{int(impreciseDistance), 0}))

	expectedInitialSpeed = Units.BallMaxSpeed * 0.9
	expectedFrames = 3.0
	impreciseDistance = expectedInitialSpeed*expectedFrames + (-Units.BallDeceleration*math.Pow(expectedFrames, 2))/2
	assert.Equal(t, expectedInitialSpeed, player.BestSpeedToTarget(physics.Point{int(impreciseDistance), 0}))

	expectedInitialSpeed = Units.BallMaxSpeed * 0.95
	expectedFrames = 5.0
	impreciseDistance = expectedInitialSpeed*expectedFrames + (-Units.BallDeceleration*math.Pow(expectedFrames, 2))/2
	assert.Equal(t, expectedInitialSpeed, player.BestSpeedToTarget(physics.Point{int(impreciseDistance), 0}))
}
