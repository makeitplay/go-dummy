package bot

import (
	"github.com/lugobots/lugo4go/v2"
)

type Role string

type TeamState string

type RegionCode struct {
	Col uint8
	Row uint8
}

type RegionMap map[TeamState]RegionCode

func DefineRegionMap(config lugo4go.Config) RegionMap {
	return roleMap[config.Number]
}
