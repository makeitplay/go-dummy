package bot

const (
	Initial       TeamState = "initial"
	UnderPressure TeamState = "under-pressure"
	Defensive     TeamState = "defensive"
	Neutral       TeamState = "neutral"
	Offensive     TeamState = "offensive"
	OnAttack      TeamState = "on-attack"
)
const (
	RegionCols = 8
	RegionRows = 4
)
const (
	Defense Role = "defense"
	Middle  Role = "middle"
	Attack  Role = "attack"
)

func DefineRole(number uint32) Role {
	if number < 6 {
		return Defense
	} else if number < 10 {
		return Middle
	}
	return Attack
}

var roleMap = map[uint32]RegionMap{
	2: {
		Initial:       {1, 0},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	3: {
		Initial:       {1, 1},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	4: {
		Initial:       {1, 2},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	5: {
		Initial:       {1, 3},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	6: {
		Initial:       {2, 0},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	7: {
		Initial:       {2, 1},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	8: {
		Initial:       {2, 2},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	9: {
		Initial:       {2, 3},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	10: {
		Initial:       {3, 1},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
	11: {
		Initial:       {3, 2},
		UnderPressure: {2, 3},
		Defensive:     {2, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {2, 3},
	},
}
