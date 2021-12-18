// Copyright 2021 Artem Mikheev

package messages

import (
	"errors"
	"fmt"

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

func EncodeFitFieldDefinition(b []byte, endian types.Endianness, fieldDef *FieldDefinition) ([]byte, error) {
	if err := validateFitFieldDefinition(fieldDef); err != nil {
		return nil, err
	}
	res := types.EncodeFitUint8(b, endian, fieldDef.DefNum)
	res = append(res, types.EncodeFitUint8(b, endian, fieldDef.Size)...)
	res = append(res, types.EncodeFitBaseType(b, endian, fieldDef.BaseType)...)
	return res, nil
}

func EncodeFitField(b []byte, endian types.Endianness, fieldDef *FieldDefinition, value interface{}) ([]byte, error) {
	if err := validateFitFieldDefinition(fieldDef); err != nil {
		return nil, err
	}

	switch fieldDef.BaseType {
	case types.FIT_TYPE_ENUM:
		return types.EncodeFitEnum(b, endian, value.(types.FitEnum)), nil
	case types.FIT_TYPE_UINT8:
		return types.EncodeFitUint8(b, endian, value.(types.FitUint8)), nil
	case types.FIT_TYPE_SINT8:
		return types.EncodeFitSint8(b, endian, value.(types.FitSint8)), nil
	case types.FIT_TYPE_UINT16:
		return types.EncodeFitUint16(b, endian, value.(types.FitUint16)), nil
	case types.FIT_TYPE_SINT16:
		return types.EncodeFitSint16(b, endian, value.(types.FitSint16)), nil
	case types.FIT_TYPE_UINT32:
		return types.EncodeFitUint32(b, endian, value.(types.FitUint32)), nil
	case types.FIT_TYPE_SINT32:
		return types.EncodeFitSint32(b, endian, value.(types.FitSint32)), nil
	case types.FIT_TYPE_UINT64:
		return types.EncodeFitUint64(b, endian, value.(types.FitUint64)), nil
	case types.FIT_TYPE_SINT64:
		return types.EncodeFitSint64(b, endian, value.(types.FitSint64)), nil
	case types.FIT_TYPE_STRING:
		strValue := value.(types.FitString)
		if err := types.ValidateFitString(strValue, fieldDef.Size); err != nil {
			return nil, err
		}
		return types.EncodeFitString(b, strValue, fieldDef.Size), nil
	default:
		return nil, types.ErrUnknownFitType
	}
}
