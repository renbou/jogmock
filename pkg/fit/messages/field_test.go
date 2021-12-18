// Copyright 2021 Artem Mikheev

package messages

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/renbou/strava-keker/pkg/fit/types"
)

func encodeAndValidateFieldDef(t *testing.T, fieldDef *FieldDefinition,
	endianness types.Endianness, valid bool, expected []byte,
) {
	var (
		actual []byte
		err    error
	)

	actual, err = EncodeFitFieldDefinition(nil, endianness, fieldDef)
	if valid && err != nil {
		t.Fatalf("Error not expected for encoding %+v with endianness %s", fieldDef, endianness)
	}

	if !valid {
		if err == nil {
			t.Fatalf("Expected error for encoding %+v with endianness %s", fieldDef, endianness)
		} else {
			return
		}
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("Expected %s but got %s for encoding %+v with endianness %s",
			hex.EncodeToString(expected), hex.EncodeToString(actual), fieldDef, endianness)
	}
}

func encodeAndValidateField(t *testing.T, fieldDef *FieldDefinition,
	endianness types.Endianness, valid bool, expected []byte, value interface{},
) {
	var (
		actual []byte
		err    error
	)

	actual, err = EncodeFitField(nil, endianness, fieldDef, value)
	if valid && err != nil {
		t.Fatalf("Error not expected for encoding %+v using fieldDef %+v with endianness %s",
			value, fieldDef, endianness)
	}

	if !valid {
		if err == nil {
			t.Fatalf("Expected error for encoding %+v using fieldDef %+v with endianness %s",
				value, fieldDef, endianness)
		} else {
			return
		}
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("Expected %s but got %s for encoding %+v using fieldDef %+v with endianness %s",
			hex.EncodeToString(expected), hex.EncodeToString(actual), value, fieldDef, endianness)
	}
}

func TestFitFieldDefinitionEncoding(t *testing.T) {
	encodeAndValidateFieldDef(t, &FieldDefinition{
		DefNum:   1,
		Size:     types.FIT_TYPE_UINT16_SIZE,
		BaseType: types.FIT_TYPE_UINT16,
	}, types.BigEndian, true, []byte{1, 2, 0x84})

	encodeAndValidateFieldDef(t, &FieldDefinition{
		DefNum:   1,
		Size:     types.FIT_TYPE_UINT16_SIZE,
		BaseType: types.FIT_TYPE_UINT16,
	}, types.LittleEndian, true, []byte{1, 2, 0x84})

	encodeAndValidateFieldDef(t, &FieldDefinition{
		DefNum:   0,
		Size:     types.FIT_TYPE_ENUM_SIZE,
		BaseType: types.FIT_TYPE_ENUM,
	}, types.BigEndian, true, []byte{0, 1, 0})

	encodeAndValidateFieldDef(t, &FieldDefinition{
		DefNum:   4,
		Size:     types.FIT_TYPE_UINT32_SIZE,
		BaseType: types.FIT_TYPE_UINT32,
	}, types.BigEndian, true, []byte{4, 4, 0x86})

	encodeAndValidateFieldDef(t, &FieldDefinition{
		DefNum:   3,
		Size:     17,
		BaseType: types.FIT_TYPE_STRING,
	}, types.BigEndian, true, []byte{3, 17, 7})

	encodeAndValidateFieldDef(t, &FieldDefinition{
		DefNum:   1,
		Size:     0,
		BaseType: types.FIT_TYPE_STRING,
	}, types.BigEndian, false, nil)
}

func TestFitFieldEncoding(t *testing.T) {
	encodeAndValidateField(t, &FieldDefinition{
		DefNum:   1,
		Size:     types.FIT_TYPE_UINT16_SIZE,
		BaseType: types.FIT_TYPE_UINT16,
	}, types.BigEndian, true, []byte{0x01, 0x09}, types.FitUint16(265))

	encodeAndValidateField(t, &FieldDefinition{
		DefNum:   3,
		Size:     17,
		BaseType: types.FIT_TYPE_STRING,
	}, types.BigEndian, true,
		[]byte{0x6c, 0x69, 0x76, 0x65, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x00},
		types.FitString("live_activity_id"))

	encodeAndValidateField(t, &FieldDefinition{
		DefNum:   3,
		Size:     4,
		BaseType: types.FIT_TYPE_STRING,
	}, types.BigEndian, false,
		[]byte{0x62, 0x65, 0x62, 0x72, 0x61, 0x00},
		types.FitString("bebra"))

	encodeAndValidateField(t, &FieldDefinition{
		DefNum:   3,
		Size:     7,
		BaseType: types.FIT_TYPE_STRING,
	}, types.BigEndian, true,
		[]byte{0x61, 0x62, 0x6f, 0x62, 0x61, 0x00, 0x00},
		types.FitString("aboba"))

	encodeAndValidateField(t, &FieldDefinition{
		DefNum:   1,
		Size:     types.FIT_TYPE_SINT32_SIZE,
		BaseType: types.FIT_TYPE_SINT32,
	}, types.BigEndian, true,
		[]byte{0x15, 0x8c, 0xcf, 0xac},
		types.FitSint32(361549740))

	encodeAndValidateField(t, &FieldDefinition{
		DefNum:   0,
		Size:     types.FIT_TYPE_SINT32_SIZE,
		BaseType: types.FIT_TYPE_SINT32,
	}, types.LittleEndian, true,
		[]byte{0x6b, 0xc1, 0x9c, 0x2a},
		types.FitSint32(714916203))
}
