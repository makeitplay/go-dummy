package dummy

import (
	"context"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBrain_GetActiveRegion(t *testing.T) {
	aldair := Dummy{}
	aldair.Player = new(client.Player)
	PlayerNumber = "3"

	aldair.TeamState = strategy.UnderPressure
	assert.Equal(t, strategy.RegionCode{0, 1}, aldair.GetActiveRegion())

	aldair.TeamState = strategy.Defensive
	assert.Equal(t, strategy.RegionCode{1, 1}, aldair.GetActiveRegion())

	PlayerNumber = "8"

	aldair.TeamState = strategy.UnderPressure
	assert.Equal(t, strategy.RegionCode{1, 2}, aldair.GetActiveRegion())

	aldair.TeamState = strategy.OnAttack
	assert.Equal(t, strategy.RegionCode{5, 2}, aldair.GetActiveRegion())
}

func TestBrain_myCurrentRegion(t *testing.T) {
	aldair := Dummy{}
	aldair.Player = new(client.Player)
	TeamPlace = arena.HomeTeam
	PlayerNumber = "3"

	aldair.Player.Coords = physics.Point{
		PosX: strategy.RegionWidth * 3.4,  //3
		PosY: strategy.RegionHeight * 2.1, //2
	}
	assert.Equal(t, strategy.RegionCode{3, 2}, aldair.myCurrentRegion())

	aldair.Player.Coords = physics.Point{
		PosX: strategy.RegionWidth * 4.2,  //4
		PosY: strategy.RegionHeight * 0.1, //0
	}
	assert.Equal(t, strategy.RegionCode{4, 0}, aldair.myCurrentRegion())

	TeamPlace = arena.AwayTeam

	aldair.Player.Coords = physics.Point{
		PosX: strategy.RegionWidth * 4.2,  //3
		PosY: strategy.RegionHeight * 0.1, //3
	}
	assert.Equal(t, strategy.RegionCode{3, 3}, aldair.myCurrentRegion())

}

func TestBrain_isItInMyActiveRegion(t *testing.T) {
	aldair := Dummy{}
	aldair.Player = new(client.Player)
	aldair.TeamState = strategy.OnAttack
	PlayerNumber = "5"

	assert.True(t, aldair.isItInMyActiveRegion(physics.Point{
		PosX: strategy.RegionWidth * 4.2,
		PosY: strategy.RegionHeight * 3.1,
	}, strategy.OnAttack))

	assert.False(t, aldair.isItInMyActiveRegion(physics.Point{
		PosX: strategy.RegionWidth * 3.2,
		PosY: strategy.RegionHeight * 3.1,
	}, strategy.OnAttack))

	assert.False(t, aldair.isItInMyActiveRegion(physics.Point{
		PosX: strategy.RegionWidth * 4.2,
		PosY: strategy.RegionHeight * 2.1,
	}, strategy.OnAttack))

}

func TestDetermineMyTeamState_NoBall(t *testing.T) {
	msg := client.GameMessage{}
	msg.GameInfo = client.GameInfo{}
	msg.GameInfo.Ball = client.Ball{}

	homePlayer := new(Dummy)
	homePlayer.Player = new(client.Player)
	msg.GameInfo.HomeTeam.Players = []*client.Player{homePlayer.Player}

	TeamPlace = arena.HomeTeam

	//awayPlayer := new(Brain)
	//awayPlayer.Player = new(client.Player)
	//awayPlayer.TeamPlace = arena.AwayTeam

	TeamBallPossession = arena.AwayTeam
	msg.GameInfo.Ball.Coords = strategy.RegionCode{0, 2}.Center(arena.HomeTeam)
	ctx, _ := client.NewGamerContext(context.Background(), &client.Configuration{TeamPlace: arena.HomeTeam, PlayerNumber: "5"})
	assert.Equal(t, strategy.UnderPressure, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.OnAttack, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{1, 1}.Center(arena.HomeTeam)
	assert.Equal(t, strategy.UnderPressure, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.OnAttack, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{2, 1}.Center(arena.HomeTeam)
	assert.Equal(t, strategy.UnderPressure, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Offensive, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{3, 1}.Center(arena.HomeTeam)
	assert.Equal(t, strategy.Defensive, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Offensive, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{4, 1}.Center(arena.HomeTeam)
	assert.Equal(t, strategy.Defensive, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{5, 1}.Center(arena.HomeTeam)
	assert.Equal(t, strategy.Neutral, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{6, 1}.Center(arena.HomeTeam)
	assert.Equal(t, strategy.Neutral, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{7, 1}.Center(arena.HomeTeam)
	assert.Equal(t, strategy.Offensive, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	TeamBallPossession = arena.HomeTeam
	msg.GameInfo.Ball.Coords = strategy.RegionCode{0, 2}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.OnAttack, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.UnderPressure, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{1, 1}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.OnAttack, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.UnderPressure, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{2, 1}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.Offensive, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.UnderPressure, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{3, 1}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.Offensive, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{4, 1}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.Neutral, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Defensive, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{5, 1}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.Neutral, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{6, 1}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.Defensive, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Neutral, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

	msg.GameInfo.Ball.Coords = strategy.RegionCode{7, 1}.Center(arena.AwayTeam)
	assert.Equal(t, strategy.Defensive, strategy.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))
	//assert.Equal(t, strategy.Offensive, awayPlayer.DetermineMyTeamState(ctx.CreateTurnContext(msg), arena.HomeTeam))

}
