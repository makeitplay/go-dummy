package bot

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

var knownTeamStates = []TeamState{
	Initial,
	UnderPressure,
	Defensive,
	Neutral,
	Offensive,
	OnAttack,
}

// Note that this method only tests if the tactic settings are well defined. Does not test if the settings are defined
// as you expect. Please create your own test to check your logic.
func TestDefineRole(t *testing.T) {
	for i := uint32(2); i <= 11; i++ {
		assert.NotEqual(t, "", DefineRole(i))
	}
}

// Note that this method only tests if the tactic settings are well defined. Does not test if the settings are defined
// as you expect. Please create your own test to check your logic.
func TestDetermineTeamState_ShouldHaveStateForAllCols(t *testing.T) {
	p, err := coach.NewPositioner(RegionCols, RegionRows, proto.Team_HOME)
	if err != nil {
		t.Fatalf("invalid settings, cannot created a NewPositioner: %s", err)
	}
	for i := uint8(0); i < RegionCols; i++ {

		ballRegion, err := p.GetRegion(i, RegionRows/2)
		if err != nil {
			t.Fatalf("could not test the team state when the ball is in col %d, row %d: %s", i, RegionRows/2, err)
		}

		_, err = DetermineTeamState(ballRegion, proto.Team_HOME, proto.Team_AWAY)
		if err != nil {
			t.Fatalf("could not defined team state when the ball is in region %s: %s", ballRegion, err)
		}
	}
}

// Note that this method only tests if the tactic settings are well defined. Does not test if the settings are defined
// as you expect. Please create your own test to check your logic.
func TestDefineRegionMap_ShouldMapAllPlayersInAllTeamStates(t *testing.T) {
	// starting from 2 because the number goalkeeper has RegionMap
	for i := uint32(2); i <= 11; i++ {
		regMap := DefineRegionMap(i)
		assert.NotNil(t, regMap, fmt.Sprintf("missing maps for player %d", i))
		for _, state := range knownTeamStates {
			_, ok := regMap[state]
			assert.True(t, ok, fmt.Sprintf("missing %s map for player %d", state, i))
		}
	}
}

func TestDefineRegionMap_AllMappedRegionShouldBeCompatibleWithOurPositioner(t *testing.T) {
	p, err := coach.NewPositioner(RegionCols, RegionRows, proto.Team_HOME)
	if err != nil {
		t.Fatalf("invalid settings, cannot created a NewPositioner: %s", err)
	}

	// starting from 2 because the number goalkeeper has RegionMap
	for i := uint32(2); i <= 11; i++ {
		regMap := DefineRegionMap(i)
		assert.NotNil(t, regMap, fmt.Sprintf("missing maps for player %d", i))
		for _, state := range knownTeamStates {
			r := regMap[state]
			_, err = p.GetRegion(r.Col, r.Row)
			assert.Nil(t, err, fmt.Sprintf("invalid region mapped for state %s player %d: %v", state, i, r))
		}
	}
}
