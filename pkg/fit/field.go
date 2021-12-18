// Copyright 2021 Artem Mikheev

package fit

import (
	"errors"
	"fmt"
	"io"

	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit/types"
)

type FieldDefinition struct {
	DefNum   types.FitUint8
	Size     types.FitUint8
	BaseType types.FitBaseType
}

// validateFitFieldDefinition validates the field definition, specifically the
// specified size. Mostly for catching accidental bugs, since the field defs aren't
// all hardcoded, etc
func validateFitFieldDefinition(fieldDef *FieldDefinition) error {
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
	if err := validateFitFieldDefinition(fieldDef); err != nil {
		return err
	}
	if err := fieldDef.DefNum.Encode(wr, endianness); err != nil {
		return err
	}
	if err := fieldDef.Size.Encode(wr, endianness); err != nil {
		return err
	}
	if err := fieldDef.BaseType.Encode(wr, endianness); err != nil {
		return err
	}
	return nil
}

type Field struct {
	Def   *FieldDefinition
	Value interface{}
}

func (field *Field) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if err := validateFitFieldDefinition(field.Def); err != nil {
		return err
	}

	switch field.Value.(type) {
	case encoding.EndianEncoder:
		return field.Value.(encoding.EndianEncoder).Encode(wr, endianness)
	case types.FitString:
		encodableStr := &types.FitEncodableString{
			field.Value.(types.FitString),
			field.Def.Size,
		}
		if err := encodableStr.Validate(); err != nil {
			return err
		}
		return encodableStr.Encode(wr, endianness)
	default:
		return types.ErrUnknownFitType
	}
}
