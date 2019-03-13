package strategy

import (
	"github.com/makeitplay/arena"
	"strconv"
)

type TeamState string

const (
	UnderPressure TeamState = "under-pressure"
	Defensive     TeamState = "defensive"
	Neutral       TeamState = "neutral"
	Offensive     TeamState = "offensive"
	OnAttack      TeamState = "on-attack"
)

// PlayerState defines states specific for players
type PlayerState string

const (
	// Supporting identifies the player supporting the team mate
	Supporting PlayerState = "supporting"
	// HoldingTheBall identifies the player holding	the ball
	HoldingTheBall PlayerState = "holding"
	// Defending identifies the player defending against the opponent team
	Defending PlayerState = "defending"
	// DisputingTheBall identifies the player disputing the ball
	DisputingTheBall PlayerState = "disputing"
)

type Region map[TeamState]RegionCode

type PlayerRule int

const (
	DefensePlayer PlayerRule = iota
	MiddlePlayer
	AttackPlayer
)

var PlayerRegionMap = map[arena.PlayerNumber]Region{
	//defense players
	"2": {
		UnderPressure: RegionCode{0, 0},
		Defensive:     RegionCode{1, 0},
		Neutral:       RegionCode{2, 0},
		Offensive:     RegionCode{3, 0},
		OnAttack:      RegionCode{4, 0},
	},
	"3": {
		UnderPressure: RegionCode{0, 1},
		Defensive:     RegionCode{1, 1},
		Neutral:       RegionCode{2, 1},
		Offensive:     RegionCode{3, 1},
		OnAttack:      RegionCode{4, 1},
	},
	"4": {
		UnderPressure: RegionCode{0, 2},
		Defensive:     RegionCode{1, 2},
		Neutral:       RegionCode{2, 2},
		Offensive:     RegionCode{3, 2},
		OnAttack:      RegionCode{4, 2},
	},
	"5": {
		UnderPressure: RegionCode{0, 3},
		Defensive:     RegionCode{1, 3},
		Neutral:       RegionCode{2, 3},
		Offensive:     RegionCode{3, 3},
		OnAttack:      RegionCode{4, 3},
	},

	//middle players
	"6": {
		UnderPressure: RegionCode{1, 0},
		Defensive:     RegionCode{2, 0},
		Neutral:       RegionCode{3, 0},
		Offensive:     RegionCode{4, 0},
		OnAttack:      RegionCode{5, 0},
	},
	"7": {
		UnderPressure: RegionCode{1, 1},
		Defensive:     RegionCode{2, 1},
		Neutral:       RegionCode{3, 1},
		Offensive:     RegionCode{4, 1},
		OnAttack:      RegionCode{5, 1},
	},
	"8": {
		UnderPressure: RegionCode{1, 2},
		Defensive:     RegionCode{2, 2},
		Neutral:       RegionCode{3, 2},
		Offensive:     RegionCode{4, 2},
		OnAttack:      RegionCode{5, 2},
	},
	"9": {
		UnderPressure: RegionCode{1, 3},
		Defensive:     RegionCode{2, 3},
		Neutral:       RegionCode{3, 3},
		Offensive:     RegionCode{4, 3},
		OnAttack:      RegionCode{5, 3},
	},

	//attack players
	"10": {
		UnderPressure: RegionCode{3, 1},
		Defensive:     RegionCode{4, 1},
		Neutral:       RegionCode{5, 1},
		Offensive:     RegionCode{6, 1},
		OnAttack:      RegionCode{6, 1},
	},
	"11": {
		UnderPressure: RegionCode{3, 2},
		Defensive:     RegionCode{4, 2},
		Neutral:       RegionCode{5, 2},
		Offensive:     RegionCode{6, 2},
		OnAttack:      RegionCode{6, 2},
	},
}

func DefinePlayerRule(number arena.PlayerNumber) PlayerRule {
	num, _ := strconv.Atoi(string(number))
	if num < 6 {
		return DefensePlayer
	} else if num < 10 {
		return MiddlePlayer
	} else {
		return AttackPlayer
	}
}
