// Copyright 2021 Artem Mikheev

package stravafit

import (
	"fmt"
	"math"
	"time"

	"github.com/renbou/jogmock/activities"
	"github.com/renbou/jogmock/fit-encoder/fit"
	"github.com/renbou/jogmock/fit-encoder/fit/types"
)

func ActivityTypeToString(activityType activities.ActivityType) string {
	if activityType == activities.RunActivity {
		return "Run"
	} else if activityType == activities.RideActivity {
		return "Ride"
	}
	panic(fmt.Sprintf("unknown activity type: %v", activityType))
}

type StravaActivity struct {
	AppVersion         uint32
	MobileAppVersion   string
	DeviceManufacturer string
	DeviceModel        string
	DeviceOsVersion    string
	Activity           *activities.Activity
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

func fitActivitySport(act *StravaActivity) uint8 {
	activityType := act.Activity.Type()
	if activityType == activities.RunActivity {
		return 1
	} else if activityType == activities.RideActivity {
		return 2
	}
	panic(fmt.Sprintf("unknown activity type: %v", activityType))
}

func fitEncodeDistanceKm(distance float64) int {
	return int(math.Round(distance * 100000.0))
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

func fitEncodeTimestamp(t time.Time) int64 {
	// fit time is offset by 631065600 seconds
	return t.Add(-time.Second * 631065600).Unix()
}

func fitEncodeDuration(duration time.Duration) int {
	return int(duration.Milliseconds())
}

// BuildFitFile creates a fit file based on the filled activity
func (act *StravaActivity) BuildFitFile() (*fit.FitFile, error) {
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

func (act *StravaActivity) writeHeader(file *fit.FitFile) error {
	// add simple file id message
	fileIdMessage, err := getFileIdMessageDefinition()
	if err != nil {
		return err
	}
	fileIdMessageDef := fileIdMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(fileIdMessageDef); err != nil {
		return err
	}
	fileIdMessageData, err := fileIdMessageDef.ConstructData(
		FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID, FIT_FILEID_TYPE_ACTIVITY,
		fitEncodeTimestamp(act.Activity.Start()))
	if err != nil {
		return nil
	}
	if err := file.AddMessage(fileIdMessageData); err != nil {
		return err
	}

	// add developer data id message and field definition messages
	dev := &devData{
		index:      0,
		appVersion: act.AppVersion,
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
	deviceInfoMessageData, err := deviceInfoMessageDef.ConstructData(
		FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID,
		act.MobileAppVersion, act.DeviceManufacturer, act.DeviceModel, act.DeviceOsVersion)
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
	activityMessageData, err := activityMessageDef.ConstructData(
		FIT_ACTIVITY_ONE_SESSION, fitEncodeDuration(act.Activity.TotalDuration()),
		FIT_ACTIVITY_EVENT_ACTIVITY, FIT_EVENT_TYPE_STOP)
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
	sessionMessageData, err := sessionMessageDef.ConstructData(
		fitEncodeDuration(act.Activity.TotalDuration()), fitEncodeTimestamp(act.Activity.Start()),
		fitEncodeDuration(act.Activity.TotalDuration()), fitEncodeDistanceKm(act.Activity.TotalDistance()),
		fitActivitySport(act),
		FIT_SESSION_EVENT_SESSION, FIT_EVENT_TYPE_STOP, FIT_SESSION_ONE_LAP, STRAVA_FIRST_LIVE_ACTIVITY_ID,
		ActivityTypeToString(act.Activity.Type()), STRAVA_AUTOPAUSE_DISABLED)
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
	lapMessageData, err := lapMessageDef.ConstructData(
		fitEncodeDuration(act.Activity.TotalDuration()), fitEncodeTimestamp(act.Activity.Start()),
		fitEncodeDuration(act.Activity.TotalDuration()), fitEncodeDistanceKm(act.Activity.TotalDistance()),
		FIT_LAP_EVENT_LAP, FIT_EVENT_TYPE_STOP, fitActivitySport(act), FIT_LAP_TRIGGER_SESSION_END)
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
	deviceInfoBatteryMessageData, err := deviceInfoBatteryMessageDef.ConstructData(
		FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID, fitEncodeTimestamp(act.Activity.Start()), FIT_BATTERY_STATUS_GOOD)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(deviceInfoBatteryMessageData); err != nil {
		return err
	}

	return nil
}

func (act *StravaActivity) writeBody(file *fit.FitFile) error {
	// add event message on start of activity
	eventMessage, err := getEventMessageDefinition()
	if err != nil {
		return err
	}
	eventMessageDef := eventMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(eventMessageDef); err != nil {
		return err
	}
	eventStartMessageData, err := eventMessageDef.ConstructData(
		FIT_EVENT_TIMER, fitEncodeTimestamp(act.Activity.Start()), FIT_EVENT_TIMER_TRIGGER_MANUAL, FIT_EVENT_TYPE_START)
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
	for _, record := range act.Activity.Records() {
		// add record normal data
		recordMessageData, err := recordMessageDef.ConstructData(
			fitEncodeCoordinate(record.Lat), fitEncodeCoordinate(record.Lon),
			fitEncodeAltitudeM(record.Altitude), fitEncodeSpeedKmH(record.Speed),
			STRAVA_NOICE_GPS_ACCURACY, fitEncodeTimestamp(record.Timestamp))
		if err != nil {
			return err
		}
		if err := file.AddMessage(recordMessageData); err != nil {
			return err
		}

		// add record distance data
		recordDistanceMessageData, err := recordDistanceMessageDef.ConstructData(
			fitEncodeTimestamp(record.Timestamp), fitEncodeDistanceKm(record.Distance))
		if err != nil {
			return err
		}
		if err := file.AddMessage(recordDistanceMessageData); err != nil {
			return err
		}
	}

	// add stop event to file
	eventStopMessageData, err := eventMessageDef.ConstructData(
		FIT_EVENT_TIMER,
		fitEncodeTimestamp(act.Activity.Start().Add(act.Activity.TotalDuration())),
		FIT_EVENT_TIMER_TRIGGER_MANUAL, FIT_EVENT_TYPE_STOP)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(eventStopMessageData); err != nil {
		return err
	}

	return nil
}

func (act *StravaActivity) writeFooter(file *fit.FitFile) error {
	// add device battery info message on end of activity
	deviceInfoBatteryMessage, err := getDeviceInfoBatteryMessageDefinition()
	if err != nil {
		return err
	}
	deviceInfoBatteryMessageDef := deviceInfoBatteryMessage.ConstructLocalMessage(0)
	if err := file.AddMessage(deviceInfoBatteryMessageDef); err != nil {
		return err
	}
	deviceInfoBatteryMessageData, err := deviceInfoBatteryMessageDef.ConstructData(
		FIT_MANUFACTURER_STRAVA, FIT_PRODUCT_STRAVA_ANDROID,
		fitEncodeTimestamp(act.Activity.Start().Add(act.Activity.TotalDuration())),
		FIT_BATTERY_STATUS_GOOD)
	if err != nil {
		return nil
	}
	if err := file.AddMessage(deviceInfoBatteryMessageData); err != nil {
		return err
	}

	return nil
}
