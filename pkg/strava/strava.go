// Copyright 2021 Artem Mikheev

package strava

import (
	"errors"
	"math"
	"time"

	"github.com/renbou/strava-keker/pkg/fit"
	"github.com/renbou/strava-keker/pkg/fit/types"
)

type (
	StravaActivityType string
)

var (
	StravaActivityTypeRun  StravaActivityType = "Run"
	StravaActivityTypeRide StravaActivityType = "Ride"
)

// ActivityOptions represent the options when starting a new activity
type ActivityOptions struct {
	// internal identifier of the app version
	AppVersion uint32
	// full mobile application version identifier
	MobileAppVersion string
	// company that produced the device, ex "Xiaomi" or Honor
	DeviceManufacturer string
	// model of your device, ex "Redmi Note 9 Pro"
	DeviceModel string
	// identifier of android os installed on the device, ex "10" or "11"
	DeviceOsVersion string
	// activity type of the to be simulated activity
	ActivityType StravaActivityType
	// when the activity is started
	StartTime time.Time
	// speed to which the average speed should be close to
	DesiredSpeed float64
}

// StravaActivity defines a single strava activity simulated
// with the wanted parameters which can be encoded into a
// fit file with the according format (reverse engineered)
type stravaActivity struct {
	appVersion         uint32
	mobileAppVersion   string
	deviceManufacturer string
	deviceModel        string
	deviceOsVersion    string
	activityType       StravaActivityType
	startTime          time.Time
	desiredSpeed       float64
	totalActiveTime    time.Duration
	totalDistance      float64
	activityRecords    []*activityRecord
}

// NewActivity validates passed ActivityOptions and constructs
// a new activity based on those options
func NewActivity(options *ActivityOptions) (*stravaActivity, error) {
	if options.AppVersion == 0 {
		return nil, errors.New("app version must not be zero")
	}
	if options.MobileAppVersion == "" {
		return nil, errors.New("full mobile app version must be set")
	}
	if options.DeviceManufacturer == "" {
		return nil, errors.New("device manufacturer must be set")
	}
	if options.DeviceModel == "" {
		return nil, errors.New("device model must be set")
	}

	switch options.DeviceOsVersion {
	case "11":
	case "10":
	case "9":
		break
	default:
		return nil, errors.New("device os version must be set to 9, 10 or 11")
	}

	switch options.ActivityType {
	case StravaActivityTypeRun:
	case StravaActivityTypeRide:
		break
	default:
		return nil, errors.New("activity type is invalid")
	}

	// fit time begins at unix 631065600
	if options.StartTime.Add(time.Second*631065600).Year() != time.Now().Year() {
		return nil, errors.New("start time year doesn't match current year, smth must be wrong")
	}

	if options.DesiredSpeed < 3 {
		return nil, errors.New("desired speed is less than 3 km/h, seems too slow")
	}

	return &stravaActivity{
		appVersion:         options.AppVersion,
		mobileAppVersion:   options.MobileAppVersion,
		deviceManufacturer: options.DeviceManufacturer,
		deviceModel:        options.DeviceModel,
		deviceOsVersion:    options.DeviceOsVersion,
		activityType:       options.ActivityType,
		startTime:          options.StartTime,
		desiredSpeed:       options.DesiredSpeed,
		activityRecords:    make([]*activityRecord, 0, 64),
	}, nil
}

// Record represents a single record which can be added to
// an activity, for which the other parameters will then be
// calculated automagically
type Record struct {
	Lat      float64
	Lon      float64
	Altitude float64
}

