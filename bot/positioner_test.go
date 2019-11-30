package bot

import (
	"github.com/lugobots/client-player-go/v2/lugo"
	"github.com/lugobots/client-player-go/v2/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPositioner(t *testing.T) {
	p, err := NewPositioner(MinCols, MinRows, proto.Team_HOME)
	assert.Nil(t, err)

	myStruct, ok := p.(*positioner)
	assert.True(t, ok)
	assert.Equal(t, lugo.FieldWidth/int(MinCols), int(myStruct.regionWidth))
	assert.Equal(t, lugo.FieldHeight/int(MinRows), int(myStruct.regionHeight))
}

func TestNewPositioner_InvalidArgs(t *testing.T) {
	p, err := NewPositioner(MinCols-1, MinRows, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMinCols, err)

	p, err = NewPositioner(MaxCols+1, MinRows, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMaxCols, err)

	p, err = NewPositioner(MinCols, MinRows-1, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMinRows, err)

	p, err = NewPositioner(MinCols, MaxCols+1, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMaxRows, err)
}

func TestRegion_Center_HomeTeam(t *testing.T) {
	type testCase struct {
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := proto.Team_HOME

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		expectedPointDefenseRight := proto.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		expectedPointDefenseLeft := proto.Point{X: +testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}
		expectedPointAttackLeft := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}
		expectedPointAttackRight := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}

		r, err := p.GetRegion(0, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointDefenseRight, r.Center(), testName)

		r, err = p.GetRegion(0, testSettings.rows-1)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointDefenseLeft, r.Center(), testName)

		r, err = p.GetRegion(testSettings.cols-1, testSettings.rows-1)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointAttackLeft, r.Center(), testName)

		r, err = p.GetRegion(testSettings.cols-1, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointAttackRight, r.Center(), testName)
	}
}

func TestRegion_Center_Away(t *testing.T) {
	type testCase struct {
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := proto.Team_AWAY

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		expectedPointDefenseRight := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}
		expectedPointDefenseLeft := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		expectedPointAttackLeft := proto.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		expectedPointAttackRight := proto.Point{X: testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}

		r, err := p.GetRegion(0, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointDefenseRight, r.Center(), testName)

		r, err = p.GetRegion(0, testSettings.rows-1)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointDefenseLeft, r.Center(), testName)

		r, err = p.GetRegion(testSettings.cols-1, testSettings.rows-1)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointAttackLeft, r.Center(), testName)

		r, err = p.GetRegion(testSettings.cols-1, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedPointAttackRight, r.Center(), testName)
	}
}

func TestPositioner_GetRegion_InvalidArgs(t *testing.T) {
	p, err := NewPositioner(10, 10, proto.Team_AWAY)
	assert.Nil(t, err)

	r, err := p.GetRegion(11, 5)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxCols, err)

	r, err = p.GetRegion(10, 5)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxCols, err)

	r, err = p.GetRegion(9, 11)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxRows, err)

	r, err = p.GetRegion(9, 10)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxRows, err)

}

func TestPositioner_GetPointRegion_HomeTeam(t *testing.T) {
	type testCase struct {
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := proto.Team_HOME

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		pointDefenseRight := proto.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		pointDefenseLeft := proto.Point{X: +testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}
		pointAttackLeft := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}
		pointAttackRight := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}

		r, err := p.GetPointRegion(pointDefenseRight)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)

		r, err = p.GetPointRegion(pointDefenseLeft)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackLeft)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackRight)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)
	}
}

func TestPositioner_GetPointRegion_AwayTeam(t *testing.T) {
	type testCase struct {
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := proto.Team_AWAY

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		pointDefenseRight := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}
		pointDefenseLeft := proto.Point{X: lugo.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		pointAttackLeft := proto.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		pointAttackRight := proto.Point{X: testSettings.regionHalfWidth, Y: lugo.FieldHeight - testSettings.regionHalfHeight}

		r, err := p.GetPointRegion(pointDefenseRight)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)

		r, err = p.GetPointRegion(pointDefenseLeft)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackLeft)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackRight)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)
	}
}
