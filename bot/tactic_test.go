package bot

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		_, ok := regMap[Initial]
		assert.True(t, ok, fmt.Sprintf("missing %s map for player %d", Initial, i))
		_, ok = regMap[UnderPressure]
		assert.True(t, ok, fmt.Sprintf("missing %s map for player %d", UnderPressure, i))
		_, ok = regMap[Defensive]
		assert.True(t, ok, fmt.Sprintf("missing %s map for player %d", Defensive, i))
		_, ok = regMap[Neutral]
		assert.True(t, ok, fmt.Sprintf("missing %s map for player %d", Neutral, i))
		_, ok = regMap[Offensive]
		assert.True(t, ok, fmt.Sprintf("missing %s map for player %d", Offensive, i))
		_, ok = regMap[OnAttack]
		assert.True(t, ok, fmt.Sprintf("missing %s map for player %d", OnAttack, i))
	}
}