// TODO first and last record mustt have speed close to zero
// TODO generate records properly
func (act *stravaActivity) AddRecord(record Record) error {
	if record.Lat < -90 || record.Lat > 90 {
		return errors.New("record latitude isn't in bounds")
	}
	if record.Lon < -180 || record.Lon > 180 {
		return errors.New("record longitude isn't in bounds")
	}
	if record.Altitude < 0 {
		return errors.New("record altitude is below sea level")
	}

	// create new activity record based on parameters
	// currently the speed is constantly set to the desired one
	actRecord := &activityRecord{
		lat:      record.Lat,
		lon:      record.Lon,
		altitude: record.Altitude,
		speed:    act.desiredSpeed,
	}

	if len(act.activityRecords) > 0 {
		// set distance using previous record
		// this should actually depend on speed+time in the future
		// (speed * 1sec), not time depending on distance like now
		prevActRecord := act.activityRecords[len(act.activityRecords)-1]
		distanceMoved := distanceBetweenCoords(
			prevActRecord.lat, prevActRecord.lon, actRecord.lat, actRecord.lon,
		)

		// currently calculating time based on distance and speed
		// should be constantly 1s
		actRecord.timestamp = prevActRecord.timestamp.Add(time.Duration(
			float64(time.Hour) * (distanceMoved / actRecord.speed),
		))
		actRecord.distance = prevActRecord.distance + distanceMoved
	} else {
		actRecord.timestamp = act.startTime
		actRecord.distance = 0
	}

	act.activityRecords = append(act.activityRecords, actRecord)
	act.totalActiveTime = actRecord.timestamp.Sub(act.startTime)
	act.totalDistance = actRecord.distance
	return nil
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// http://en.wikipedia.org/wiki/Haversine_formula
func distanceBetweenCoords(lat1, lon1, lat2, lon2 float64) float64 {
	var la1, lo1, la2, lo2, earthRadiusM float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	earthRadiusM = 6378100

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * earthRadiusM * math.Asin(math.Sqrt(h)) / 1000
}

// activityRecord represents a single record of an activity
type activityRecord struct {
	timestamp time.Time
	lat       float64
	lon       float64
	altitude  float64
	speed     float64
	distance  float64
}

const (
	FIT_MANUFACTURER_STRAVA        = 265
	FIT_PRODUCT_STRAVA_ANDROID     = 102
	FIT_FILEID_TYPE_ACTIVITY       = 4
	FIT_ACTIVITY_ONE_SESSION       = 1
	FIT_ACTIVITY_EVENT_ACTIVITY    = 26
	FIT_EVENT_TYPE_START           = 0
	FIT_EVENT_TYPE_STOP            = 1
	FIT_SESSION_EVENT_SESSION      = 8
	FIT_SESSION_ONE_LAP            = 1
	STRAVA_FIRST_LIVE_ACTIVITY_ID  = 0
	STRAVA_AUTOPAUSE_DISABLED      = 0
	FIT_LAP_EVENT_LAP              = 9
	FIT_LAP_TRIGGER_SESSION_END    = 7
	FIT_BATTERY_STATUS_GOOD        = 3
	FIT_EVENT_TIMER                = 0
	FIT_EVENT_TIMER_TRIGGER_MANUAL = 0
	STRAVA_NOICE_GPS_ACCURACY      = 4
)

func fitActivitySport(act *stravaActivity) uint8 {
	switch act.activityType {
	case StravaActivityTypeRide:
		return 2
	case StravaActivityTypeRun:
		return 1
	default:
		panic("invalid strava activity type")
	}
}

func fitEncodeDistanceKm(distance float64) int {
	return int(math.Round(distance * 10000.0))
}

func fitEncodeSpeedKmH(speed float64) int {
	return int(math.Round(speed * (float64(1000) / 3600) * 1000))
}

func fitEncodeAltitudeM(altitude float64) int {
	return int(math.Round((altitude + 500.0) * 5.0))
}

func fitEncodeCoordinate(coord float64) int {
	return int(
		math.Round(coord * float64(int64(1)<<31) / 180.0),
	)
}

func fitEncodeDuration(duration time.Duration) int {
	return int(duration.Seconds())
}

// BuildFitFile creates a fit file based on the filled activity
func (act *stravaActivity) BuildFitFile() (*fit.FitFile, error) {
	file := new(fit.FitFile)
	if err := act.writeHeader(file); err != nil {
		return nil, err
	}
	if err := act.writeBody(file); err != nil {
		return nil, err
	}
	if err := act.writeFooter(file); err != nil {
		return nil, err
	}
	return file, nil
}

func (act *stravaActivity) writeHeader(file *fit.FitFile) error {
	// add simple file id message
	fileIdMessage, err := getFileIdMessageDefinition()
	if err != nil {
		return err
	}
	fileIdMessageDef := fileIdMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(fileIdMessageDef); err != nil {
		return err
	}
	fileIdMessageData, err := fileIdMessageDef.ConstructData(FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID, FIT_FILEID_TYPE_ACTIVITY, time.Now().Unix())
	if err != nil {
		return nil
	}
	if err := file.AddMessage(fileIdMessageData); err != nil {
		return err
	}

	// add developer data id message and field definition messages
	dev := &devData{
		index:      0,
		appVersion: act.appVersion,
		fields:     make(map[string]*devField),
	}
	if err := dev.addField("live_activity_id", types.FIT_TYPE_UINT64); err != nil {
		return err
	}
	if err := dev.addField("activity_type", types.FIT_TYPE_STRING); err != nil {
		return err
	}
	if err := dev.addField("autopause_enabled", types.FIT_TYPE_ENUM); err != nil {
		return err
	}
	if err := dev.addField("mobile_app_version", types.FIT_TYPE_STRING); err != nil {
		return err
	}
	if err := dev.addField("device_model", types.FIT_TYPE_STRING); err != nil {
		return err
	}
	if err := dev.addField("device_os_version", types.FIT_TYPE_STRING); err != nil {
		return err
	}
	if err := dev.addField("device_manufacturer", types.FIT_TYPE_STRING); err != nil {
		return err
	}
	devMessages, err := dev.constructAllMessages(0)
	if err != nil {
		return err
	}
	for _, message := range devMessages {
		if err := file.AddMessage(message); err != nil {
			return err
		}
	}

	// add device info message
	deviceInfoMessage, err := getDeviceInfoMessageDefinition(act, dev)
	if err != nil {
		return err
	}
	deviceInfoMessageDef := deviceInfoMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(deviceInfoMessageDef); err != nil {
		return err
	}
	deviceInfoMessageData, err := deviceInfoMessageDef.ConstructData(FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID, act.mobileAppVersion, act.deviceManufacturer, act.deviceModel, act.deviceOsVersion)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(deviceInfoMessageData); err != nil {
		return err
	}

	// add activity message
	activityMessage, err := getActivityMessageDefinition()
	if err != nil {
		return err
	}
	activityMessageDef := activityMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(activityMessageDef); err != nil {
		return err
	}
	activityMessageData, err := activityMessageDef.ConstructData(FIT_ACTIVITY_ONE_SESSION, fitEncodeDuration(act.totalActiveTime), FIT_ACTIVITY_EVENT_ACTIVITY, FIT_EVENT_TYPE_STOP)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(activityMessageData); err != nil {
		return err
	}

	// add session message
	sessionMessage, err := getSessionMessageDefinition(act, dev)
	if err != nil {
		return err
	}
	sessionMessageDef := sessionMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(sessionMessageDef); err != nil {
		return err
	}
	sessionMessageData, err := sessionMessageDef.ConstructData(fitEncodeDuration(act.totalActiveTime), act.startTime.Unix(), fitEncodeDuration(act.totalActiveTime), fitEncodeDistanceKm(act.totalDistance), fitActivitySport(act), FIT_SESSION_EVENT_SESSION, FIT_EVENT_TYPE_STOP, FIT_SESSION_ONE_LAP, STRAVA_FIRST_LIVE_ACTIVITY_ID, act.activityType, STRAVA_AUTOPAUSE_DISABLED)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(sessionMessageData); err != nil {
		return err
	}

	// add lap message
	lapMessage, err := getLapMessageDefinition()
	if err != nil {
		return err
	}
	lapMessageDef := lapMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(lapMessageDef); err != nil {
		return err
	}
	lapMessageData, err := lapMessageDef.ConstructData(fitEncodeDuration(act.totalActiveTime), act.startTime.Unix(), fitEncodeDuration(act.totalActiveTime), fitEncodeDistanceKm(act.totalDistance), FIT_LAP_EVENT_LAP, FIT_EVENT_TYPE_STOP, fitActivitySport(act), FIT_LAP_TRIGGER_SESSION_END)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(lapMessageData); err != nil {
		return err
	}

	// add device battery info message on start of activity
	deviceInfoBatteryMessage, err := getDeviceInfoBatteryMessageDefinition()
	if err != nil {
		return err
	}
	deviceInfoBatteryMessageDef := deviceInfoBatteryMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(deviceInfoBatteryMessageDef); err != nil {
		return err
	}
	deviceInfoBatteryMessageData, err := deviceInfoBatteryMessageDef.ConstructData(FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID, act.startTime.Unix(), FIT_BATTERY_STATUS_GOOD)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(deviceInfoBatteryMessageData); err != nil {
		return err
	}

	return nil
}

func (act *stravaActivity) writeBody(file *fit.FitFile) error {
	// add event message on start of activity
	eventMessage, err := getEventMessageDefinition()
	if err != nil {
		return err
	}
	eventMessageDef := eventMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(eventMessageDef); err != nil {
		return err
	}
	eventStartMessageData, err := eventMessageDef.ConstructData(FIT_EVENT_TIMER, act.startTime.Unix(), FIT_EVENT_TIMER_TRIGGER_MANUAL, FIT_EVENT_TYPE_START)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(eventStartMessageData); err != nil {
		return err
	}

	// add normal record definition message
	recordMessage, err := getRecordMessageDefinition()
	if err != nil {
		return err
	}
	recordMessageDef := recordMessage.ConstructLocalMessage(1)
	if err := file.AddMessage(recordMessageDef); err != nil {
		return err
	}

	// add distance record definition message
	recordDistanceMessage, err := getRecordDistanceMessageDefinition()
	if err != nil {
		return err
	}
	recordDistanceMessageDef := recordDistanceMessage.ConstructLocalMessage(2)
	if err := file.AddMessage(recordDistanceMessageDef); err != nil {
		return err
	}

	// add all records to file
	for _, record := range act.activityRecords {
		// add record normal data
		recordMessageData, err := recordMessageDef.ConstructData(fitEncodeCoordinate(record.lat), fitEncodeCoordinate(record.lon), fitEncodeAltitudeM(record.altitude), fitEncodeSpeedKmH(record.speed), STRAVA_NOICE_GPS_ACCURACY, record.timestamp.Unix())
		if err != nil {
			return err
		}
		if err := file.AddMessage(recordMessageData); err != nil {
			return err
		}

		// add record distance data
		recordDistanceMessageData, err := recordDistanceMessageDef.ConstructData(record.timestamp.Unix(), fitEncodeDistanceKm(record.distance))
		if err != nil {
			return err
		}
		if err := file.AddMessage(recordDistanceMessageData); err != nil {
			return err
		}
	}

	// add stop event to file
	eventStopMessageData, err := eventMessageDef.ConstructData(FIT_EVENT_TIMER, act.startTime.Add(act.totalActiveTime).Unix(), FIT_EVENT_TIMER_TRIGGER_MANUAL, FIT_EVENT_TYPE_STOP)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(eventStopMessageData); err != nil {
		return err
	}

	return nil
}

func (act *stravaActivity) writeFooter(file *fit.FitFile) error {
	// add device battery info message on end of activity
	deviceInfoBatteryMessage, err := getDeviceInfoBatteryMessageDefinition()
	if err != nil {
		return err
	}
	deviceInfoBatteryMessageDef := deviceInfoBatteryMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(deviceInfoBatteryMessageDef); err != nil {
		return err
	}
	deviceInfoBatteryMessageData, err := deviceInfoBatteryMessageDef.ConstructData(FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID, act.startTime.Add(act.totalActiveTime).Unix(), FIT_BATTERY_STATUS_GOOD)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(deviceInfoBatteryMessageData); err != nil {
		return err
	}

	return nil
}
