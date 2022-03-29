// Copyright 2021 Artem Mikheev

package fit

import (
	"encoding/hex"
	"testing"

	"github.com/renbou/jogmock/fit-encoder/encoding"
	"github.com/renbou/jogmock/fit-encoder/fit/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageEncoding(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

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
	encodeAndValidate(a, localFileIdMsgDef, encoding.BigEndian, true, []byte{
		0x40, 0x00, 0x01, 0x00, 0x00, 0x04, 0x01, 0x02, 0x84, 0x02, 0x02, 0x84, 0x00, 0x01, 0x00, 0x04, 0x04, 0x86,
	})

	localFileIdMsg, err := localFileIdMsgDef.ConstructData(265, 102, 4, 1007562558)
	r.NoError(err, "construction of valid file id")
	encodeAndValidate(a, localFileIdMsg, encoding.BigEndian, true, []byte{
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
	encodeAndValidate(a, localDeveloperDataIdMsgDef, encoding.BigEndian, true, []byte{
		0x40, 0x00, 0x01, 0x00, 0xCF, 0x02, 0x03, 0x01, 0x02, 0x04, 0x04, 0x86,
	})

	localDeveloperDataIdMsg, err := localDeveloperDataIdMsgDef.ConstructData(0, 1221988)
	r.NoError(err, "construction of valid developer data id")
	encodeAndValidate(a, localDeveloperDataIdMsg, encoding.BigEndian, true, []byte{
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
	encodeAndValidate(a, localDeviceInfoMsgDef, encoding.BigEndian, true, []byte{
		0x60, 0x00, 0x01, 0x00, 0x17, 0x02, 0x02, 0x02, 0x84, 0x04, 0x02, 0x84, 0x04,
		0x03, 0x11, 0x00, 0x06, 0x07, 0x00, 0x04, 0x11, 0x00, 0x05, 0x03, 0x00,
	})

	localDeviceInfoMsg, err := localDeviceInfoMsgDef.ConstructData(
		265, 102, "230.10 (1221988)", "Xiaomi", "Redmi Note 9 Pro", "10",
	)
	r.NoError(err, "construction of valid local device id")
	localDeviceInfoMsgBytes, _ := hex.DecodeString(
		"00010900663233302E313020283132323139383829005869616F6D69005265646D69204E6F746520392050726F00313000",
	)
	encodeAndValidate(a, localDeviceInfoMsg, encoding.BigEndian, true, localDeviceInfoMsgBytes)

	// Error tests
	a.ErrorIs(localDeviceInfoMsgDef.Encode(nil, encoding.Endianness(123)),
		encoding.ErrUnknownEndianness,
	)
	a.ErrorIs(localDeviceInfoMsg.Encode(nil, encoding.Endianness(123)),
		encoding.ErrUnknownEndianness,
	)

	_, err = localDeviceInfoMsgDef.ConstructData(types.FitUint16(265))
	r.Error(err, "construction of invalid local device info (invalid length)")

	_, err = localDeviceInfoMsgDef.ConstructData(
		types.FitUint16(265), "102",
		types.FitString("230.10 (1221988)"), types.FitString("Xiaomi"),
		types.FitString("Redmi Note 9 Pro"), "10",
	)
	r.Error(err, "construction of invalid local device info (invalid field types)")

	_, err = localDeviceInfoMsgDef.ConstructData(
		types.FitUint16(265), types.FitUint16(102),
		types.FitString("230.10 (1221988)"), types.FitString("Xiaomi"),
		types.FitString("Redmi Note 9 Pro"), types.FitUint16(1),
	)
	r.Error(err, "construction of invalid local device info (invalid def field types)")

	fakeMsgDef := &DefinitionMessage{
		GlobalMsgNum: 123,
		FieldDefs: []*FieldDefinition{
			{
				DefNum:   3,
				Size:     4,
				BaseType: types.FIT_TYPE_INVALID,
			},
		},
		DevFieldDefs: nil,
	}
	localFakeMsgDef := fakeMsgDef.ConstructLocalMessage(10)
	_, err = localFakeMsgDef.ConstructData("invalid fit type")
	a.Error(err)
}

func TestInvalidMessageHeader(t *testing.T) {
	a := assert.New(t)

	a.ErrorIs(encodeMessageHeader(nil, defMsgType, true, 20), ErrInvalidLocalMsgType)
	a.ErrorIs(encodeMessageHeader(nil, 2, true, 4), ErrInvalidMsgType)
	a.ErrorIs(encodeMessageHeader(nil, dataMsgType, true, 4), ErrInvalidMsgSpecific)
}
