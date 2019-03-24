package strategy

import (
	"errors"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"math"
	"sort"
)

func DetermineMyState(turn client.TurnContext) PlayerState {
	if turn.GameMsg().Ball().Holder == nil {
		return DisputingTheBall
	} else if turn.GameMsg().Ball().Holder.TeamPlace == turn.Player().TeamPlace {
		if turn.Player().IHoldTheBall(turn.GameMsg().Ball()) {
			return HoldingTheBall
		}
		return Supporting
	}
	return Defending
}

func DetermineMyTeamState(turn client.TurnContext, TeamBallPossession arena.TeamPlace) TeamState {
	ballRegionCode := GetRegionCode(turn.GameMsg().Ball().Coords, turn.Player().TeamPlace)
	if TeamBallPossession != turn.Player().TeamPlace {
		if ballRegionCode.X < 3 {
			return UnderPressure
		} else if ballRegionCode.X < 5 {
			return Defensive
		} else if ballRegionCode.X < 7 {
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

func FindBestPointInterceptBall(ball client.Ball, player *client.Player) (speed float64, target physics.Point) {
	if ball.Velocity.Speed == 0 {
		return units.PlayerMaxSpeed, ball.Coords
	} else {
		calcBallPos := func(frame int) *physics.Point {
			//S = So + VT + (aT^2)/2
			V := ball.Velocity.Speed
			T := float64(frame)
			a := -units.BallDeceleration
			distance := V*T + (a*math.Pow(T, 2))/2
			if distance <= 0 {
				return nil
			}
			vectorToBal, _ := ball.Velocity.Direction.Copy().SetLength(distance)
			ballTarget := vectorToBal.TargetFrom(ball.Coords)
			return &ballTarget
		}
		frames := 1
		lastBallPosition := ball.Coords
		for {
			ballLocation := calcBallPos(frames)
			if ballLocation == nil {
				break
			}
			minDistanceToTouch := ballLocation.DistanceTo(player.Coords) - ((units.BallSize + units.PlayerSize) / 2)

			if minDistanceToTouch <= float64(units.PlayerMaxSpeed*frames) {
				if frames > 1 {
					return units.PlayerMaxSpeed, *ballLocation
				} else {
					return player.Coords.DistanceTo(*ballLocation), *ballLocation
				}
			}
			lastBallPosition = *ballLocation
			frames++
		}
		return units.PlayerMaxSpeed, lastBallPosition
	}
}

// watchOpponentOnMyRoute returns a list for obstacle between the player an it's target sorted by the distance to it
func WatchOpponentOnMyRoute(origin, target physics.Point, margin float64, opponentTeam client.Team) ([]physics.Point, error) {
	collisionPoints := []physics.Point{}
	distanceToTarget := origin.DistanceTo(target)
	vectorExpected, err := physics.NewVector(origin, target)
	if err != nil {
		return nil, fmt.Errorf("cannot find obstacles between the two points: %s", err)
	}
	for _, opponent := range opponentTeam.Players {
		if opponent.Coords.DistanceTo(origin) < distanceToTarget {
			collisionPoint := opponent.VectorCollides(*vectorExpected, origin, margin)
			if collisionPoint != nil {
				collisionPoints = append(collisionPoints, *collisionPoint)
			}
		}

	}
	sort.Slice(collisionPoints, func(i, j int) bool {
		return collisionPoints[i].DistanceTo(origin) < collisionPoints[j].DistanceTo(origin)
	})

	return collisionPoints, nil
}

// QuadraticResults resolves a quadratic function returning the x1 and x2
func QuadraticResults(a, b, c float64) (float64, float64, error) {
	if a == 0 {
		return 0, 0, errors.New("a cannot be zero")
	}
	// delta: B^2 -4.A.C
	delta := math.Pow(b, 2) - 4*a*c
	// quadratic formula: -b +/- sqrt(delta)/2a
	t1 := (-b + math.Sqrt(delta)) / (2 * a)
	t2 := (-b - math.Sqrt(delta)) / (2 * a)
	return t1, t2, nil
}
