package dummy

import (
	"context"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"os"
	"testing"
)

var Controller client.Controller

func TestMain(m *testing.M) {
	ctx, stop := context.WithCancel(context.Background())

	defer stop()
	serverConfig := new(client.Configuration)
	serverConfig.WSPort = "8080"
	serverConfig.WSHost = "localhost"
	serverConfig.UUID = "local"
	integrationCtx, ctrl, err := client.NewTestController(ctx, *serverConfig)
	if err == nil {
		Controller = ctrl
		//Controller.SetFrameInterval(500 * time.Millisecond)
		go func() {
			select {
			case <-integrationCtx.Done():
				log.Fatal("integration test was interrupted by the controller")
			}
		}()
	}
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func resetIntagrationState() {
	if Controller != nil {
		Controller.ResetScore()
	}
}

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

func TestIntegration_KeeperShouldCatch_atCenterFromCenter(t *testing.T) {
	if Controller == nil {
		t.Skip("no integration available")
	}

	ballInitialPosition := arena.HomeTeamGoal.Center
	ballInitialPosition.PosX += Units.GoalZoneRange + 1 //avoiding the ball be kicked out after "BallTimeInGoalZone" frames
	ballVelocity := physics.NewZeroedVelocity(physics.West)
	ballVelocity.Speed = units.BallMaxSpeed

	if _, err := Controller.SetPlayerPos(arena.HomeTeam, arena.GoalkeeperNumber, arena.HomeTeamGoal.Center); err != nil {
		t.Fatal(err)
	}
	gameStat, err := Controller.SetBallProperties(ballVelocity, ballInitialPosition)
	if err != nil {
		t.Fatal(err)
	}
	ballState := gameStat.Ball()

	for ballState.Velocity.Speed > 0 {
		turnCtx, err := Controller.GetGamerCtx(arena.HomeTeam, arena.GoalkeeperNumber)
		if err != nil {
			t.Fatal(err)
		}
		player := &Dummy{
			GameMsg: turnCtx.GameMsg(),
			Player:  turnCtx.Player(),
		}

		_, orderList := player.orderForGoalkeeper()
		Controller.SendOrders(arena.HomeTeam, arena.GoalkeeperNumber, orderList)
		newState, _ := Controller.NextTurn()
		ballState = newState.Ball()
	}
	if ballState.Holder == nil {
		t.Fatal("should had caught the ball. Middle")
	}
	assert.Equal(t, ballState.Holder.TeamPlace, arena.HomeTeam)
	assert.Equal(t, ballState.Holder.Number, arena.GoalkeeperNumber)
}

func TestIntegration_KeeperShouldCatch_atTopFromCenter(t *testing.T) {

	ballInitialPosition := arena.HomeTeamGoal.Center
	ballInitialPosition.PosX += Units.GoalZoneRange + Units.BallSize //avoiding the ball be kicked out after "BallTimeInGoalZone" frames
	ballVelocity := physics.NewZeroedVelocity(physics.West)
	ballVelocity.Speed = units.BallMaxSpeed

	shootTarget := arena.HomeTeamGoal.TopPole
	shootTarget.PosY -= 1
	shootTop, _ := physics.NewVector(ballInitialPosition, shootTarget)
	ballVelocity = physics.NewZeroedVelocity(*shootTop)
	ballVelocity.Speed = units.BallMaxSpeed

	gameStat, err := Controller.SetBallProperties(ballVelocity, ballInitialPosition)
	if err != nil {
		t.Fatal(err)
	}
	ballState := gameStat.Ball()

	for ballState.Velocity.Speed > 0 {
		turnCtx, err := Controller.GetGamerCtx(arena.HomeTeam, arena.GoalkeeperNumber)
		if err != nil {
			t.Fatal(err)
		}
		player := &Dummy{
			GameMsg: turnCtx.GameMsg(),
			Player:  turnCtx.Player(),
		}

		_, orderList := player.orderForGoalkeeper()
		Controller.SendOrders(arena.HomeTeam, arena.GoalkeeperNumber, orderList)
		newState, _ := Controller.NextTurn()
		ballState = newState.Ball()
	}
	if ballState.Holder == nil {
		t.Fatal("should had caught the ball. Top")
	}

	assert.Equal(t, ballState.Holder.TeamPlace, arena.HomeTeam)
	assert.Equal(t, ballState.Holder.Number, arena.GoalkeeperNumber)
}

//ctrl.ResetScore()
//ctrl.SendOrders(arena.HomeTeam, arena.GoalkeeperNumber, []orders.Order{
//
//})
//
//for gameStat.Ball().Velocity.Speed > 0 {
//	ctrl.ResetScore()
//	gameStat, err = ctrl.NextTurn()
//	if err != nil {
//		t.Fatal(err)
//	}
//	v := gameStat.Ball().Velocity
//	logrus.Infof("Ball velocity %v", &v)
//	ctrl.SetGameTurn( 100)
//}
//}
