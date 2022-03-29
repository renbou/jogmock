// Copyright 2022 Artem Mikheev

// Package activities provides activity construction and generation.

package activities

import (
	"errors"
	"fmt"
	"time"

	"github.com/renbou/jogmock/activities/internal/randutil"
	"github.com/renbou/jogmock/activities/internal/wavegen"
)

// SpeedOptions represents the options used to generate
// the slopes of the speed perlin-noise wave during the activity.
type SpeedOptions struct {
	// absolute value of the tangent of slopes,
	// generated slopes will have their tangent in the range (-Slope, Slope)
	Slope float64
	// amplitude of the speed during a slope in km/h
	Amplitude float64
	// minimum duration of the slope in seconds
	MinDuration int
	// minimum duration of the slope in seconds
	MaxDuration int
}

func (o *SpeedOptions) validate() error {
	if o.Slope > 1.57 {
		return errors.New("speed slope is too high (over pi/2)")
	}

	if o.Amplitude < 0 {
		return errors.New("amplitude cannot be negative")
	}

	if o.MinDuration > o.MaxDuration {
		return errors.New("speed slope min duration cannot be higher than the max duration")
	}

	if o.MinDuration < 5 {
		return errors.New("speed slope duration cannot be less than 5 seconds")
	}

	return nil
}

// ActivityOptions contains the optional and required options
// for initializing a new Activity. Most useful optional fields
// are CommonSpeed and RareSpeed, which are used to configure
// the generation of speed wave slopes for the activity. Since during
// an activity the speed is never constant, perlin-noise waves are used
// to generate the speed, using the CommonSpeed options most of the time
// and RareSpeed occasionally to add some diversity. Thus,
// CommonSpeed should configure smooth and slow transitions, whereas
// RareSpeed should configure rough and immediate transitions.
type ActivityOptions struct {
	// Required, the activity type. Currently only Run and Ride are supported
	Type ActivityType
	// Required, the activity starting time
	Start time.Time
	// Required, the desired average speed of the activity in km/h
	DesiredSpeed float64
	// Options for generating common speed slopes
	CommonSpeed *SpeedOptions
	// Options for generating common speed slopes
	RareSpeed *SpeedOptions
	// The chance of RareSpeed options being selected
	// for the generation of the next slope
	RareSpeedChance float64
	// Duration of the speed fade-in and fade-out which happen during
	// activity start and end. This effectively changes the first FadeDuration
	// records at the beginning and the last FadeDuration at the end,
	// making it seem like the activity speed began and ended smoothly at 0.
	FadeDuration time.Duration
}

const (
	minDesiredSpeed    float64 = 1
	maxRareSpeedChance float64 = 0.5
	minFadeDuration            = 20 * time.Second
	// minFadeDuration - fadeRandomAmplitude should be > 0
	fadeRandomAmplitude            = 10
	DefaultRareSpeedChance float64 = 0.1
	DefaultFadeDuration            = 45 * time.Second
)

var (
	DefaultRunCommonSpeed = SpeedOptions{
		Slope:       0.4,
		Amplitude:   1.5,
		MinDuration: 20,
		MaxDuration: 50,
	}

	DefaultRunRareSpeed = SpeedOptions{
		Slope:       1,
		Amplitude:   2.5,
		MinDuration: 15,
		MaxDuration: 30,
	}

	DefaultRideCommonSpeed = SpeedOptions{
		Slope:       0.6,
		Amplitude:   2.5,
		MinDuration: 30,
		MaxDuration: 60,
	}

	DefaultRideRareSpeed = SpeedOptions{
		Slope:       1.2,
		Amplitude:   4,
		MinDuration: 20,
		MaxDuration: 40,
	}
)

