package brain

import (
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/go-dummy/strategy"
)
func (b *Brain) GetActiveRegion(state strategy.TeamState) strategy.RegionCode {
	return strategy.PlayerRegionMap[b.Number][state]
}
func (b *Brain) myCurrentRegion() strategy.RegionCode {
	return strategy.GetRegionCode(b.Coords,  b.TeamPlace)
}

func (b *Brain) isItInMyActiveRegion(coords Physics.Point, strategyState strategy.TeamState) bool {
	myRegionCode := b.GetActiveRegion(strategyState)
	regionCode := strategy.GetRegionCode(coords, b.TeamPlace)
	return myRegionCode == regionCode
}

func (b *Brain) isItInMyCurrentRegion(coords Physics.Point, strategyState strategy.TeamState) bool {
	myRegionCode := b.myCurrentRegion()
	regionCode := strategy.GetRegionCode(coords, b.TeamPlace)
	return myRegionCode == regionCode
}

func (b *Brain) myActiveRegionCenter(strategyState strategy.TeamState) Physics.Point {
	myRegionCode := b.GetActiveRegion(strategyState)
	return strategy.GetRegionCenter(myRegionCode,  b.TeamPlace)
}



