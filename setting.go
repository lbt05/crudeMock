package main

type Setting struct {
	delayDistribution DelayDistribution `json:"delayDistribution"`
}

type DelayDistribution struct {
	Type string `json:"type"`
	//::TODO going to implement later
	//Median int     `json:"median"`
	//Sigma  float32 `json:"sigma"`
	Lower int `json:"lower"`
	Upper int `json:"upper"`
}

func (distribution DelayDistribution) getDelay(fixedDelay int) int {
	gen := boolGen()
	if fixedDelay > 0 {
		return fixedDelay
	} else {
		if gen.Bool() {
			return distribution.Upper
		} else {
			return distribution.Lower
		}
	}
}