func (o *ActivityOptions) validateAndSetDefaults() error {
	if !o.Type.IsKnown() {
		return ErrUnknownActivityType
	}

	now := time.Now()
	if o.Start.Year() < now.Year()-10 {
		return errors.New("start time is over 10 years")
	}
	if o.Start.Year() > now.Year() {
		return errors.New("start time is in the future")
	}

	if o.DesiredSpeed < minDesiredSpeed {
		return fmt.Errorf("desired speed is less than %fkm/h", minDesiredSpeed)
	}

	if o.RareSpeedChance != 0 {
		if o.RareSpeedChance > maxRareSpeedChance {
			return fmt.Errorf("rare speed chance is too high (over %f%%)", maxRareSpeedChance*100)
		} else if o.RareSpeedChance < 0 {
			return errors.New("rare speed chance cannot be negative")
		}
	} else {
		o.RareSpeedChance = DefaultRareSpeedChance
	}

	if o.FadeDuration != 0 {
		if o.FadeDuration < minFadeDuration {
			return fmt.Errorf("fade duration too small (less than %s seconds)", minFadeDuration)
		}
	} else {
		o.FadeDuration = DefaultFadeDuration
	}

	if o.CommonSpeed == nil {
		if o.Type == RunActivity {
			o.CommonSpeed = &DefaultRunCommonSpeed
		} else {
			o.CommonSpeed = &DefaultRideCommonSpeed
		}
	}

	if o.RareSpeed == nil {
		if o.Type == RunActivity {
			o.RareSpeed = &DefaultRunRareSpeed
		} else {
			o.RareSpeed = &DefaultRideRareSpeed
		}
	}

	return nil
}

// Activity represents a single activity of any valid type
type Activity struct {
	activityType    ActivityType
	startTime       time.Time
	desiredSpeed    float64
	wave            wavegen.Wave
	fadeInDuration  time.Duration
	fadeOutDuration time.Duration
	records         []Record
}

func randomiseFade(fadeDuration time.Duration) time.Duration {
	return fadeDuration + time.Second*time.Duration(randutil.IntInRange(-fadeRandomAmplitude, fadeRandomAmplitude))
}

// NewActivity initializes a new activity with the given options
func NewActivity(options *ActivityOptions) (*Activity, error) {
	if err := options.validateAndSetDefaults(); err != nil {
		return nil, err
	}

	return &Activity{
		activityType: options.Type,
		startTime:    options.Start,
		desiredSpeed: options.DesiredSpeed,
		wave: wavegen.Wave{
			CommonSlope: wavegen.SlopeOptions{
				Slope:     options.CommonSpeed.Slope,
				Amplitude: options.CommonSpeed.Amplitude,
				MinPeriod: options.CommonSpeed.MinDuration,
				MaxPeriod: options.CommonSpeed.MaxDuration,
			},
			RareSlope: wavegen.SlopeOptions{
				Slope:     options.RareSpeed.Slope,
				Amplitude: options.RareSpeed.Amplitude,
				MinPeriod: options.RareSpeed.MinDuration,
				MaxPeriod: options.RareSpeed.MaxDuration,
			},
			RareSlopeChance: options.RareSpeedChance,
			Average:         options.DesiredSpeed,
		},
		fadeInDuration:  randomiseFade(options.FadeDuration),
		fadeOutDuration: randomiseFade(options.FadeDuration),
	}, nil
}

func (a *Activity) lastRecord() *Record {
	if len(a.records) == 0 {
		return nil
	}
	return &a.records[len(a.records)-1]
}

// TotalDistance returns the total distance traversed during the activity
func (a *Activity) TotalDistance() float64 {
	if r := a.lastRecord(); r == nil {
		return 0
	} else {
		return r.Distance
	}
}

// TotalTime returns the total duration of the activity
func (a *Activity) TotalDuration() time.Duration {
	if r := a.lastRecord(); r == nil {
		return 0
	} else {
		return r.Timestamp.Sub(a.startTime)
	}
}

