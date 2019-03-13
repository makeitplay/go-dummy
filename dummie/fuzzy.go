package dummie

import "github.com/makeitplay/arena/units"

const (
	// DisputingMust means: the player MUST disput the ball
	DisputingMust = iota
	// DisputingShould means: the player SHOULD disput the ball if it can
	DisputingShould
	// DisputingMay means: the player MAY disput the ball if it want
	DisputingMay
	// DisputingShouldNot means: the player SHOULD not disput, but it's ok if it does
	DisputingShouldNot
	// DisputingMustNot means: the player MUST not disput because it should let someone elso do
	DisputingMustNot
)

const (
	DistanceBeside  = units.FieldWidth / 10
	DistanceNear    = units.FieldWidth / 8
	DistanceFar     = units.FieldWidth / 6
	DistanceDistant = units.FieldWidth / 4
)
