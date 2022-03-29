// Copyright 2021 Artem Mikheev

package fit

import (
	"errors"
	"fmt"
	"io"

	"github.com/renbou/jogmock/strava-mock/pkg/encoding"
	"github.com/renbou/jogmock/strava-mock/pkg/fit/types"
	"github.com/renbou/jogmock/strava-mock/pkg/io/maybe"
)

type FieldDefinition struct {
	DefNum   types.FitUint8
	Size     types.FitUint8
	BaseType types.FitBaseType
}

// validate validates the field definition, specifically the
// specified size. Mostly for catching accidental bugs, since the field defs aren't
// all hardcoded, etc
func (fieldDef *FieldDefinition) validate() error {
	if _, ok := types.FitTypeMap[fieldDef.BaseType]; !ok {
		return types.ErrUnknownFitType
	}

	if fieldDef.BaseType == types.FIT_TYPE_STRING {
		if fieldDef.Size < 1 {
			return errors.New("field definition size with base string type must be at least 1")
		}
	} else {
		if properSize, ok := types.FitTypeSize[fieldDef.BaseType]; !ok {
			return fmt.Errorf("field definition specifies base type %v for which valid size is unknown",
				fieldDef.BaseType)
		} else if properSize != fieldDef.Size {
			return fmt.Errorf("field definition with base type %v has unexpected size %v (instead of %v)",
				fieldDef.BaseType, fieldDef.Size, properSize)
		}
	}
	return nil
}

func (fieldDef *FieldDefinition) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if err := fieldDef.validate(); err != nil {
		return err
	}
	mwr := &maybe.Writer{Writer: wr}
	fieldDef.DefNum.Encode(mwr, endianness)
	fieldDef.Size.Encode(mwr, endianness)
	fieldDef.BaseType.Encode(mwr, endianness)
	return mwr.Error()
}

type Field struct {
	Def   *FieldDefinition
	Value interface{}
}

func (field *Field) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if err := field.Def.validate(); err != nil {
		return err
	}
	if err := field.Def.BaseType.ValidateValue(field.Value); err != nil {
		return err
	}

	switch field.Value.(type) {
	case encoding.EndianEncoder:
		return field.Value.(encoding.EndianEncoder).Encode(wr, endianness)
	case types.FitString:
		encodableStr := &types.FitEncodableString{
			FitString: field.Value.(types.FitString),
			Length:    field.Def.Size,
		}
		if err := encodableStr.Validate(); err != nil {
			return err
		}
		return encodableStr.Encode(wr, endianness)
	default:
		return types.ErrUnknownFitType
	}
}
