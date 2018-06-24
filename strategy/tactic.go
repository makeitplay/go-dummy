package strategy

import (
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/client-player-go/Game"
	"github.com/makeitplay/commons/Units"
)

type TeamState string

const UnderPressure TeamState = "under-pressure"
const Defensive TeamState = "defensive"
const Neutral TeamState = "neutral"
const Offensive TeamState = "offensive"
const OnAttack TeamState = "on-attack"

func DetermineTeamState(msg Game.GameMessage, place Units.TeamPlace) TeamState {
	ballRegionCode := GetRegionCode(msg.GameInfo.Ball.Coords, place)
	if msg.GameInfo.Ball.Holder == nil || msg.GameInfo.Ball.Holder.TeamPlace != place {
		if ballRegionCode.X < 2 {
			return UnderPressure
		} else if ballRegionCode.X < 4 {
			return Defensive
		} else if ballRegionCode.X < 6 {
			return Neutral
		} else {
			return Offensive
		}
	} else {
		if ballRegionCode.X < 2 {
			return Defensive
		} else if ballRegionCode.X < 4 {
			return Neutral
		} else if ballRegionCode.X < 6 {
			return Offensive
		} else {
			return OnAttack
		}
	}
}

type Region map[TeamState]RegionCode

var PlayerRegionMap = map[BasicTypes.PlayerNumber]Region{
	//defense players
	"2": {
		UnderPressure: RegionCode{0, 0},
		Defensive:     RegionCode{1, 0},
		Neutral:       RegionCode{2, 0},
		Offensive:     RegionCode{3, 0},
		OnAttack:      RegionCode{4, 0},
	},
	"3": {
		UnderPressure: RegionCode{0, 1},
		Defensive:     RegionCode{1, 1},
		Neutral:       RegionCode{2, 1},
		Offensive:     RegionCode{3, 1},
		OnAttack:      RegionCode{4, 1},
	},
	"4": {
		UnderPressure: RegionCode{0, 2},
		Defensive:     RegionCode{1, 2},
		Neutral:       RegionCode{2, 2},
		Offensive:     RegionCode{3, 2},
		OnAttack:      RegionCode{4, 2},
	},
	"5": {
		UnderPressure: RegionCode{0, 3},
		Defensive:     RegionCode{1, 3},
		Neutral:       RegionCode{2, 3},
		Offensive:     RegionCode{3, 3},
		OnAttack:      RegionCode{4, 3},
	},

	//middle players
	"6": {
		UnderPressure: RegionCode{1, 0},
		Defensive:     RegionCode{2, 0},
		Neutral:       RegionCode{3, 0},
		Offensive:     RegionCode{4, 0},
		OnAttack:      RegionCode{5, 0},
	},
	"7": {
		UnderPressure: RegionCode{1, 1},
		Defensive:     RegionCode{2, 1},
		Neutral:       RegionCode{3, 1},
		Offensive:     RegionCode{4, 1},
		OnAttack:      RegionCode{5, 1},
	},
	"8": {
		UnderPressure: RegionCode{1, 2},
		Defensive:     RegionCode{2, 2},
		Neutral:       RegionCode{3, 2},
		Offensive:     RegionCode{4, 2},
		OnAttack:      RegionCode{5, 2},
	},
	"9": {
		UnderPressure: RegionCode{1, 3},
		Defensive:     RegionCode{2, 3},
		Neutral:       RegionCode{3, 3},
		Offensive:     RegionCode{4, 3},
		OnAttack:      RegionCode{5, 3},
	},

	//attack players
	"10": {
		UnderPressure: RegionCode{2, 1},
		Defensive:     RegionCode{3, 1},
		Neutral:       RegionCode{4, 1},
		Offensive:     RegionCode{5, 1},
		OnAttack:      RegionCode{6, 1},
	},
	"11": {
		UnderPressure: RegionCode{2, 2},
		Defensive:     RegionCode{3, 2},
		Neutral:       RegionCode{4, 2},
		Offensive:     RegionCode{5, 2},
		OnAttack:      RegionCode{6, 2},
	},
}
