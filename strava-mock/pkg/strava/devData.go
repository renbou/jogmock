// Copyright 2021 Artem Mikheev

package strava

import (
	"errors"

	"github.com/renbou/jogmock/strava-mock/pkg/encoding"
	"github.com/renbou/jogmock/strava-mock/pkg/fit"
	"github.com/renbou/jogmock/strava-mock/pkg/fit/types"
)

// devData is a helper for working with developer data
// messages and fields, helps in simplifying the dev data logic
type devData struct {
	index          uint8
	appVersion     uint32
	fields         map[string]*devField
	alreadyEncoded bool
}

func (dev *devData) addField(name string, fitType types.FitBaseType) error {
	if dev.alreadyEncoded {
		return errors.New("developer data has already been encoded")
	}
	if _, ok := dev.fields[name]; ok {
		return errors.New("field already exists")
	}
	dev.fields[name] = &devField{
		name:    name,
		defNum:  uint8(len(dev.fields)),
		fitType: fitType,
	}
	return nil
}

// Construct developer data id message and all of the assigned developer
// field definition messages using availableLocalMsgNum.
// After this you can reuse availableLocalMsgNum
func (dev *devData) constructAllMessages(availableLocalMsgNum uint8) ([]encoding.EndianEncoder, error) {
	// even if we fail we encoded some data soo...
	defer func() {
		dev.alreadyEncoded = true
	}()

	messages := make([]encoding.EndianEncoder, 0)

	developerDataIdMessage := &fit.DefinitionMessage{
		GlobalMsgNum: 207,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   3,
				Size:     types.FIT_TYPE_UINT8_SIZE,
				BaseType: types.FIT_TYPE_UINT8,
			},
			{
				DefNum:   4,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
		},
		DevFieldDefs: nil,
	}
	localDeveloperDataIdMsgDef := developerDataIdMessage.ConstructLocalMessage(availableLocalMsgNum)
	messages = append(messages, localDeveloperDataIdMsgDef)

	localDeveloperDataIdMsg, err := localDeveloperDataIdMsgDef.ConstructData(dev.index, dev.appVersion)
	if err != nil {
		return nil, err
	}
	messages = append(messages, localDeveloperDataIdMsg)

	if len(dev.fields) == 0 {
		return nil, errors.New("no developer fields to encode")
	}

	maxNameLen := 0
	for _, field := range dev.fields {
		if curLen := len(field.name); curLen > maxNameLen {
			maxNameLen = curLen
		}
	}

	if maxNameLen > 250 {
		return nil, errors.New(
			"unable to encode fields with name length > 250",
		)
	}

	fieldDescriptionMessage := &fit.DefinitionMessage{
		GlobalMsgNum: 206,
		FieldDefs: []*fit.FieldDefinition{
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_UINT8_SIZE,
				BaseType: types.FIT_TYPE_UINT8,
			},
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_UINT8_SIZE,
				BaseType: types.FIT_TYPE_UINT8,
			},
			{
				DefNum:   2,
				Size:     types.FIT_TYPE_UINT8_SIZE,
				BaseType: types.FIT_TYPE_UINT8,
			},
			{
				DefNum:   3,
				Size:     types.FitUint8(maxNameLen) + 1,
				BaseType: types.FIT_TYPE_STRING,
			},
		},
		DevFieldDefs: nil,
	}
	localFieldDescriptionMsgDef := fieldDescriptionMessage.ConstructLocalMessage(availableLocalMsgNum)
	messages = append(messages, localFieldDescriptionMsgDef)

	for _, field := range dev.fields {
		localFieldDescriptionMsg, err := localFieldDescriptionMsgDef.ConstructData(dev.index, field.defNum, field.fitType, field.name)
		if err != nil {
			return nil, err
		}
		messages = append(messages, localFieldDescriptionMsg)
	}

	return messages, nil
}

// getFieldDefinition returns a fit.DevFieldDefinition for
// the defined developer field with given name
// size parameter is only needed for string fields
func (dev *devData) getFieldDefinition(name string, size ...int) (*fit.DevFieldDefinition, error) {
	if !dev.alreadyEncoded {
		return nil, errors.New(
			"developer data id and field definition messages not yet encoded",
		)
	}

	field, ok := dev.fields[name]
	if !ok {
		return nil, errors.New("no such developer field defined")
	}

	var fieldSize uint
	if field.fitType == types.FIT_TYPE_STRING {
		if len(size) != 1 {
			return nil, errors.New(
				"must specify single size for string-type field",
			)
		}
		if size[0] < 1 {
			return nil, errors.New(
				"string-typed field must have size >= 1",
			)
		}
		fieldSize = uint(size[0] + 1)
	} else {
		fitTypeSize, ok := types.FitTypeSize[field.fitType]
		if !ok {
			return nil, errors.New(
				"encountered fit field with unknown size",
			)
		}
		fieldSize = uint(fitTypeSize)
	}

	return &fit.DevFieldDefinition{
		Field: &fit.FieldDescriptionStub{
			DevDataIndex: types.FitUint8(dev.index),
			DefNum:       types.FitUint8(field.defNum),
			BaseType:     field.fitType,
		},
		Size: types.FitUint8(fieldSize),
		DevId: &fit.DeveloperDataIdStub{
			DevDataIndex: types.FitUint8(dev.index),
		},
	}, nil
}

// devField represents a single developer field
// owned by some developer index
type devField struct {
	name    string
	defNum  uint8
	fitType types.FitBaseType
}
