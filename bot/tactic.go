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
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	3: {
		Initial:       {1, 1},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	4: {
		Initial:       {1, 2},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	5: {
		Initial:       {1, 3},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	6: {
		Initial:       {2, 0},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	7: {
		Initial:       {2, 1},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	8: {
		Initial:       {2, 2},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	9: {
		Initial:       {2, 3},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	10: {
		Initial:       {3, 1},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	11: {
		Initial:       {3, 2},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
}
