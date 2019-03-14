package brain

import (
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/the-dummies-go/strategy"
	"math"
)

// QuadraticResults resolves a quadratic function returning the x1 and x2
func QuadraticResults(a, b, c float64) (float64, float64) {
	// delta: B^2 -4.A.C
	delta := math.Pow(b, 2) - 4*a*c
	// quadratic formula: -b +/- sqrt(delta)/2a
	t1 := (-b + math.Sqrt(delta)) / (2 * a)
	t2 := (-b - math.Sqrt(delta)) / (2 * a)
	return t1, t2
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
