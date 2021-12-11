// Copyright 2021 Artem Mikheev

// Package types defines the base fit protocol types
// and provides functions for encoding most of the types.

package types

import (
	"encoding/binary"
	"fmt"
)

type Endianness int

const (
	LittleEndian Endianness = iota
	BigEndian
)

func EncodeUint8(b []byte, endian Endianness, value uint8) []byte {
	encoded := []byte{value}
	return append(b, encoded...)
}

func EncodeSint8(b []byte, endian Endianness, value int8) []byte {
	return EncodeUint8(b, endian, uint8(value))
}

func EncodeEnum(b []byte, endian Endianness, value Enum) []byte {
	return EncodeUint8(b, endian, uint8(value))
}

func EncodeUint16(b []byte, endian Endianness, value uint16) []byte {
	encoded := make([]byte, 2)
	switch endian {
	case LittleEndian:
		binary.LittleEndian.PutUint16(encoded, value)
	case BigEndian:
		binary.BigEndian.PutUint16(encoded, value)
	default:
		panic(fmt.Sprintf("unknown endianness %d", endian))
	}
	return append(b, encoded...)
}

func EncodeSint16(b []byte, endian Endianness, value int16) []byte {
	return EncodeUint16(b, endian, uint16(value))
}

func EncodeUint32(b []byte, endian Endianness, value uint32) []byte {
	encoded := make([]byte, 4)
	switch endian {
	case LittleEndian:
		binary.LittleEndian.PutUint32(encoded, value)
	case BigEndian:
		binary.BigEndian.PutUint32(encoded, value)
	default:
		panic(fmt.Sprintf("unknown endianness %d", endian))
	}
	return append(b, encoded...)
}

func EncodeSint32(b []byte, endian Endianness, value int32) []byte {
	return EncodeUint32(b, endian, uint32(value))
}

func EncodeUint64(b []byte, endian Endianness, value uint64) []byte {
	encoded := make([]byte, 8)
	switch endian {
	case LittleEndian:
		binary.LittleEndian.PutUint64(encoded, value)
	case BigEndian:
		binary.BigEndian.PutUint64(encoded, value)
	default:
		panic(fmt.Sprintf("unknown endianness %d", endian))
	}
	return append(b, encoded...)
}

func EncodeSint64(b []byte, endian Endianness, value int64) []byte {
	return EncodeUint64(b, endian, uint64(value))
}

func EncodeString(b []byte, value string) []byte {
	return append(append(b, []byte(value)...), 0)
}

func EncodeFitType(b []byte, endian Endianness, value FitType) []byte {
	return EncodeEnum(b, endian, Enum(value))
}
