// Copyright 2021 Artem Mikheev

// Package types defines the base fit protocol types
// and provides functions for encoding most of the types.

package types

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/renbou/jogmock/strava-mock/pkg/encoding"
)

type (
	FitEnum            uint8
	FitSint8           int8
	FitUint8           uint8
	FitSint16          int16
	FitUint16          uint16
	FitSint32          int32
	FitUint32          uint32
	FitString          string
	FitEncodableString struct {
		FitString
		Length FitUint8
	}
	FitFloat32  float32
	FitFloat64  float64
	FitByte     byte
	FitSint64   int64
	FitUint64   uint64
	FitBaseType FitEnum
)

func (value FitUint8) Encode(wr io.Writer, endianness encoding.Endianness) error {
	_, err := wr.Write([]byte{uint8(value)})
	return err
}

func (value FitSint8) Encode(wr io.Writer, endianness encoding.Endianness) error {
	return FitUint8(value).Encode(wr, endianness)
}

func (value FitEnum) Encode(wr io.Writer, endianness encoding.Endianness) error {
	return FitUint8(value).Encode(wr, endianness)
}

func (value FitUint16) Encode(wr io.Writer, endianness encoding.Endianness) error {
	encoded := make([]byte, 2)
	switch endianness {
	case encoding.LittleEndian:
		binary.LittleEndian.PutUint16(encoded, uint16(value))
	case encoding.BigEndian:
		binary.BigEndian.PutUint16(encoded, uint16(value))
	default:
		panic(encoding.ErrUnknownEndianness)
	}
	_, err := wr.Write(encoded)
	return err
}

func (value FitSint16) Encode(wr io.Writer, endianness encoding.Endianness) error {
	return FitUint16(value).Encode(wr, endianness)
}

func (value FitUint32) Encode(wr io.Writer, endianness encoding.Endianness) error {
	encoded := make([]byte, 4)
	switch endianness {
	case encoding.LittleEndian:
		binary.LittleEndian.PutUint32(encoded, uint32(value))
	case encoding.BigEndian:
		binary.BigEndian.PutUint32(encoded, uint32(value))
	default:
		panic(encoding.ErrUnknownEndianness)
	}
	_, err := wr.Write(encoded)
	return err
}

func (value FitSint32) Encode(wr io.Writer, endianness encoding.Endianness) error {
	return FitUint32(value).Encode(wr, endianness)
}

func (value FitUint64) Encode(wr io.Writer, endianness encoding.Endianness) error {
	encoded := make([]byte, 8)
	switch endianness {
	case encoding.LittleEndian:
		binary.LittleEndian.PutUint64(encoded, uint64(value))
	case encoding.BigEndian:
		binary.BigEndian.PutUint64(encoded, uint64(value))
	default:
		panic(encoding.ErrUnknownEndianness)
	}
	_, err := wr.Write(encoded)
	return err
}

func (value FitSint64) Encode(wr io.Writer, endianness encoding.Endianness) error {
	return FitUint64(value).Encode(wr, endianness)
}

func (str *FitEncodableString) Validate() error {
	// a fit string must have an additional null byte
	if int(str.Length) < len(str.FitString)+1 {
		return fmt.Errorf("invalid fit string size %d (< %d expected)", str.Length, len(str.FitString)+1)
	}
	return nil
}

func (str *FitEncodableString) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if err := str.Validate(); err != nil {
		// which is why you should validate strings before calling this...
		panic(err)
	}
	data := []byte(str.FitString)
	for i := len(str.FitString); i < int(str.Length); i++ {
		data = append(data, 0)
	}
	_, err := wr.Write(data)
	return err
}

func (value FitBaseType) Encode(wr io.Writer, endianness encoding.Endianness) error {
	return FitEnum(value).Encode(wr, endianness)
}

func (baseType FitBaseType) ValidateValue(value interface{}) error {
	if FitTypeMap[baseType] != reflect.TypeOf(value) {
		return ErrFitBaseTypeMismatch
	}
	return nil
}

// IsImplemented returns true if t is an implemented fit type, false otherwise
func IsFitType(t reflect.Type) bool {
	return implementedFitTypes[t]
}
