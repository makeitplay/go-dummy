package dummy

import "github.com/lugobots/arena/units"

type FuzzyScale int

const (
	MustNot FuzzyScale = iota
	ShouldNot
	May
	Should
	Must
)

const (
	DistanceBeside  = units.FieldWidth / 10
	DistanceNear    = units.FieldWidth / 8
	DistanceFar     = units.FieldWidth / 6
	DistanceDistant = units.FieldWidth / 4
)
