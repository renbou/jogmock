// Copyright 2021 Artem Mikheev

package fit

import (
	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit/types"
)

type MessageType uint8

const (
	DataMessageType       MessageType = 0
	DefinitionMessageType MessageType = 1
)

type MessageHeader struct {
	msgType MessageType
}

type DefinitionMessage struct {
	endianness   encoding.Endianness
	msgNum       types.FitUint16
	fieldDefs    []FieldDefinition
	devFieldDefs []DevFieldDefinition
}
