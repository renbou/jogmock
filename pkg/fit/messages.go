// Copyright 2021 Artem Mikheev

package fit

import (
	"io"

	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit/types"
	"github.com/renbou/strava-keker/pkg/io/maybe"
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

	mwr := &maybe.MaybeWriter{Writer: wr}

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

// ConstructData constructs a data message with the given values as fields
func (localDefMsg *localDefinitionMessage) ConstructData(values ...interface{}) (*localDataMessage, error) {
	if len(values) != len(localDefMsg.def.FieldDefs)+len(localDefMsg.def.DevFieldDefs) {
		return nil, ErrFieldNumMismatch
	}
	localDataMsg := &localDataMessage{
		messageHeader: localDefMsg.messageHeader,
		def:           localDefMsg.def,
	}

	index := 0
	for _, fieldDef := range localDefMsg.def.FieldDefs {
		if err := fieldDef.BaseType.ValidateValue(values[index]); err != nil {
			return nil, err
		}
		localDataMsg.fields = append(localDataMsg.fields, &Field{
			Def:   fieldDef,
			Value: values[index],
		})
		index++
	}

	for _, devFieldDef := range localDefMsg.def.DevFieldDefs {
		if err := devFieldDef.Field.BaseType.ValidateValue(values[index]); err != nil {
			return nil, err
		}
		localDataMsg.devFields = append(localDataMsg.devFields, &DevField{
			Def:   devFieldDef,
			Value: values[index],
		})
		index++
	}
	return localDataMsg, nil
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

	mwr := &maybe.MaybeWriter{Writer: wr}

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
