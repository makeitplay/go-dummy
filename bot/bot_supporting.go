package bot

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"sort"
)

const supportersCount = 3

func (b Bot) OnSupporting(ctx context.Context, turn coach.TurnData) error {
	debugMsg := ""
	b.LastBallHolder = turn.Snapshot.Ball.Holder.Number
	supportPlayers := findClosestPlayers(turn.Snapshot.Ball.Holder, turn.Snapshot)
	shouldSupport := false
	for _, p := range supportPlayers {
		if p.Number == turn.Me.Number {
			shouldSupport = true
			break
		}
	}
	if shouldSupport {
		target, err := FindSpotToAssist(supportPlayers, turn.Snapshot.Ball.Holder, turn.Me, b.Positioner)
		if err != nil {
			debugMsg = "I could not find a good spot to be in"
			b.log.Errorf("error finding a good stop to support the holder: %s", err)
		} else {
			moveOrder, err := field.MakeOrderMoveMaxSpeed(*turn.Me.Position, target)
			debugMsg = "getting better position"
			if target.DistanceTo(*turn.Me.Position) < DistanceNear {
				debugMsg = "already in position"
				moveOrder, err = field.MakeOrderMove(*turn.Me.Position, *turn.Snapshot.Ball.Position, 0)
			}
			if err != nil {
				return fmt.Errorf("was not able to move this turn: %s", err)
			}
			return send(ctx, turn, []proto.PlayerOrder{moveOrder}, debugMsg)
		}
	}
	// @see Readme- ignored errors
	ballRegion, _ := b.Positioner.GetPointRegion(*turn.Snapshot.Ball.Position)

	// @see Readme- ignored errors
	teamState, _ := DetermineTeamState(ballRegion, turn.Me.TeamSide, b.BallPossessionTeam)

	// @see Readme- ignored errors
	currentReg, _ := b.Positioner.GetPointRegion(*turn.Me.Position)

	expectedRegion := GetMyRegion(teamState, b.Positioner, turn.Me.Number)

	//b.log.Debugf("C %s, Ex %s", currentReg, expectedRegion)

	var moveOrder *proto.Order_Move
	var err error
	if currentReg.String() != expectedRegion.String() {
		moveOrder, err = field.MakeOrderMoveMaxSpeed(*turn.Me.Position, expectedRegion.Center())
		if err != nil {
			return fmt.Errorf("was not able to move this turn: %s", err)
		}
		debugMsg = fmt.Sprintf("moving to my region: %v", expectedRegion)
	} else {
		moveOrder, err = field.MakeOrderMove(*turn.Me.Position, *turn.Snapshot.Ball.Position, 0)
		if err != nil {
			return fmt.Errorf("was not able to move this turn: %s", err)
		}
		debugMsg = fmt.Sprintf("holding position (%s %s)", expectedRegion, currentReg)
	}
	return send(ctx, turn, []proto.PlayerOrder{moveOrder}, debugMsg)
}

func findClosestPlayers(ballHolder *proto.Player, snapshot *proto.GameSnapshot) []*proto.Player {
	ballPosition := *snapshot.Ball.Holder.Position
	l := make([]*proto.Player, len(field.GetTeam(snapshot, ballHolder.TeamSide).Players))
	copy(l, field.GetTeam(snapshot, ballHolder.TeamSide).Players)
	sort.Slice(l, func(i, j int) bool {
		return l[i].Position.DistanceTo(ballPosition) < l[j].Position.DistanceTo(ballPosition)
	})
	return l[0:supportersCount]
}

type pair struct {
	r coach.Region
	p *proto.Player
}

type arrange struct {
	Arr []pair
	Sum float64
}

// this logic is based on having only 3 supporters.
// Suggestion for improvements: find a strategic sport instead of just be besides the player
func FindSpotToAssist(supporters []*proto.Player, holder *proto.Player, me *proto.Player, positioner coach.Positioner) (proto.Point, error) {
	// Strategy: having one guy behind and two in front of the player holder
	// the most close to the defensive goal, will be behind

	holderRegion, _ := positioner.GetPointRegion(*holder.Position)
	strategicRegions := make([]coach.Region, supportersCount) // this magic number is equals to the number of supporters plus some extra possibilities
	// those are the regions we hopefully will have someone in
	strategicRegions[0] = holderRegion.Front().Left()
	strategicRegions[1] = holderRegion.Front().Right()
	strategicRegions[2] = holderRegion.Left()
	if strategicRegions[2].String() == holderRegion.String() {
		strategicRegions[2] = holderRegion.Right()
	}
	arr := []arrange{}
	for _, r := range strategicRegions {
		arrA := arrange{
			Arr: []pair{},
		}
		for _, p := range supporters {
			arrA.Arr = append(arrA.Arr, pair{r: r, p: p})
			d := r.Center()
			arrA.Sum += d.DistanceTo(*p.Position)
		}
		arr = append(arr, arrA)
	}

	// find closest region from me
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Sum < arr[j].Sum
	})

	bestArrangement := arr[0]
	for _, r := range bestArrangement.Arr {
		if r.p.Number == me.Number {
			return r.r.Center(), nil
		}
	}
	return proto.Point{}, fmt.Errorf("player should be among the supporters")
}
