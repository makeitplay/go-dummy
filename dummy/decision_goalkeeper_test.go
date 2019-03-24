package dummy

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNearestGoalPoint_InFrondOfTheGoal(t *testing.T) {
	goal := arena.HomeTeamGoal
	ball := client.Ball{}
	ball.Coords = goal.Center

	ball.Coords.PosX += rand.Int() % units.FieldWidth
	ball.Coords.PosY += units.PlayerMaxSpeed

	expectedPoint := goal.Center
	expectedPoint.PosY = ball.Coords.PosY

	actualDistance := NearestGoalPoint(ball, goal)

	assert.Equal(t, expectedPoint, actualDistance)
}

func TestNearestGoalPoint_AboveOnMap(t *testing.T) {
	goal := arena.HomeTeamGoal
	ball := client.Ball{}
	ball.Coords = goal.TopPole
	ball.Coords.PosY += units.BallSize

	expectedPoint := goal.TopPole
	expectedPoint.PosY -= units.BallSize / 2
	actualDistance := NearestGoalPoint(ball, goal)

	assert.Equal(t, expectedPoint, actualDistance)
}

func TestNearestGoalPoint_BelowOnMap(t *testing.T) {
	goal := arena.HomeTeamGoal
	ball := client.Ball{}
	ball.Coords = goal.BottomPole
	ball.Coords.PosY -= units.BallSize

	expectedPoint := goal.BottomPole

	expectedPoint.PosY += units.BallSize / 2
	actualDistance := NearestGoalPoint(ball, goal)

	assert.Equal(t, expectedPoint, actualDistance)
}

func TestFindThreatenedSpot_BallNotComing_Holder(t *testing.T) {
	goal := arena.HomeTeamGoal
	ball := client.Ball{}
	ball.Holder = &client.Player{}

	_, _, coming := findThreatenedSpot(ball, goal)
	assert.False(t, coming)
}

func TestFindThreatenedSpot_BallNotComing_Stopped(t *testing.T) {
	goal := arena.HomeTeamGoal

	ball := client.Ball{}
	ball.Velocity = physics.NewZeroedVelocity(physics.West)
	ball.Coords = goal.Center
	ball.Coords.PosX += Units.GoalZoneRange

	_, _, coming := findThreatenedSpot(ball, goal)
	assert.False(t, coming)
}

func TestFindThreatenedSpot_BallNotComing_DiffDirection(t *testing.T) {
	goal := arena.HomeTeamGoal

	ball := client.Ball{}
	ball.Velocity = physics.NewZeroedVelocity(physics.North)
	ball.Velocity.Speed = units.BallMaxSpeed
	ball.Coords = goal.Center
	ball.Coords.PosX += Units.GoalZoneRange

	_, _, coming := findThreatenedSpot(ball, goal)
	assert.False(t, coming)
}
func TestFindThreatenedSpot_BallNotComing_OppositeDirection(t *testing.T) {
	goal := arena.HomeTeamGoal

	ball := client.Ball{}
	ball.Velocity = physics.NewZeroedVelocity(physics.East)
	ball.Velocity.Speed = units.BallMaxSpeed
	ball.Coords = goal.Center
	ball.Coords.PosX += Units.GoalZoneRange

	_, _, coming := findThreatenedSpot(ball, goal)
	assert.False(t, coming)
}

func TestFindThreatenedSpot_BallNotComing_DoesNotHitGoal(t *testing.T) {
	goal := arena.HomeTeamGoal

	ball := client.Ball{}
	ball.Coords = goal.Center
	ball.Coords.PosX += Units.GoalZoneRange

	wrongTarget := goal.TopPole
	wrongTarget.PosY += units.BallSize * 2
	wrongShoot, _ := physics.NewVector(ball.Coords, wrongTarget)

	ball.Velocity = physics.NewZeroedVelocity(*wrongShoot)
	ball.Velocity.Speed = units.BallMaxSpeed
	_, _, coming := findThreatenedSpot(ball, goal)
	assert.False(t, coming)
}

func TestFindThreatenedSpot_BestPosition_TopPoleNoRush(t *testing.T) {

	serverConfig := new(client.Configuration)
	serverConfig.TeamPlace = arena.HomeTeam
	serverConfig.PlayerNumber = arena.GoalkeeperNumber
	serverConfig.WSHost = "localhost"
	serverConfig.WSPort = "8080"

	gamer := &client.Gamer{}

	PlayerNumber = serverConfig.PlayerNumber
	TeamPlace = serverConfig.TeamPlace
	MyRule = strategy.DefinePlayerRule(serverConfig.PlayerNumber)
	TeamBallPossession = TeamPlace
	ClientResponder = gamer

	gamer.OnAnnouncement = func(turnTx client.TurnContext) {
		logrus.Warn("cool")

	}
	gamer.Play(GetInitialRegion().Center(serverConfig.TeamPlace), serverConfig)

	ctrl := client.NewTestController()

	//goal := arena.HomeTeamGoal
	//
	//ball := client.Ball{}
	//ball.Coords = goal.TopPole
	//ball.Coords.PosY -= units.BallSize
	//ball.Coords.PosX += units.GoalZoneRange
	//ball.Velocity = physics.NewZeroedVelocity(physics.West)
	//ball.Velocity.Speed = units.BallMaxSpeed
	//
	//target, timeToReach, coming := findThreatenedSpot(ball, goal)
	//
	//expectedTarget := ball.Coords
	//expectedTarget.PosX = goal.Center.PosX
	//
	//optimumWatchingPosition(goal, target, timeToReach)
	//assert.True(t, coming)
	//assert.Equal(t, expectedTarget, target)
}
