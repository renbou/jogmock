// Copyright 2022 Artem Mikheev

package randutil

import "math/rand"

func Float64InRange(a, b float64) float64 {
	return rand.Float64()*(b-a) + a
}

func IntInRange(a, b int) int {
	return a + rand.Intn(b-a+1)
}
