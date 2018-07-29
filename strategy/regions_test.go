package strategy

import (
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRegionCode(t *testing.T) {
	code := GetRegionCode(Physics.Point{0, 0}, Units.HomeTeam)
	assert.Equal(t, 0, code.X)
	assert.Equal(t, 0, code.Y)

	code = GetRegionCode(Physics.Point{0, 0}, Units.AwayTeam)
	assert.Equal(t, 7, code.X)
	assert.Equal(t, 3, code.Y)

	code = GetRegionCode(Physics.Point{Units.CourtWidth, Units.CourtHeight}, Units.AwayTeam)
	assert.Equal(t, 0, code.X)
	assert.Equal(t, 0, code.Y)

	code = GetRegionCode(Physics.Point{Units.CourtWidth, Units.CourtHeight}, Units.HomeTeam)
	assert.Equal(t, 7, code.X)
	assert.Equal(t, 3, code.Y)

	code = GetRegionCode(Physics.Point{0, Units.CourtHeight}, Units.HomeTeam)
	assert.Equal(t, 0, code.X)
	assert.Equal(t, 3, code.Y)

	code = GetRegionCode(Physics.Point{0, Units.CourtHeight}, Units.AwayTeam)
	assert.Equal(t, 7, code.X)
	assert.Equal(t, 0, code.Y)
}

func TestGetRegionCenter(t *testing.T) {
	halfRegionHeight := RegionHeight / 2
	halfRegionWidth := RegionWidth / 2
	center := RegionCode{0, 0}.Center(Units.HomeTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, halfRegionHeight, center.PosY)

	center = RegionCode{0, 0}.Center(Units.AwayTeam)
	assert.Equal(t, Units.CourtWidth-halfRegionWidth, center.PosX)
	assert.Equal(t, Units.CourtHeight-halfRegionHeight, center.PosY)

	center = RegionCode{7, 3}.Center(Units.HomeTeam)
	assert.Equal(t, Units.CourtWidth-halfRegionWidth, center.PosX)
	assert.Equal(t, Units.CourtHeight-halfRegionHeight, center.PosY)

	center = RegionCode{7, 3}.Center(Units.AwayTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, halfRegionHeight, center.PosY)

	center = RegionCode{0, 3}.Center(Units.HomeTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, Units.CourtHeight-halfRegionHeight, center.PosY)

	center = RegionCode{7, 0}.Center(Units.AwayTeam)
	assert.Equal(t, halfRegionWidth, center.PosX)
	assert.Equal(t, Units.CourtHeight-halfRegionHeight, center.PosY)

}
