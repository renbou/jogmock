// Copyright 2021 Artem Mikheev

package encoding

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	EncodableString string
	EncodableInt    int
)

func (value EncodableString) Encode(wr io.Writer, endianness Endianness) error {
	_, err := wr.Write([]byte(value))
	return err
}

func (value EncodableInt) Encode(wr io.Writer, endianness Endianness) error {
	b := make([]byte, 4)
	switch endianness {
	case LittleEndian:
		binary.LittleEndian.PutUint32(b, uint32(value))
	case BigEndian:
		binary.BigEndian.PutUint32(b, uint32(value))
	default:
		panic(ErrUnknownEndianness)
	}
	_, err := wr.Write(b)
	return err
}

func TestEncoder(t *testing.T) {
	a := assert.New(t)

	buffer := new(bytes.Buffer)
	leEncoder := NewEncoder(buffer, LittleEndian)
	if a.NoError(leEncoder.Encode(EncodableString("encodable string"))) &&
		a.NoError(leEncoder.Encode(EncodableInt(305419896))) {
		a.Equal([]byte{
			'e', 'n', 'c', 'o', 'd', 'a', 'b', 'l', 'e',
			' ', 's', 't', 'r', 'i', 'n', 'g', 0x78, 0x56, 0x34, 0x12,
		}, buffer.Bytes(), "invalid little endian encoding")
	}
	a.Equal(leEncoder.String(), "Encoder(LittleEndian)")

	buffer = new(bytes.Buffer)
	beEncoder := NewEncoder(buffer, BigEndian)
	if a.NoError(beEncoder.Encode(EncodableString("encodable string"))) &&
		a.NoError(beEncoder.Encode(EncodableInt(305419896))) {
		a.Equal([]byte{
			'e', 'n', 'c', 'o', 'd', 'a', 'b', 'l', 'e',
			' ', 's', 't', 'r', 'i', 'n', 'g', 0x12, 0x34, 0x56, 0x78,
		}, buffer.Bytes(), "invalid big endian encoding")
	}
	a.Equal(beEncoder.String(), "Encoder(BigEndian)")
}

func TestEncoderUnknownEndianness(t *testing.T) {
	a := assert.New(t)

	a.PanicsWithError(ErrUnknownEndianness.Error(), func() {
		NewEncoder(nil, Endianness(123))
	})
}
