// Copyright 2022 Artem Mikheev

package activities

import (
	"fmt"
	"time"

	"github.com/renbou/jogmock/activities/internal/trigonometry"
)

// Record represents a single activity record.
// Only Lat, Lon and Altitude must be specified for
// building the activity, the other parameters will be
// calculated automagically.
type Record struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
	// altitude in meters
	Altitude  float64   `json:"altitude"`
	Timestamp time.Time `json:"timestamp"`
	Speed     float64   `json:"speed"`
	Distance  float64   `json:"distance"`
}

// DistanceTo calculates the distance between this and another record on the globe, in kilometres
func (r *Record) DistanceTo(other *Record) float64 {
	return trigonometry.DistanceBetweenCoords(
		trigonometry.Degree(r.Lat), trigonometry.Degree(r.Lon),
		trigonometry.Degree(other.Lat), trigonometry.Degree(other.Lon),
	)
}

// String implements the Stringer interface
func (r *Record) String() string {
	return fmt.Sprintf("@%d:%d:%d.%v: lat=%v lon=%v alt=%v speed=%v distance=%v",
		r.Timestamp.Hour(), r.Timestamp.Minute(), r.Timestamp.Second(), r.Timestamp.Nanosecond(),
		r.Lat, r.Lon, r.Altitude, r.Speed, r.Distance,
	)
}

// RecordInBetween calculates the record c that lies in between records a, b such that
// a.DistanceTo(c)/a.DistanceTo(b) == fraction
func RecordInBetween(a, b *Record, fraction float64) Record {
	altitude := a.Altitude + (b.Altitude-a.Altitude)*fraction
	lat, lon := trigonometry.CoordsBetween(
		trigonometry.Degree(a.Lat), trigonometry.Degree(a.Lon),
		trigonometry.Degree(b.Lat), trigonometry.Degree(b.Lon),
		fraction)
	return Record{
		Lat:      float64(lat),
		Lon:      float64(lon),
		Altitude: altitude,
	}
}
