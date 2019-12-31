package bot

import (
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShootingEvaluation_MustNotShooIfItIsTooFar(t *testing.T) {
	fieldCenter := field.FieldCenter()
	ball := &proto.Ball{Position: &fieldCenter}
	me := &proto.Player{TeamSide: proto.Team_HOME}
	snapshot := &proto.GameSnapshot{Ball: ball}
	assert.Equal(t, MustNot, ShootingEvaluation(me, snapshot))
}

func TestCountCloseOpponents(t *testing.T) {
	fieldCenter := field.FieldCenter()

	me := &proto.Player{Position: &fieldCenter}
	teamMate := &proto.Player{Position: &proto.Point{
		X: fieldCenter.X + 1000,
		Y: fieldCenter.Y,
	},
	}
	opponent := &proto.Player{Position: &proto.Point{
		X: fieldCenter.X + 500,
		Y: fieldCenter.Y,
	}}

	opponents := []*proto.Player{opponent}

	// case: between them
	assert.Equal(t, 1, isObstacleForPassing(me, *teamMate.Position, opponents))

	// case: at same point as me
	opponent.Position.X = fieldCenter.X
	assert.Equal(t, 1, isObstacleForPassing(me, *teamMate.Position, opponents))

	// case: at same point as the team mae
	opponent.Position.X = teamMate.Position.X
	opponent.Position.Y = teamMate.Position.Y
	assert.Equal(t, 1, isObstacleForPassing(me, *teamMate.Position, opponents))

	// case behind me
	opponent.Position.X = me.Position.X - (field.BallSize * 3)
	opponent.Position.Y = me.Position.Y
	assert.Equal(t, 0, isObstacleForPassing(me, *teamMate.Position, opponents))

	// case behind my team mate
	opponent.Position.X = teamMate.Position.X + 10
	opponent.Position.Y = teamMate.Position.Y
	assert.Equal(t, 0, isObstacleForPassing(me, *teamMate.Position, opponents))

	// case between us, in a little angle
	opponent.Position.X = fieldCenter.X + 500
	opponent.Position.Y = fieldCenter.Y + field.BallSize
	assert.Equal(t, 1, isObstacleForPassing(me, *teamMate.Position, opponents))

	// case between us, in a high angle
	opponent.Position.X = fieldCenter.X + 500
	opponent.Position.Y = fieldCenter.Y + field.BallSize*2
	assert.Equal(t, 0, isObstacleForPassing(me, *teamMate.Position, opponents))
}
