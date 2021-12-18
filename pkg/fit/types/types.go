// Copyright 2021 Artem Mikheev

// Package types defines the base fit protocol types
// and provides functions for encoding most of the types.

package types

import (
	"encoding/binary"
	"fmt"
)

type Endianness int

//go:generate stringer -type Endianness
const (
	LittleEndian Endianness = iota
	BigEndian
	NonEndian
)

func EncodeFitUint8(b []byte, endian Endianness, value FitUint8) []byte {
	encoded := []byte{uint8(value)}
	return append(b, encoded...)
}

func EncodeFitSint8(b []byte, endian Endianness, value FitSint8) []byte {
	return EncodeFitUint8(b, endian, FitUint8(value))
}

func EncodeFitEnum(b []byte, endian Endianness, value FitEnum) []byte {
	return EncodeFitUint8(b, endian, FitUint8(value))
}

func EncodeFitUint16(b []byte, endian Endianness, value FitUint16) []byte {
	encoded := make([]byte, 2)
	switch endian {
	case LittleEndian:
		binary.LittleEndian.PutUint16(encoded, uint16(value))
	case BigEndian:
		binary.BigEndian.PutUint16(encoded, uint16(value))
	default:
		panic(fmt.Sprintf("unknown endianness %d", endian))
	}
	return append(b, encoded...)
}

func EncodeFitSint16(b []byte, endian Endianness, value FitSint16) []byte {
	return EncodeFitUint16(b, endian, FitUint16(value))
}

func EncodeFitUint32(b []byte, endian Endianness, value FitUint32) []byte {
	encoded := make([]byte, 4)
	switch endian {
	case LittleEndian:
		binary.LittleEndian.PutUint32(encoded, uint32(value))
	case BigEndian:
		binary.BigEndian.PutUint32(encoded, uint32(value))
	default:
		panic(fmt.Sprintf("unknown endianness %d", endian))
	}
	return append(b, encoded...)
}

func EncodeFitSint32(b []byte, endian Endianness, value FitSint32) []byte {
	return EncodeFitUint32(b, endian, FitUint32(value))
}

func EncodeFitUint64(b []byte, endian Endianness, value FitUint64) []byte {
	encoded := make([]byte, 8)
	switch endian {
	case LittleEndian:
		binary.LittleEndian.PutUint64(encoded, uint64(value))
	case BigEndian:
		binary.BigEndian.PutUint64(encoded, uint64(value))
	default:
		panic(fmt.Sprintf("unknown endianness %d", endian))
	}
	return append(b, encoded...)
}

func EncodeFitSint64(b []byte, endian Endianness, value FitSint64) []byte {
	return EncodeFitUint64(b, endian, FitUint64(value))
}

func ValidateFitString(value FitString, length FitUint8) error {
	// a fit string must have an additional null byte
	if int(length) < len(value)+1 {
		return fmt.Errorf("invalid fit string size %d (< %d expected)", length, len(value)+1)
	}
	return nil
}

func EncodeFitString(b []byte, value FitString, length FitUint8) []byte {
	if err := ValidateFitString(value, length); err != nil {
		// which is why you should validate strings before calling this...
		panic(err)
	}
	res := append(b, []byte(value)...)
	for i := len(value); i < int(length); i++ {
		res = append(res, 0)
	}
	return res
}

func EncodeFitBaseType(b []byte, endian Endianness, value FitBaseType) []byte {
	return EncodeFitEnum(b, endian, FitEnum(value))
}
