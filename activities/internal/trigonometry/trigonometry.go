// Copyright 2022 Artem Mikheev

package trigonometry

import "math"

type (
	Degree float64
	Radian float64
)

// haversine(θ) function
func haversine(theta Radian) float64 {
	return math.Pow(math.Sin(float64(theta)/2), 2)
}

func (d Degree) ToRadians() Radian {
	return Radian(d * math.Pi / 180)
}

func (r Radian) ToDegrees() Degree {
	return Degree(r * 180 / math.Pi)
}

func RadCos(r Radian) float64 {
	return math.Cos(float64(r))
}

func RadSin(r Radian) float64 {
	return math.Sin(float64(r))
}

// http://en.wikipedia.org/wiki/Haversine_formula
func DistanceBetweenCoords(lat1, lon1, lat2, lon2 Degree) float64 {
	la1, lo1 := lat1.ToRadians(), lon1.ToRadians()
	la2, lo2 := lat2.ToRadians(), lon2.ToRadians()

	earthRadiusM := float64(6378100)

	h := haversine(la2-la1) + RadCos(la1)*RadCos(la2)*haversine(lo2-lo1)

	return 2 * earthRadiusM * math.Asin(math.Sqrt(h)) / 1000
}

func CoordsBetween(lat1, lon1, lat2, lon2 Degree, fraction float64) (Degree, Degree) {
	// https://stackoverflow.com/questions/33907276/calculate-point-between-two-coordinates-based-on-a-percentage
	la1, lo1 := lat1.ToRadians(), lon1.ToRadians()
	la2, lo2 := lat2.ToRadians(), lon2.ToRadians()

	deltaLat := la2 - la1
	deltaLon := lo2 - lo1

	calcA := math.Pow(RadSin(deltaLat/2), 2) + RadCos(la1)*RadCos(la2)*math.Pow(RadSin(deltaLon/2), 2)
	calcB := 2 * math.Atan2(math.Sqrt(calcA), math.Sqrt(1-calcA))

	A := math.Sin((1-fraction)*calcB) / math.Sin(calcB)
	B := math.Sin(fraction*calcB) / math.Sin(calcB)

	x := A*RadCos(la1)*RadCos(lo1) + B*RadCos(la2)*RadCos(lo2)
	y := A*RadCos(la1)*RadSin(lo1) + B*RadCos(la2)*RadSin(lo2)
	z := A*RadSin(la1) + B*RadSin(la2)

	lat3 := Radian(math.Atan2(z, math.Sqrt(x*x+y*y)))
	lon3 := Radian(math.Atan2(y, x))

	return lat3.ToDegrees(), lon3.ToDegrees()
}
