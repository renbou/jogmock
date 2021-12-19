// Copyright 2021 Artem Mikheev

package fit

import (
	"errors"
	"io"
	"reflect"

	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit/types"
	"github.com/renbou/strava-keker/pkg/io/maybe"
)

// DeveloperDataIdStub is a struct representing only the "required" part
// of a DeveloperDataId message - the developer data index
type DeveloperDataIdStub struct {
	DevDataIndex types.FitUint8
}

// FieldDescriptionStub is a struct representing only the "required" part
// of a FieldDescription message - the developer data index and field def number
type FieldDescriptionStub struct {
	DevDataIndex types.FitUint8
	DefNum       types.FitUint8
	BaseType     types.FitBaseType
}

// DevFieldDefinition struct encapsulates the developer field definition
// logic in a way that makes usages explicitly specify the field description
// and developer data id
type DevFieldDefinition struct {
	Field *FieldDescriptionStub
	Size  types.FitUint8
	DevId *DeveloperDataIdStub
}

func (devFieldDef *DevFieldDefinition) validate() error {
	if devFieldDef.DevId.DevDataIndex != devFieldDef.Field.DevDataIndex {
		return errors.New("developer field definition developer data index mismatch")
	}

	return (&FieldDefinition{
		DefNum:   devFieldDef.Field.DefNum,
		Size:     devFieldDef.Size,
		BaseType: devFieldDef.Field.BaseType,
	}).validate()
}

func (devFieldDef *DevFieldDefinition) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if err := devFieldDef.validate(); err != nil {
		return err
	}
	mwr := &maybe.MaybeWriter{Writer: wr}
	devFieldDef.Field.DefNum.Encode(mwr, endianness)
	devFieldDef.Size.Encode(mwr, endianness)
	devFieldDef.DevId.DevDataIndex.Encode(mwr, endianness)
	return mwr.Error()
}

type DevField struct {
	Def   *DevFieldDefinition
	Value interface{}
}

func (devField *DevField) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if err := devField.Def.validate(); err != nil {
		return err
	}
	if types.FitTypeMap[devField.Def.Field.BaseType] != reflect.TypeOf(devField.Value) {
		return ErrFieldTypeMismatch
	}

	switch devField.Value.(type) {
	case encoding.EndianEncoder:
		return devField.Value.(encoding.EndianEncoder).Encode(wr, endianness)
	case types.FitString:
		encodableStr := &types.FitEncodableString{
			FitString: devField.Value.(types.FitString),
			Length:    devField.Def.Size,
		}
		if err := encodableStr.Validate(); err != nil {
			return err
		}
		return encodableStr.Encode(wr, endianness)
	default:
		return types.ErrUnknownFitType
	}
}
