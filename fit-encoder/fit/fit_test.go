// Copyright 2021 Artem Mikheev

package fit

import (
	"encoding/hex"
	"testing"

	"github.com/renbou/jogmock/strava-mock/pkg/encoding"
	"github.com/renbou/jogmock/strava-mock/pkg/fit/types"
	"github.com/stretchr/testify/assert"
)

func TestFitFileEncoding(t *testing.T) {
	a := assert.New(t)

	fitFile := new(FitFile)

	fileIdMsgDef := &DefinitionMessage{
		GlobalMsgNum: 0,
		FieldDefs: []*FieldDefinition{
			{
				DefNum:   1,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   2,
				Size:     types.FIT_TYPE_UINT16_SIZE,
				BaseType: types.FIT_TYPE_UINT16,
			},
			{
				DefNum:   0,
				Size:     types.FIT_TYPE_ENUM_SIZE,
				BaseType: types.FIT_TYPE_ENUM,
			},
			{
				DefNum:   4,
				Size:     types.FIT_TYPE_UINT32_SIZE,
				BaseType: types.FIT_TYPE_UINT32,
			},
		},
		DevFieldDefs: nil,
	}

	localFileIdMsgDef := fileIdMsgDef.ConstructLocalMessage(0)
	fitFile.AddMessage(localFileIdMsgDef)

	localFileIdMsg, err := localFileIdMsgDef.ConstructData(265, 102, 4, 1007562558)
	if err != nil {
		t.Fatalf("Unexpected error during data construction for file id: %v", err)
	}
	fitFile.AddMessage(localFileIdMsg)

	fitFileBytes, _ := hex.DecodeString(
		"0e2054081c0000002e464954b8884000010000040102840202840001000404860001090066043c0e2f3ede5f",
	)
	encodeAndValidate(a, fitFile, encoding.BigEndian, true, fitFileBytes)
}
