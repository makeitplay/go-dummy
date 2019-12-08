package bot

import (
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/field"
)

// IMPORTANT: all this constant sets below may be changed (see each set instructions). However, any change will
// affect the tactic defined in tactic.go file. So you must go there and adapt your tactics to your new settings.

type TeamState string

type Role string

// Do not remove, rename, or add constants here.
// You however, may increase or decrease their values to change the precision of the Positioner.
// These values defines how the field will be divided by the Positioner to create a field map.
const (
	RegionCols = 8
	RegionRows = 4
)

const (
	Initial       TeamState = "initial"
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
	DistanceNear = field.FieldWidth / 8
	DistanceFar  = DistanceNear * 3
)

type RegionCode struct {
	Col uint8
	Row uint8
}

type RegionMap map[TeamState]RegionCode

func DefineRegionMap(config lugo4go.Config) RegionMap {
	return roleMap[config.Number]
}