// intermediateRecord builds an intermediate record which is a result
// of moving from prev to next with given speed, and returns the new record
// as well as a boolean representing whether the next record was reached
// (in which case the returned record's coordinates will be equal to next's)
func intermediateRecord(prev *Record, next *Record, speed float64) (intermediate Record, reached bool) {
	distanceBetweenRecords := prev.DistanceTo(next)
	distanceTraveledIn1Second := speed / 3600

	if distanceTraveledIn1Second > distanceBetweenRecords {
		// if we reach the point, then we need to just copy it and correct the timestamp
		reached = true
		intermediate = Record{
			Lat:      next.Lat,
			Lon:      next.Lon,
			Altitude: next.Altitude,
			Distance: prev.Distance + distanceBetweenRecords,
			Speed:    speed,
			Timestamp: prev.Timestamp.Add(time.Duration(
				float64(time.Hour) * (distanceBetweenRecords / speed),
			)),
		}
	} else {
		// otherwise calculate the intermediate point values
		reached = false
		fraction := distanceTraveledIn1Second / distanceBetweenRecords
		intermediate = RecordInBetween(prev, next, fraction)
		intermediate.Distance = prev.Distance + distanceTraveledIn1Second
		intermediate.Speed = speed
		intermediate.Timestamp = prev.Timestamp.Add(time.Second)
	}
	return
}

func calculateFadeSegmentSpeed(fadeDuration time.Duration, speed float64) float64 {
	fadeSeconds := float64(fadeDuration / time.Second)
	return speed / fadeSeconds
}

var (
	ErrInvalidLatitude  = errors.New("latitude isn't in bounds (-90, 90)")
	ErrInvalidLongitude = errors.New("longitude isn't in bounds (-180, 180)")
)

// AddRecord adds the given record to the activity. Only the Lat, Lon and Altitude
// fields of the Record are used, other fields are calculated depending on the activity parameters
func (a *Activity) AddRecord(record *Record) error {
	if record.Lat < -90 || record.Lat > 90 {
		return ErrInvalidLatitude
	} else if record.Lon < -180 || record.Lon > 180 {
		return ErrInvalidLongitude
	}

	if len(a.records) == 0 {
		// the first record with speed set to 0 for fade-in
		a.records = append(a.records, Record{
			Lat:       record.Lat,
			Lon:       record.Lon,
			Altitude:  record.Altitude,
			Timestamp: a.startTime,
			Distance:  0,
			Speed:     0,
		})
		return nil
	}

	var intermediate Record
	var speed float64
	reached := false
	for !reached {
		if a.TotalDuration() < a.fadeInDuration {
			// add the fade-in records if we haven't reached our desired speed yet
			fadeInSegmentSpeed := calculateFadeSegmentSpeed(a.fadeInDuration, a.desiredSpeed)
			speed = fadeInSegmentSpeed * float64(a.TotalDuration()/time.Second)
		} else {
			speed = a.wave.Next()
		}
		intermediate, reached = intermediateRecord(a.lastRecord(), record, speed)
		a.records = append(a.records, intermediate)
	}

	return nil
}

// addFadeOut should be called once we are done adding records
func (a *Activity) addFadeOut() error {
	if a.TotalDuration() < a.fadeOutDuration {
		return errors.New("unable to add fade-out to an activity shorter than the fade duration")
	}

	// remove the last records and add them with a faded-out speed
	fadeOutSeconds := int(a.fadeOutDuration / time.Second)
	recBeforeFadeout := &a.records[len(a.records)-fadeOutSeconds]
	fadeOutRecords := a.records[len(a.records)-fadeOutSeconds+1:]
	a.records = a.records[:len(a.records)-fadeOutSeconds+1]

	fadeOutSegmentSpeed := calculateFadeSegmentSpeed(a.fadeOutDuration, recBeforeFadeout.Speed)
	calcSegmentSpeed := func() float64 {
		return fadeOutSegmentSpeed * (float64(fadeOutSeconds) - float64(
			(a.lastRecord().Timestamp.Sub(recBeforeFadeout.Timestamp))/time.Second,
		))
	}

	prev := recBeforeFadeout
	for _, record := range fadeOutRecords {
		var intermediate Record
		reached := false
		for !reached {
			intermediate, reached = intermediateRecord(prev, &record, calcSegmentSpeed())
			a.records = append(a.records, intermediate)
		}
	}

	return nil
}

// BuildRecords finalizes the activity by adding the fade-out
// and returns the constructed list of records
func (a *Activity) BuildRecords() ([]Record, error) {
	if err := a.addFadeOut(); err != nil {
		return nil, err
	}
	return a.records, nil
}
