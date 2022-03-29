// Copyright 2021 Artem Mikheev

package fit

import (
	"bytes"
	"testing"

	"github.com/renbou/jogmock/fit-encoder/encoding"
	"github.com/renbou/jogmock/fit-encoder/fit/types"
	"github.com/stretchr/testify/assert"
)

func encodeAndValidate(a *assert.Assertions, value encoding.EndianEncoder,
	endianness encoding.Endianness, valid bool, expected []byte,
) {
	var err error

	buffer := new(bytes.Buffer)
	encoder := encoding.NewEncoder(buffer, endianness)

	err = encoder.Encode(value)
	if valid && !a.NoErrorf(err, "no error expected for %+v with %s", value, encoder) {
		return
	} else if !valid {
		a.Error(err, "expected error for %+v with %s", value, encoder)
		return
	}

	a.Equal(expected, buffer.Bytes(), "invalid encoding of %+v with %s", value, encoder)
}

func TestFitFieldDefinitionEncoding(t *testing.T) {
	a := assert.New(t)

	encodeAndValidate(a, &FieldDefinition{
		DefNum:   1,
		Size:     types.FIT_TYPE_UINT16_SIZE,
		BaseType: types.FIT_TYPE_UINT16,
	}, encoding.BigEndian, true, []byte{1, 2, 0x84})

	encodeAndValidate(a, &FieldDefinition{
		DefNum:   1,
		Size:     types.FIT_TYPE_UINT16_SIZE,
		BaseType: types.FIT_TYPE_UINT16,
	}, encoding.LittleEndian, true, []byte{1, 2, 0x84})

	encodeAndValidate(a, &FieldDefinition{
		DefNum:   0,
		Size:     types.FIT_TYPE_ENUM_SIZE,
		BaseType: types.FIT_TYPE_ENUM,
	}, encoding.BigEndian, true, []byte{0, 1, 0})

	encodeAndValidate(a, &FieldDefinition{
		DefNum:   4,
		Size:     types.FIT_TYPE_UINT32_SIZE,
		BaseType: types.FIT_TYPE_UINT32,
	}, encoding.BigEndian, true, []byte{4, 4, 0x86})

	encodeAndValidate(a, &FieldDefinition{
		DefNum:   3,
		Size:     17,
		BaseType: types.FIT_TYPE_STRING,
	}, encoding.BigEndian, true, []byte{3, 17, 7})

	encodeAndValidate(a, &FieldDefinition{
		DefNum:   1,
		Size:     0,
		BaseType: types.FIT_TYPE_STRING,
	}, encoding.BigEndian, false, nil)

	encodeAndValidate(a, &FieldDefinition{
		DefNum:   1,
		Size:     0,
		BaseType: types.FIT_TYPE_INVALID,
	}, encoding.BigEndian, false, nil)
}

func TestFitFieldEncoding(t *testing.T) {
	a := assert.New(t)

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   1,
			Size:     types.FIT_TYPE_UINT16_SIZE,
			BaseType: types.FIT_TYPE_UINT16,
		},
		Value: types.FitUint16(265),
	}, encoding.BigEndian, true, []byte{0x01, 0x09})

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   3,
			Size:     17,
			BaseType: types.FIT_TYPE_STRING,
		},
		Value: types.FitString("live_activity_id"),
	}, encoding.BigEndian, true,
		[]byte{0x6c, 0x69, 0x76, 0x65, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x00})

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   3,
			Size:     4,
			BaseType: types.FIT_TYPE_STRING,
		},
		Value: types.FitString("bebra"),
	}, encoding.BigEndian, false,
		[]byte{0x62, 0x65, 0x62, 0x72, 0x61, 0x00})

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   3,
			Size:     7,
			BaseType: types.FIT_TYPE_STRING,
		},
		Value: types.FitString("aboba"),
	}, encoding.BigEndian, true,
		[]byte{0x61, 0x62, 0x6f, 0x62, 0x61, 0x00, 0x00})

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   1,
			Size:     types.FIT_TYPE_SINT32_SIZE,
			BaseType: types.FIT_TYPE_SINT32,
		},
		Value: types.FitSint32(361549740),
	}, encoding.BigEndian, true,
		[]byte{0x15, 0x8c, 0xcf, 0xac})

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   0,
			Size:     types.FIT_TYPE_SINT32_SIZE,
			BaseType: types.FIT_TYPE_SINT32,
		},
		Value: types.FitSint32(714916203),
	}, encoding.LittleEndian, true,
		[]byte{0x6b, 0xc1, 0x9c, 0x2a})

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   3,
			Size:     4,
			BaseType: types.FIT_TYPE_INVALID,
		},
		Value: "invalid type",
	}, encoding.LittleEndian, false, nil)
}

func TestFitFieldErrors(t *testing.T) {
	a := assert.New(t)

	encodeAndValidate(a, &Field{
		Def: &FieldDefinition{
			DefNum:   0,
			Size:     types.FIT_TYPE_SINT32_SIZE,
			BaseType: types.FIT_TYPE_SINT32,
		},
		Value: types.FitString("aboba"),
	}, encoding.LittleEndian, false, nil)
}
