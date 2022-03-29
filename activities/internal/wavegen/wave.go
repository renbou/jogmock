// Copyright 2022 Artem Mikheev

package wavegen

import (
	"math/rand"

	"github.com/renbou/jogmock/activities/internal/randutil"
)

// SlopeOptions are the options used to generate
// a single wave slope. Pretty much a more general version
// of SpeedOptions
type SlopeOptions struct {
	Slope     float64
	Amplitude float64
	MinPeriod int
	MaxPeriod int
}

type slopeType int

const (
	commonSlope = iota
	rareSlope
)

// Wave is a single perlin-noise wave which will be generated using the given options
type Wave struct {
	CommonSlope     SlopeOptions
	RareSlope       SlopeOptions
	RareSlopeChance float64
	// the average value of the slope
	Average float64
	// roughly equals the period of the current slope
	slopePeriod int
	// index of the current segment in the slope
	slopeIndex int
	// previous and current slope tangent
	prevSlope float64
	curSlope  float64
	// current slope type
	slopeType slopeType
}

func (w *Wave) generateSlopeType() slopeType {
	if rand.Float64() < w.RareSlopeChance {
		return rareSlope
	} else {
		return commonSlope
	}
}

func (w *Wave) generateSlope() float64 {
	if w.slopeType == rareSlope {
		// rare slope should be generated as the maximum common slope
		// size with an added random slope to the maximum rare slope
		randomRange := w.RareSlope.Slope - w.CommonSlope.Slope
		rnd := randutil.Float64InRange(-randomRange, randomRange)
		if rnd <= 0 {
			return -w.CommonSlope.Slope + rnd
		} else {
			return w.CommonSlope.Slope + rnd
		}
	} else {
		// common slope is simply a random slope with the given bounds
		return randutil.Float64InRange(-w.CommonSlope.Slope, w.CommonSlope.Slope)
	}
}

func (w *Wave) generateSlopePeriod() int {
	var min, max int
	if w.slopeType == rareSlope {
		min, max = w.RareSlope.MinPeriod, w.RareSlope.MaxPeriod
	} else {
		min, max = w.CommonSlope.MinPeriod, w.CommonSlope.MaxPeriod
	}
	return randutil.IntInRange(min, max)
}

func (w *Wave) nextCycle() {
	w.prevSlope = w.curSlope
	w.slopeIndex = 0

	w.slopeType = w.generateSlopeType()
	w.curSlope = w.generateSlope()
	w.slopePeriod = w.generateSlopePeriod()
}

func (w *Wave) currentAmplitude() float64 {
	if w.slopeType == rareSlope {
		return w.RareSlope.Amplitude
	} else {
		return w.CommonSlope.Amplitude
	}
}

func (w *Wave) Next() float64 {
	if w.slopeIndex == w.slopePeriod {
		w.nextCycle()
	}

	fraction := float64(w.slopeIndex) / float64(w.slopePeriod)
	loPos := w.prevSlope * fraction
	hiPos := -w.curSlope * (1 - fraction)
	curve := fraction * fraction * (3.0 - 2.0*fraction)
	valueModifier := (loPos * (1 - curve)) + (hiPos * curve)
	ret := w.Average + valueModifier*w.currentAmplitude()

	w.slopeIndex++
	return ret
}
