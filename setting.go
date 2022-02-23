package main

import (
	match "github.com/alexpantyukhin/go-pattern-match"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
)

type Setting struct {
	DelayDistribution DelayDistribution `json:"delayDistribution"`
	AccessLog         bool
}

type DelayDistribution struct {
	Type   string  `json:"type"`
	Median int     `json:"median"`
	Sigma  float64 `json:"sigma"`
	Lower  int     `json:"lower"`
	Upper  int     `json:"upper"`
}

func (distribution DelayDistribution) getDelay(fixedDelay int) int {

	if fixedDelay > 0 {
		return fixedDelay
	} else {
		//match delayDistribution type
		isMatched, mr := match.Match(distribution.Type).
			When("uniform", distribution.diceUniformDelay()).
			When("lognormal", distribution.diceNormal()).
			Result()
		if isMatched {
			return mr.(int)
		} else {
			return fixedDelay
		}
	}
}

func (distribution DelayDistribution) diceUniformDelay() int {
	gen := boolGen()
	if gen.Bool() {
		return distribution.Upper
	} else {
		return distribution.Lower
	}
}

func (distribution DelayDistribution) diceNormal() int {
	return int(math.Exp(distuv.Normal{Mu: 0.0, Sigma: 1}.Rand()*distribution.Sigma) * float64(distribution.Median))
}
