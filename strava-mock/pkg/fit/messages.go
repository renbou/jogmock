// Copyright 2021 Artem Mikheev

package fit

import (
	"fmt"
	"io"
	"reflect"

	"github.com/renbou/jogmock/strava-mock/pkg/encoding"
	"github.com/renbou/jogmock/strava-mock/pkg/fit/types"
	"github.com/renbou/jogmock/strava-mock/pkg/io/maybe"
)

const (
	dataMsgType uint8 = iota
	defMsgType
)

var endiannessByteMap = map[encoding.Endianness]byte{
	encoding.LittleEndian: 0,
	encoding.BigEndian:    1,
}

func encodeMessageHeader(wr io.Writer, msgType uint8, msgSpecific bool, localMsgType uint8) error {
	if localMsgType > 15 {
		return ErrInvalidLocalMsgType
	}
	messageHeader := localMsgType

	if msgType > 1 {
		return ErrInvalidMsgType
	}
	messageHeader |= uint8(msgType) << 6

	if msgSpecific {
		if msgType == dataMsgType {
			return ErrInvalidMsgSpecific
		}
		messageHeader |= 1 << 5
	}

	_, err := wr.Write([]byte{messageHeader})
	return err
}

type messageHeader struct {
	localMsgType uint8
}

type DefinitionMessage struct {
	GlobalMsgNum types.FitUint16
	FieldDefs    []*FieldDefinition
	DevFieldDefs []*DevFieldDefinition
}

func (defMsg *DefinitionMessage) ConstructLocalMessage(localMsgType uint8) *localDefinitionMessage {
	return &localDefinitionMessage{
		messageHeader: messageHeader{
			localMsgType,
		},
		def: defMsg,
	}
}

// localDefinitionMessage is a DefinitionMessage tied to a
// specific LocalMsgType, and is thus encodable
type localDefinitionMessage struct {
	messageHeader
	def *DefinitionMessage
}

func (localDefMsg *localDefinitionMessage) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if !endianness.IsKnown() {
		return encoding.ErrUnknownEndianness
	}

	mwr := &maybe.Writer{Writer: wr}

	// Encode the single-byte message header
	encodeMessageHeader(mwr, defMsgType, len(localDefMsg.def.DevFieldDefs) > 0, localDefMsg.localMsgType)

	// Write reserved byte and endianness (architecture) byte
	mwr.Write([]byte{0, endiannessByteMap[endianness]})

	// Encode global message num
	localDefMsg.def.GlobalMsgNum.Encode(mwr, endianness)

	// Encode number of field defs as a single byte
	types.FitUint8(len(localDefMsg.def.FieldDefs)).Encode(mwr, endianness)

	// Encode all field defs
	for _, fieldDef := range localDefMsg.def.FieldDefs {
		fieldDef.Encode(mwr, endianness)
	}

	// Encode dev fields if we need to
	if len(localDefMsg.def.DevFieldDefs) > 0 {
		types.FitUint8(len(localDefMsg.def.DevFieldDefs)).Encode(mwr, endianness)
		for _, devFieldDef := range localDefMsg.def.DevFieldDefs {
			devFieldDef.Encode(mwr, endianness)
		}
	}
	return mwr.Error()
}

// constructDataImpl is the actual implementation of ConstructData
// which takes only defined fit base types (in fit.types)
func (localDefMsg *localDefinitionMessage) constructDataImpl(values []interface{}) (*localDataMessage, error) {
	localDataMsg := &localDataMessage{
		messageHeader: localDefMsg.messageHeader,
		def:           localDefMsg.def,
	}

	index := 0
	for _, fieldDef := range localDefMsg.def.FieldDefs {
		localDataMsg.fields = append(localDataMsg.fields, &Field{
			Def:   fieldDef,
			Value: values[index],
		})
		index++
	}

	for _, devFieldDef := range localDefMsg.def.DevFieldDefs {
		localDataMsg.devFields = append(localDataMsg.devFields, &DevField{
			Def:   devFieldDef,
			Value: values[index],
		})
		index++
	}
	return localDataMsg, nil
}

// convertValue converts an arbitrary value to the respective
// fit base type value
func convertValue(value interface{}, fitBaseType types.FitBaseType) (interface{}, error) {
	currentValue := reflect.ValueOf(value)
	currentType := currentValue.Type()
	if types.IsFitType(currentType) {
		// if current value is already a fit type, validate it
		if err := fitBaseType.ValidateValue(value); err != nil {
			return nil, err
		}
		return value, nil
	} else {
		if currentFitType, ok := types.FitTypeMap[fitBaseType]; !ok {
			return nil, fmt.Errorf("unable to convert value to field with unknown base type %v", fitBaseType)
		} else if !currentValue.CanConvert(currentFitType) {
			return nil, fmt.Errorf("unable to convert value %s of type %s to field with base type %s",
				currentValue, currentType, currentFitType,
			)
		} else {
			return currentValue.Convert(currentFitType).Interface(), nil
		}
	}
}

func (localDefMsg *localDefinitionMessage) ConstructData(values ...interface{}) (*localDataMessage, error) {
	if len(values) != len(localDefMsg.def.FieldDefs)+len(localDefMsg.def.DevFieldDefs) {
		return nil, ErrFieldNumMismatch
	}

	index := 0
	convertedValues := make([]interface{}, len(values))
	for _, fieldDef := range localDefMsg.def.FieldDefs {
		converted, err := convertValue(values[index], fieldDef.BaseType)
		if err != nil {
			return nil, err
		}
		convertedValues[index] = converted
		index++
	}
	for _, devFieldDef := range localDefMsg.def.DevFieldDefs {
		converted, err := convertValue(values[index], devFieldDef.Field.BaseType)
		if err != nil {
			return nil, err
		}
		convertedValues[index] = converted
		index++
	}
	return localDefMsg.constructDataImpl(convertedValues)
}

type localDataMessage struct {
	messageHeader
	def       *DefinitionMessage
	fields    []*Field
	devFields []*DevField
}

func (localDataMsg *localDataMessage) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if !endianness.IsKnown() {
		return encoding.ErrUnknownEndianness
	}

	mwr := &maybe.Writer{Writer: wr}

	// Encode the single-byte message header
	encodeMessageHeader(mwr, dataMsgType, false, localDataMsg.localMsgType)

	// Encode fields
	for _, field := range localDataMsg.fields {
		field.Encode(mwr, endianness)
	}
	for _, devField := range localDataMsg.devFields {
		devField.Encode(mwr, endianness)
	}

	return mwr.Error()
}
