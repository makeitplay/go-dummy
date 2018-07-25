package brain

import (
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/go-dummy/strategy"
	"github.com/makeitplay/client-player-go/Game"
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

func (b *Brain) GetActiveRegionCenter(strategyState strategy.TeamState) Physics.Point {
	myRegionCode := b.GetActiveRegion(strategyState)
	return myRegionCode.Center(b.TeamPlace)
}

func (b *Brain) DetermineMyTeamState(msg Game.GameMessage) strategy.TeamState {
	ballRegionCode := strategy.GetRegionCode(msg.GameInfo.Ball.Coords, b.TeamPlace)
	if TeamBallPossession != b.TeamPlace {
		if ballRegionCode.X < 3 {
			return strategy.UnderPressure
		} else if ballRegionCode.X < 5 {
			return strategy.Defensive
		} else if ballRegionCode.X < 7 {
			return strategy.Neutral
		} else {
			return strategy.Offensive
		}
	} else {
		if ballRegionCode.X < 2 {
			return strategy.Defensive
		} else if ballRegionCode.X < 4 {
			return strategy.Neutral
		} else if ballRegionCode.X < 6 {
			return strategy.Offensive
		} else {
			return strategy.OnAttack
		}
	}
}

func (b *Brain) GetPlayersInRegion(regionCode strategy.RegionCode, team Game.Team) []*Game.Player {
	var playerList []*Game.Player
	for _, player := range team.Players {
		if b.ID() != player.ID() && strategy.GetRegionCode(player.Coords, b.TeamPlace) == regionCode {
			playerList = append(playerList, player)
		}
	}
	return playerList
}