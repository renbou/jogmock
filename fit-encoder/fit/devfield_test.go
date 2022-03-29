// Copyright 2021 Artem Mikheev

package fit

import (
	"testing"

	"github.com/renbou/jogmock/strava-mock/pkg/encoding"
	"github.com/renbou/jogmock/strava-mock/pkg/fit/types"
	"github.com/stretchr/testify/assert"
)

func TestFitDevFieldDefinitionEncoding(t *testing.T) {
	a := assert.New(t)

	encodeAndValidate(a, &DevFieldDefinition{
		Field: &FieldDescriptionStub{
			DevDataIndex: 0,
			DefNum:       0,
			BaseType:     types.FIT_TYPE_UINT64,
		},
		Size: types.FIT_TYPE_UINT64_SIZE,
		DevId: &DeveloperDataIdStub{
			DevDataIndex: 0,
		},
	}, encoding.BigEndian, true, []byte{0, 8, 0})

	encodeAndValidate(a, &DevFieldDefinition{
		Field: &FieldDescriptionStub{
			DevDataIndex: 1,
			DefNum:       2,
			BaseType:     types.FIT_TYPE_UINT16,
		},
		Size: types.FIT_TYPE_UINT16_SIZE,
		DevId: &DeveloperDataIdStub{
			DevDataIndex: 1,
		},
	}, encoding.LittleEndian, true, []byte{2, 2, 1})

	encodeAndValidate(a, &DevFieldDefinition{
		Field: &FieldDescriptionStub{
			DevDataIndex: 1,
			DefNum:       5,
			BaseType:     types.FIT_TYPE_STRING,
		},
		Size: 0,
		DevId: &DeveloperDataIdStub{
			DevDataIndex: 1,
		},
	}, encoding.LittleEndian, false, nil)

	encodeAndValidate(a, &DevFieldDefinition{
		Field: &FieldDescriptionStub{
			DevDataIndex: 3,
			DefNum:       5,
			BaseType:     types.FIT_TYPE_SINT32,
		},
		Size: 3,
		DevId: &DeveloperDataIdStub{
			DevDataIndex: 3,
		},
	}, encoding.BigEndian, false, nil)

	encodeAndValidate(a, &DevFieldDefinition{
		Field: &FieldDescriptionStub{
			DevDataIndex: 3,
			DefNum:       4,
			BaseType:     types.FIT_TYPE_SINT32,
		},
		Size: types.FIT_TYPE_SINT32_SIZE,
		DevId: &DeveloperDataIdStub{
			DevDataIndex: 2,
		},
	}, encoding.BigEndian, false, nil)
}

func TestFitDevFieldEncoding(t *testing.T) {
	a := assert.New(t)

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 0,
				DefNum:       0,
				BaseType:     types.FIT_TYPE_UINT64,
			},
			Size: types.FIT_TYPE_UINT64_SIZE,
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 0,
			},
		},
		Value: types.FitUint64(0),
	}, encoding.BigEndian, true, []byte{0, 0, 0, 0, 0, 0, 0, 0})

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 0,
				DefNum:       0,
				BaseType:     types.FIT_TYPE_SINT32,
			},
			Size: types.FIT_TYPE_SINT32_SIZE,
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 0,
			},
		},
		Value: types.FitSint32(361549740),
	}, encoding.BigEndian, true, []byte{0x15, 0x8c, 0xcf, 0xac})

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 1,
				DefNum:       2,
				BaseType:     types.FIT_TYPE_UINT16,
			},
			Size: types.FIT_TYPE_UINT16_SIZE,
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 1,
			},
		},
		Value: types.FitUint16(64),
	}, encoding.LittleEndian, true, []byte{0x40, 0x00})

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 1,
				DefNum:       5,
				BaseType:     types.FIT_TYPE_STRING,
			},
			Size: 0,
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 1,
			},
		},
		Value: types.FitString("aboba"),
	}, encoding.LittleEndian, false, nil)

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 3,
				DefNum:       5,
				BaseType:     types.FIT_TYPE_SINT32,
			},
			Size: 3,
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 3,
			},
		},
		Value: types.FitSint32(234),
	}, encoding.BigEndian, false, nil)

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 3,
				DefNum:       4,
				BaseType:     types.FIT_TYPE_SINT32,
			},
			Size: types.FIT_TYPE_SINT32_SIZE,
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 2,
			},
		},
		Value: types.FitSint32(-1234),
	}, encoding.BigEndian, false, nil)

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 1,
				DefNum:       2,
				BaseType:     types.FIT_TYPE_STRING,
			},
			Size: types.FitUint8(len("not a fit string") + 1),
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 1,
			},
		},
		Value: "not a fit string",
	}, encoding.LittleEndian, false, nil)

	encodeAndValidate(a, &DevField{
		Def: &DevFieldDefinition{
			Field: &FieldDescriptionStub{
				DevDataIndex: 1,
				DefNum:       2,
				BaseType:     types.FIT_TYPE_STRING,
			},
			Size: 3,
			DevId: &DeveloperDataIdStub{
				DevDataIndex: 1,
			},
		},
		Value: types.FitString("fit string with invalid size"),
	}, encoding.LittleEndian, false, nil)
}
