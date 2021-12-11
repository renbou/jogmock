// Copyright 2021 Artem Mikheev

package types

type Enum uint8
type FitType Enum

const (
	FIT_TYPE_INVALID FitType = 0xFF
	FIT_TYPE_ENUM    FitType = 0
	FIT_TYPE_SINT8   FitType = 1
	FIT_TYPE_UINT8   FitType = 2
	FIT_TYPE_SINT16  FitType = 131
	FIT_TYPE_UINT16  FitType = 132
	FIT_TYPE_SINT32  FitType = 133
	FIT_TYPE_UINT32  FitType = 134
	FIT_TYPE_STRING  FitType = 7
	FIT_TYPE_FLOAT32 FitType = 136
	FIT_TYPE_FLOAT64 FitType = 137
	FIT_TYPE_UINT8Z  FitType = 10
	FIT_TYPE_UINT16Z FitType = 139
	FIT_TYPE_UINT32Z FitType = 140
	FIT_TYPE_BYTE    FitType = 13
	FIT_TYPE_SINT64  FitType = 142
	FIT_TYPE_UINT64  FitType = 143
	FIT_TYPE_UINT64Z FitType = 144
	FIT_TYPE_COUNT   FitType = 17
)
