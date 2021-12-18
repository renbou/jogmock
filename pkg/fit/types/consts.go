// Copyright 2021 Artem Mikheev

package types

type FitEnum uint8
type FitSint8 int8
type FitUint8 uint8
type FitSint16 int16
type FitUint16 uint16
type FitSint32 int32
type FitUint32 uint32
type FitString string
type FitFloat32 float32
type FitFloat64 float64
type FitByte byte
type FitSint64 int64
type FitUint64 uint64
type FitBaseType FitEnum

const (
	FIT_TYPE_BYTE_SIZE    FitUint8 = 1
	FIT_TYPE_ENUM_SIZE    FitUint8 = FIT_TYPE_BYTE_SIZE
	FIT_TYPE_UINT8_SIZE   FitUint8 = 1
	FIT_TYPE_SINT8_SIZE   FitUint8 = FIT_TYPE_UINT8_SIZE
	FIT_TYPE_UINT16_SIZE  FitUint8 = 2
	FIT_TYPE_SINT16_SIZE  FitUint8 = FIT_TYPE_UINT16_SIZE
	FIT_TYPE_UINT32_SIZE  FitUint8 = 4
	FIT_TYPE_SINT32_SIZE  FitUint8 = FIT_TYPE_UINT32_SIZE
	FIT_TYPE_STRING_SIZE  FitUint8 = 1
	FIT_TYPE_FLOAT32_SIZE FitUint8 = 4
	FIT_TYPE_FLOAT64_SIZE FitUint8 = 8
	FIT_TYPE_UINT64_SIZE  FitUint8 = 8
	FIT_TYPE_SINT64_SIZE  FitUint8 = FIT_TYPE_UINT64_SIZE
)

const (
	FIT_TYPE_INVALID FitBaseType = 0xFF
	FIT_TYPE_ENUM    FitBaseType = 0
	FIT_TYPE_SINT8   FitBaseType = 1
	FIT_TYPE_UINT8   FitBaseType = 2
	FIT_TYPE_SINT16  FitBaseType = 131
	FIT_TYPE_UINT16  FitBaseType = 132
	FIT_TYPE_SINT32  FitBaseType = 133
	FIT_TYPE_UINT32  FitBaseType = 134
	FIT_TYPE_STRING  FitBaseType = 7
	FIT_TYPE_FLOAT32 FitBaseType = 136
	FIT_TYPE_FLOAT64 FitBaseType = 137
	FIT_TYPE_UINT8Z  FitBaseType = 10
	FIT_TYPE_UINT16Z FitBaseType = 139
	FIT_TYPE_UINT32Z FitBaseType = 140
	FIT_TYPE_BYTE    FitBaseType = 13
	FIT_TYPE_SINT64  FitBaseType = 142
	FIT_TYPE_UINT64  FitBaseType = 143
	FIT_TYPE_UINT64Z FitBaseType = 144
	FIT_TYPE_COUNT   FitBaseType = 17
)
