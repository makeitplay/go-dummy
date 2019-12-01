package bot

import (
	"github.com/lugobots/client-player-go/v2/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAngleWithRoute(t *testing.T) {
	type testCase struct {
		direction     proto.Vector
		from          proto.Point
		obstacle      proto.Point
		expectedAngle float64
	}

	testCases := map[string]testCase{
		"in front":   {direction: proto.North(), from: proto.Point{}, obstacle: proto.Point{Y: 100}, expectedAngle: 0},
		"behind":     {direction: proto.North(), from: proto.Point{Y: 100}, obstacle: proto.Point{}, expectedAngle: 180},
		"Right side": {direction: proto.North(), from: proto.Point{}, obstacle: proto.Point{X: 1}, expectedAngle: -90},
		"Left side":  {direction: proto.North(), from: proto.Point{X: 1}, obstacle: proto.Point{}, expectedAngle: 90},
	}

	for caseName, def := range testCases {
		assert.Equal(t, def.expectedAngle, AngleWithRoute(def.direction, def.from, def.obstacle), caseName)
	}
}
