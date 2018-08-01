package brain

import (
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/the-dummies-go/strategy"
	"math"
	"sort"
)

// PointCollection creates a list of points
type PointCollection []Physics.Point

// Len implements the
func (s PointCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// SortByDistance allows to sort a list of points
// @todo this struct may be removed if we use the sort.Slice function
type SortByDistance struct {
	PointCollection
	From Physics.Point
}

// watchOpponentOnMyRoute returns a list for obstacle between the player an it's target sorted by the distance to it
func watchOpponentOnMyRoute(status client.GameInfo, player *client.Player, target Physics.Point) PointCollection {
	opponentTeam := player.GetOpponentTeam(status)
	collisionPoints := SortByDistance{From: player.Coords}

	vectorExpected := Physics.NewVector(player.Coords, target)
	for _, opponent := range opponentTeam.Players {
		collisionPoint := opponent.VectorCollides(*vectorExpected, player.Coords, float64(player.Size)/2)

		if collisionPoint != nil {
			collisionPoints.PointCollection = append(collisionPoints.PointCollection, *collisionPoint)
		}
	}
	return collisionPoints.PointCollection
}

// QuadraticResults resolves a quadratic function returning the x1 and x2
func QuadraticResults(a, b, c float64) (float64, float64) {
	// delta: B^2 -4.A.C
	delta := math.Pow(b, 2) - 4*a*c
	// quadratic formula: -b +/- sqrt(delta)/2a
	t1 := (-b + math.Sqrt(delta)) / (2 * a)
	t2 := (-b - math.Sqrt(delta)) / (2 * a)
	return t1, t2
}

// FindBestPointInRegionToAssist finds the best point to support the ball holder from within a region
func FindBestPointInRegionToAssist(gameMessage client.GameMessage, region strategy.RegionCode, assisted *client.Player) (target Physics.Point) {
	centerPoint := region.Center(assisted.TeamPlace)
	vctToCenter := Physics.NewVector(assisted.Coords, centerPoint).SetLength(strategy.RegionWidth)
	obstacles := watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, vctToCenter.TargetFrom(assisted.Coords))
	if len(obstacles) == 0 {
		return vctToCenter.TargetFrom(assisted.Coords)
	} else {
		initialVector := vctToCenter
		avoidObstacles := func(ang float64) bool {
			tries := 3
			for tries > 0 {
				vctToCenter.AddAngleDegree(ang)
				target = vctToCenter.TargetFrom(assisted.Coords)
				if region != strategy.GetRegionCode(target, assisted.TeamPlace) {
					//too far
					tries = 0
				}
				obstacles = watchOpponentOnMyRoute(gameMessage.GameInfo, assisted, target)
				tries--
				if len(obstacles) <= 0 {
					return true
				}
			}
			return false
		}

		if !avoidObstacles(10) && !avoidObstacles(-10) {
			target = initialVector.TargetFrom(assisted.Coords)
		}
	}
	return
}

