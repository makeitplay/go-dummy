package bot

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/proto"
)

func DefineRole(number uint32) Role {
	// starting from 2 because the number goalkeeper has no role
	switch number {
	case 2, 3, 4, 5:
		return Defense
	case 6, 7, 8, 9:
		return Middle
	case 10, 11:
		return Attack
	}
	return ""
}

func DetermineTeamState(ballRegion coach.Region, teamSide, possession proto.Team_Side) (s TeamState, e error) {
	regionCol := ballRegion.Col()
	if possession == teamSide {
		switch regionCol {
		case 6, 7:
			return OnAttack, nil
		case 4, 5:
			return Offensive, nil
		case 3:
			return Neutral, nil
		case 2:
			return Defensive, nil
		case 0, 1:
			return UnderPressure, nil
		}

	} else {
		switch regionCol {
		case 6, 7:
			return OnAttack, nil
		case 4, 5:
			return Offensive, nil
		case 3:
			return Neutral, nil
		case 2:
			return Defensive, nil
		case 0, 1:
			return UnderPressure, nil
		}
	}
	return "", fmt.Errorf("unknown team state for ball in %d col, tor possion with %s", regionCol, possession)
}

var roleMap = map[uint32]RegionMap{
	// starting from 2 because the number goalkeeper has RegionMap
	2: {
		Initial:       {1, 0},
		UnderPressure: {0, 1},
		Defensive:     {1, 0},
		Neutral:       {2, 0},
		Offensive:     {1, 0},
		OnAttack:      {1, 1},
	},
	3: {
		Initial:       {1, 1},
		UnderPressure: {1, 1},
		Defensive:     {1, 1},
		Neutral:       {1, 1},
		Offensive:     {3, 1},
		OnAttack:      {4, 1},
	},
	4: {
		Initial:       {1, 2},
		UnderPressure: {1, 2},
		Defensive:     {1, 2},
		Neutral:       {1, 2},
		Offensive:     {3, 2},
		OnAttack:      {4, 2},
	},
	5: {
		Initial:       {1, 3},
		UnderPressure: {0, 2},
		Defensive:     {1, 3},
		Neutral:       {2, 3},
		Offensive:     {1, 3},
		OnAttack:      {1, 2},
	},
	6: {
		Initial:       {2, 0},
		UnderPressure: {1, 0},
		Defensive:     {2, 1},
		Neutral:       {3, 1},
		Offensive:     {4, 1},
		OnAttack:      {5, 0},
	},
	7: {
		Initial:       {2, 1},
		UnderPressure: {2, 1},
		Defensive:     {3, 1},
		Neutral:       {4, 1},
		Offensive:     {5, 1},
		OnAttack:      {6, 1},
	},
	8: {
		Initial:       {2, 2},
		UnderPressure: {2, 2},
		Defensive:     {3, 2},
		Neutral:       {4, 2},
		Offensive:     {5, 2},
		OnAttack:      {6, 2},
	},
	9: {
		Initial:       {2, 3},
		UnderPressure: {1, 3},
		Defensive:     {2, 2},
		Neutral:       {3, 2},
		Offensive:     {4, 2},
		OnAttack:      {5, 3},
	},
	10: {
		Initial:       {3, 1},
		UnderPressure: {3, 1},
		Defensive:     {4, 1},
		Neutral:       {5, 1},
		Offensive:     {6, 1},
		OnAttack:      {7, 1},
	},
	11: {
		Initial:       {3, 2},
		UnderPressure: {3, 2},
		Defensive:     {4, 2},
		Neutral:       {6, 2},
		Offensive:     {6, 2},
		OnAttack:      {7, 2},
	},
}
