package bot

import (
	"github.com/lugobots/lugo4go/v2"
)

type Role string

type RegionCode struct {
	Col uint8
	Row uint8
}
type TeamState string

type RegionMap map[TeamState]RegionCode

const (
	UnderPressure TeamState = "under-pressure"
	Defensive     TeamState = "defensive"
	Neutral       TeamState = "neutral"
	Offensive     TeamState = "offensive"
	OnAttack      TeamState = "on-attack"
)

const (
	Defense Role = "defense"
	Middle  Role = "middle"
	Attack  Role = "attack"
)

const (
	RegionCols = 6
	RegionRows = 4
)

var roleMap = map[uint32]RegionMap{
	2: {
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
}

func DefineRegionMap(config lugo4go.Config) RegionMap {
	return roleMap[config.Number]
}

func DefineRole(config lugo4go.Config) Role {

}
