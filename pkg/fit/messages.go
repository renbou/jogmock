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

type MessageHeader struct {
	LocalMsgType uint8
}

type DefinitionMessage struct {
	MessageHeader
	GlobalMsgNum types.FitUint16
	FieldDefs    []*FieldDefinition
	DevFieldDefs []*DevFieldDefinition
}

func (defMsg *DefinitionMessage) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if !endianness.IsKnown() {
		return encoding.ErrUnknownEndianness
	}

	mwr := &maybe.MaybeWriter{Writer: wr}

	// Encode the single-byte message header
	encodeMessageHeader(mwr, defMsgType, len(defMsg.DevFieldDefs) > 0, defMsg.LocalMsgType)

	// Write reserved byte and endianness (architecture) byte
	mwr.Write([]byte{0, endiannessByteMap[endianness]})

	// Encode global message num
	defMsg.GlobalMsgNum.Encode(mwr, endianness)

	// Encode number of field defs as a single byte
	types.FitUint8(len(defMsg.FieldDefs)).Encode(mwr, endianness)

	// Encode all field defs
	for _, fieldDef := range defMsg.FieldDefs {
		fieldDef.Encode(mwr, endianness)
	}

	// Encode dev fields if we need to
	if len(defMsg.DevFieldDefs) > 0 {
		types.FitUint8(len(defMsg.DevFieldDefs)).Encode(mwr, endianness)
		for _, devFieldDef := range defMsg.DevFieldDefs {
			devFieldDef.Encode(mwr, endianness)
		}
	}
	return mwr.Error()
}

// ConstructData constructs a data message with the given values as fields
func (defMsg *DefinitionMessage) ConstructData(values ...interface{}) (*DataMessage, error) {
	if len(values) != len(defMsg.FieldDefs)+len(defMsg.DevFieldDefs) {
		return nil, ErrFieldNumMismatch
	}
	dataMsg := &DataMessage{
		Def: defMsg,
	}

	index := 0
	for _, fieldDef := range defMsg.FieldDefs {
		dataMsg.Fields = append(dataMsg.Fields, &Field{
			Def:   fieldDef,
			Value: values[index],
		})
		index++
	}

	for _, devFieldDef := range defMsg.DevFieldDefs {
		dataMsg.DevFields = append(dataMsg.DevFields, &DevField{
			Def:   devFieldDef,
			Value: values[index],
		})
		index++
	}
	return dataMsg, nil
}

type DataMessage struct {
	Def       *DefinitionMessage
	Fields    []*Field
	DevFields []*DevField
}

func (dataMsg *DataMessage) Encode(wr io.Writer, endianness encoding.Endianness) error {
	return nil
}

// Global messages
var DeveloperDataId = &DefinitionMessage{}
