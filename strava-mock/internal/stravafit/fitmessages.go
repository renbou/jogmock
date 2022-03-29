// Copyright 2021 Artem Mikheev

package stravafit

import (
	"github.com/renbou/jogmock/fit-encoder/fit"
	"github.com/renbou/jogmock/fit-encoder/fit/types"
)

// getters for all needed strava fit message definitions

func getFileIdMessageDefinition() (*fit.DefinitionMessage, error) {
	return &fit.DefinitionMessage{
		GlobalMsgNum: 0,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   2,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   4,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
		},
		DevFieldDefs: nil,
	}, nil
}

func getDeviceInfoMessageDefinition(act *StravaActivity, dev *devData) (*fit.DefinitionMessage, error) {
	mobileAppVersionDevField, err := dev.getFieldDefinition("mobile_app_version", len(act.MobileAppVersion))
	if err != nil {
		return nil, err
	}
	deviceManufacturerDevField, err := dev.getFieldDefinition("device_manufacturer", len(act.DeviceManufacturer))
	if err != nil {
		return nil, err
	}
	deviceModelDevField, err := dev.getFieldDefinition("device_model", len(act.DeviceModel))
	if err != nil {
		return nil, err
	}
	deviceOsVersion, err := dev.getFieldDefinition("device_os_version", len(act.DeviceOsVersion))
	if err != nil {
		return nil, err
	}
	return &fit.DefinitionMessage{
		GlobalMsgNum: 23,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   2,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   4,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
		},
		DevFieldDefs: []*fit.DevFieldDefinition{
			mobileAppVersionDevField,
			deviceManufacturerDevField,
			deviceModelDevField,
			deviceOsVersion,
		},
	}, nil
}

func getActivityMessageDefinition() (*fit.DefinitionMessage, error) {
	return &fit.DefinitionMessage{
		GlobalMsgNum: 34,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   3,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   4,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
		},
		DevFieldDefs: nil,
	}, nil
}

func getSessionMessageDefinition(act *StravaActivity, dev *devData) (*fit.DefinitionMessage, error) {
	liveActivityIdDevField, err := dev.getFieldDefinition("live_activity_id")
	if err != nil {
		return nil, err
	}
	activityTypeDevField, err := dev.getFieldDefinition("activity_type", len(ActivityTypeToString(act.Activity.Type())))
	if err != nil {
		return nil, err
	}
	autopauseEnabledDevField, err := dev.getFieldDefinition("autopause_enabled")
	if err != nil {
		return nil, err
	}
	return &fit.DefinitionMessage{
		GlobalMsgNum: 18,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   7,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   2,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   8,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   9,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   5,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   26,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
		},
		DevFieldDefs: []*fit.DevFieldDefinition{
			liveActivityIdDevField,
			activityTypeDevField,
			autopauseEnabledDevField,
		},
	}, nil
}

func getLapMessageDefinition() (*fit.DefinitionMessage, error) {
	return &fit.DefinitionMessage{
		GlobalMsgNum: 19,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   7,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   2,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   8,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   9,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   25,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   24,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
		},
		DevFieldDefs: nil,
	}, nil
}

func getDeviceInfoBatteryMessageDefinition() (*fit.DefinitionMessage, error) {
	return &fit.DefinitionMessage{
		GlobalMsgNum: 23,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   2,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   4,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   253,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   11,
				Size:     types.FIT_TYPE_UINT8_SIZE,
				BaseType: types.FIT_TYPE_UINT8,
			},
		},
		DevFieldDefs: nil,
	}, nil
}

func getEventMessageDefinition() (*fit.DefinitionMessage, error) {
	return &fit.DefinitionMessage{
		GlobalMsgNum: 21,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   253,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   3,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
		},
		DevFieldDefs: nil,
	}, nil
}

func getRecordMessageDefinition() (*fit.DefinitionMessage, error) {
	return &fit.DefinitionMessage{
		GlobalMsgNum: 20,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_SINT32_SIZE,
				BaseType: types.FIT_TYPE_SINT32,
			},
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_SINT32_SIZE,
				BaseType: types.FIT_TYPE_SINT32,
			},
			{
				DefNum:   78,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   6,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   31,
				Size:     types.FIT_TYPE_UINT8_SIZE,
				BaseType: types.FIT_TYPE_UINT8,
			},
			{
				DefNum:   253,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
		},
		DevFieldDefs: nil,
	}, nil
}

func getRecordDistanceMessageDefinition() (*fit.DefinitionMessage, error) {
	return &fit.DefinitionMessage{
		GlobalMsgNum: 20,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   253,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
			{
				DefNum:   5,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
		},
		DevFieldDefs: nil,
	}, nil
}
