// Copyright 2021 Artem Mikheev

package fit

import (
	"encoding/hex"
	"testing"

	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit/types"
)

func TestMessageEncoding(t *testing.T) {
	fileIdMsgDef := &DefinitionMessage{
		GlobalMsgNum: 0,
		FieldDefs: []*FieldDefinition{
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
	}

	localFileIdMsgDef := fileIdMsgDef.ConstructLocalMessage(0)
	encodeAndValidate(t, localFileIdMsgDef, encoding.BigEndian, true, []byte{
		0x40, 0x00, 0x01, 0x00, 0x00, 0x04, 0x01, 0x02, 0x84, 0x02, 0x02, 0x84, 0x00, 0x01, 0x00, 0x04, 0x04, 0x86,
	})

	localFileIdMsg, err := localFileIdMsgDef.ConstructData(
		types.FitUint16(265), types.FitUint16(102), types.FitEnum(4), types.FitUint32(1007562558),
	)
	if err != nil {
		t.Fatalf("Unexpected error during data construction for file id: %v", err)
	}
	encodeAndValidate(t, localFileIdMsg, encoding.BigEndian, true, []byte{
		0x00, 0x01, 0x09, 0x00, 0x66, 0x04, 0x3C, 0x0E, 0x2F, 0x3E,
	})

	developerDataIdMsgDef := &DefinitionMessage{
		GlobalMsgNum: 207,
		FieldDefs: []*FieldDefinition{
			{
				DefNum:   3,
				Size:     types.FIT_TYPE_UINT8_SIZE,
				BaseType: types.FIT_TYPE_UINT8,
			},
			{
				DefNum:   4,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
		},
		DevFieldDefs: nil,
	}

	localDeveloperDataIdMsgDef := developerDataIdMsgDef.ConstructLocalMessage(0)
	encodeAndValidate(t, localDeveloperDataIdMsgDef, encoding.BigEndian, true, []byte{
		0x40, 0x00, 0x01, 0x00, 0xCF, 0x02, 0x03, 0x01, 0x02, 0x04, 0x04, 0x86,
	})

	localDeveloperDataIdMsg, err := localDeveloperDataIdMsgDef.ConstructData(
		types.FitUint8(0), types.FitUint32(1221988),
	)
	if err != nil {
		t.Fatalf("Unexpected error during data construction for developer data id: %v", err)
	}
	encodeAndValidate(t, localDeveloperDataIdMsg, encoding.BigEndian, true, []byte{
		0x00, 0x00, 0x00, 0x12, 0xA5, 0x64,
	})

	deviceInfoMsgDef := &DefinitionMessage{
		GlobalMsgNum: 23,
		FieldDefs: []*FieldDefinition{
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
		DevFieldDefs: []*DevFieldDefinition{
			{
				Field: &FieldDescriptionStub{
					DevDataIndex: 0,
					DefNum:       3,
					BaseType:     types.FIT_TYPE_STRING,
				},
				Size: 17,
				DevId: &DeveloperDataIdStub{
					DevDataIndex: 0,
				},
			},
			{
				Field: &FieldDescriptionStub{
					DevDataIndex: 0,
					DefNum:       6,
					BaseType:     types.FIT_TYPE_STRING,
				},
				Size: 7,
				DevId: &DeveloperDataIdStub{
					DevDataIndex: 0,
				},
			},
			{
				Field: &FieldDescriptionStub{
					DevDataIndex: 0,
					DefNum:       4,
					BaseType:     types.FIT_TYPE_STRING,
				},
				Size: 17,
				DevId: &DeveloperDataIdStub{
					DevDataIndex: 0,
				},
			},
			{
				Field: &FieldDescriptionStub{
					DevDataIndex: 0,
					DefNum:       5,
					BaseType:     types.FIT_TYPE_STRING,
				},
				Size: 3,
				DevId: &DeveloperDataIdStub{
					DevDataIndex: 0,
				},
			},
		},
	}

	localDeviceInfoMsgDef := deviceInfoMsgDef.ConstructLocalMessage(0)
	encodeAndValidate(t, localDeviceInfoMsgDef, encoding.BigEndian, true, []byte{
		0x60, 0x00, 0x01, 0x00, 0x17, 0x02, 0x02, 0x02, 0x84, 0x04, 0x02, 0x84, 0x04,
		0x03, 0x11, 0x00, 0x06, 0x07, 0x00, 0x04, 0x11, 0x00, 0x05, 0x03, 0x00,
	})

	localDeviceInfoMsg, err := localDeviceInfoMsgDef.ConstructData(
		types.FitUint16(265), types.FitUint16(102),
		types.FitString("230.10 (1221988)"), types.FitString("Xiaomi"),
		types.FitString("Redmi Note 9 Pro"), types.FitString("10"),
	)
	if err != nil {
		t.Fatalf("Unexpected error during data construction for developer data id: %v", err)
	}
	localDeviceInfoMsgBytes, _ := hex.DecodeString(
		"00010900663233302E313020283132323139383829005869616F6D69005265646D69204E6F746520392050726F00313000",
	)
	encodeAndValidate(t, localDeviceInfoMsg, encoding.BigEndian, true, localDeviceInfoMsgBytes)
}
