package strategy

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"math"
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
