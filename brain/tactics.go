package brain

import (
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
)

// GetActiveRegion finds the region when the player should be on this team state
func (b *Brain) GetActiveRegion(state strategy.TeamState) strategy.RegionCode {
	return strategy.PlayerRegionMap[b.Number][state]
}

// myCurrentRegion finds the current region where the player is in
func (b *Brain) myCurrentRegion() strategy.RegionCode {
	return strategy.GetRegionCode(b.Coords, b.TeamPlace)
}

// isItInMyActiveRegion find out whether a point is within the active region
// @todo do I need this method?
func (b *Brain) isItInMyActiveRegion(coords physics.Point, strategyState strategy.TeamState) bool {
	myRegionCode := b.GetActiveRegion(strategyState)
	regionCode := strategy.GetRegionCode(coords, b.TeamPlace)
	return myRegionCode == regionCode
}

// isItInMyCurrentRegion find out whether a point is within the current region
// @todo do I need this method?
func (b *Brain) isItInMyCurrentRegion(coords physics.Point, strategyState strategy.TeamState) bool {
	myRegionCode := b.myCurrentRegion()
	regionCode := strategy.GetRegionCode(coords, b.TeamPlace)
	return myRegionCode == regionCode
}

// GetActiveRegionCenter gets the center of the active region
func (b *Brain) GetActiveRegionCenter(strategyState strategy.TeamState) physics.Point {
	myRegionCode := b.GetActiveRegion(strategyState)
	return myRegionCode.Center(b.TeamPlace)
}

// GetPlayersInRegion list all player in a specific region
func (b *Brain) GetPlayersInRegion(regionCode strategy.RegionCode, team client.Team) []*client.Player {
	var playerList []*client.Player
	for _, player := range team.Players {
		if b.ID() != player.ID() && strategy.GetRegionCode(player.Coords, b.TeamPlace) == regionCode {
			playerList = append(playerList, player)
		}
	}
	return playerList
}
