package dummie

import (
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
)

// GetActiveRegion finds the region when the player should be on this team state
func GetInitialRegion() strategy.RegionCode {
	return strategy.PlayerRegionMap[PlayerNumber][strategy.UnderPressure]
}

// GetActiveRegion finds the region when the player should be on this team state
func (d *Dummie) GetActiveRegion() strategy.RegionCode {
	return strategy.PlayerRegionMap[PlayerNumber][d.TeamState]
}

// myCurrentRegion finds the current region where the player is in
func (d *Dummie) myCurrentRegion() strategy.RegionCode {
	return strategy.GetRegionCode(d.Player.Coords, TeamPlace)
}

// isItInMyActiveRegion find out whether a point is within the active region
// @todo do I need this method?
func (d *Dummie) isItInMyActiveRegion(coords physics.Point, strategyState strategy.TeamState) bool {
	myRegionCode := d.GetActiveRegion()
	regionCode := strategy.GetRegionCode(coords, TeamPlace)
	return myRegionCode == regionCode
}

// isItInMyCurrentRegion find out whether a point is within the current region
// @todo do I need this method?
func (d *Dummie) isItInMyCurrentRegion(coords physics.Point, strategyState strategy.TeamState) bool {
	myRegionCode := d.myCurrentRegion()
	regionCode := strategy.GetRegionCode(coords, TeamPlace)
	return myRegionCode == regionCode
}

// GetActiveRegionCenter gets the center of the active region
func (d *Dummie) GetActiveRegionCenter() physics.Point {
	myRegionCode := d.GetActiveRegion()
	return myRegionCode.Center(TeamPlace)
}

// GetPlayersInRegion list all player in a specific region
func (d *Dummie) GetPlayersInRegion(regionCode strategy.RegionCode, team client.Team) []*client.Player {
	var playerList []*client.Player
	for _, player := range team.Players {
		if d.Player.ID() != player.ID() && strategy.GetRegionCode(player.Coords, TeamPlace) == regionCode {
			playerList = append(playerList, player)
		}
	}
	return playerList
}
