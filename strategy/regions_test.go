package strategy

import (
	"github.com/lugobots/arena"
	"github.com/lugobots/arena/physics"
	"github.com/lugobots/arena/units"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRegionCode(t *testing.T) {
	code := GetRegionCode(physics.Point{0, 0}, arena.HomeTeam)
	assert.Equal(t, 0, code.X)
	assert.Equal(t, 0, code.Y)

	code = GetRegionCode(physics.Point{0, 0}, arena.AwayTeam)
	assert.Equal(t, 7, code.X)
	assert.Equal(t, 3, code.Y)

	code = GetRegionCode(physics.Point{units.FieldWidth, units.FieldHeight}, arena.AwayTeam)
	assert.Equal(t, 0, code.X)
	assert.Equal(t, 0, code.Y)

	code = GetRegionCode(physics.Point{units.FieldWidth, units.FieldHeight}, arena.HomeTeam)
	assert.Equal(t, 7, code.X)
	assert.Equal(t, 3, code.Y)

	code = GetRegionCode(physics.Point{0, units.FieldHeight}, arena.HomeTeam)
	assert.Equal(t, 0, code.X)
	assert.Equal(t, 3, code.Y)

	code = GetRegionCode(physics.Point{0, units.FieldHeight}, arena.AwayTeam)
	assert.Equal(t, 7, code.X)
	assert.Equal(t, 0, code.Y)
}

func TestGetRegionCenter(t *testing.T) {
	halfRegionHeight := RegionHeight / 2
	halfRegionWidth := RegionWidth / 2

	center := RegionCode{0, 0}.Center(arena.HomeTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, halfRegionHeight, center.PosY)

	center = RegionCode{0, 0}.Center(arena.AwayTeam)
	assert.Equal(t, units.FieldWidth-halfRegionWidth, center.PosX)
	assert.Equal(t, units.FieldHeight-halfRegionHeight, center.PosY)

	center = RegionCode{7, 3}.Center(arena.HomeTeam)
	assert.Equal(t, units.FieldWidth-halfRegionWidth, center.PosX)
	assert.Equal(t, units.FieldHeight-halfRegionHeight, center.PosY)

	center = RegionCode{7, 3}.Center(arena.AwayTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, halfRegionHeight, center.PosY)

	center = RegionCode{0, 3}.Center(arena.HomeTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, units.FieldHeight-halfRegionHeight, center.PosY)

	center = RegionCode{7, 0}.Center(arena.AwayTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, units.FieldHeight-halfRegionHeight, center.PosY)

}