// FindSpotToAssist finds a good region to support the ball holder
func FindSpotToAssist(gameMessage client.GameMessage, assisted *client.Player, assistant *Brain, offensively bool) strategy.RegionCode {
	var availableSpots []strategy.RegionCode
	var spotList []strategy.RegionCode
	if offensively {
		spotList = ListSpotsCandidatesToOffensiveAssistance(assisted, assistant)
	} else {
		spotList = ListSpotsCandidatesToDefensiveAssistance(assisted, assistant)
	}
	for _, region := range spotList {
		mateInTheRegion := assistant.GetPlayersInRegion(region, assistant.FindMyTeamStatus(gameMessage.GameInfo))
		if len(mateInTheRegion) == 0 {
			availableSpots = append(availableSpots, region)
		} else if region == assistant.GetActiveRegion(TeamState) {
			// eu to no meu canto, me deixe em paz
			availableSpots = append(availableSpots, region)
		} else {
			frankenstein := Brain{Player: mateInTheRegion[0]}
			isHimTheOwner := region == frankenstein.GetActiveRegion(TeamState)
			if !isHimTheOwner && assistant.myCurrentRegion() == region {
				// two invasors disputing
				myDistanceToTheBall := assistant.Coords.DistanceTo(assisted.Coords)
				invasorDistanceToTheBall := assistant.Coords.DistanceTo(mateInTheRegion[0].Coords)
				if myDistanceToTheBall < invasorDistanceToTheBall {
					availableSpots = append(availableSpots, region)
				}
			}
		}
	}
	sort.Slice(availableSpots, func(a, b int) bool {
		teamStatus := assistant.GetOpponentTeam(gameMessage.GameInfo)
		opponentsInA := len(assistant.GetPlayersInRegion(availableSpots[a], teamStatus))
		opponentsInB := len(assistant.GetPlayersInRegion(availableSpots[b], teamStatus))

		distanceToA := math.Round(assistant.Coords.DistanceTo(availableSpots[a].Center(assistant.TeamPlace)) / strategy.RegionWidth)
		distanceToB := math.Round(assistant.Coords.DistanceTo(availableSpots[b].Center(assistant.TeamPlace)) / strategy.RegionWidth)

		distanceAToAssistant := math.Round(assisted.Coords.DistanceTo(availableSpots[a].Center(assistant.TeamPlace)) / strategy.RegionWidth)
		distanceBToAssistant := math.Round(assisted.Coords.DistanceTo(availableSpots[b].Center(assistant.TeamPlace)) / strategy.RegionWidth)

		APoints := distanceToB - distanceToA
		APoints += float64(opponentsInB - opponentsInA)
		APoints += distanceBToAssistant - distanceAToAssistant
		APoints += float64(availableSpots[a].X-availableSpots[b].X) * 2.5
		return APoints >= 0
	})

	if len(availableSpots) > 0 {
		return availableSpots[0]
	}
	return assistant.GetActiveRegion(TeamState)
}

// ListSpotsCandidatesToOffensiveAssistance List the best regions around the ball holder to help him in a offensive support (closer to the attack)
func ListSpotsCandidatesToOffensiveAssistance(assisted *client.Player, assistant *Brain) []strategy.RegionCode {
	spotCollection := []strategy.RegionCode{}
	currentRegion := strategy.GetRegionCode(assisted.Coords, assistant.TeamPlace)

	front := currentRegion.Forwards()
	if front != currentRegion {
		spotCollection = append(spotCollection, front)
	}

	assistantActiveRegion := assistant.GetActiveRegion(TeamState)

	goodRegionA := front.Left()
	if currentRegion != front && goodRegionA.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, goodRegionA)
	}
	goodRegionB := front.Right()
	if currentRegion != front && goodRegionB.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, goodRegionB)
	}

	fairRegionA := currentRegion.Left()
	if currentRegion != fairRegionA && fairRegionA.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, fairRegionA)
	}
	fairRegionB := currentRegion.Right()
	if currentRegion != fairRegionB && fairRegionB.ChessDistanceTo(assistantActiveRegion) < 2 {
		spotCollection = append(spotCollection, fairRegionB)
	}
	return spotCollection
}

// ListSpotsCandidatesToDefensiveAssistance List the best regions around the ball holder to help him in a defensive support
func ListSpotsCandidatesToDefensiveAssistance(assisted *client.Player, assistant *Brain) []strategy.RegionCode {
	spotCollection := []strategy.RegionCode{}
	currentRegion := strategy.GetRegionCode(assisted.Coords, assistant.TeamPlace)

	back := currentRegion.Backwards()
	if back != currentRegion {
		spotCollection = append(spotCollection, back)
	}

	goodRegionA := back.Left()
	if currentRegion != back {
		spotCollection = append(spotCollection, goodRegionA)
	}
	goodRegionB := back.Right()
	if currentRegion != back {
		spotCollection = append(spotCollection, goodRegionB)
	}

	return spotCollection
}
