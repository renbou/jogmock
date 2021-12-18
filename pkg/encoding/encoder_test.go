// Copyright 2021 Artem Mikheev

package encoding

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	fitTesting "github.com/renbou/strava-keker/internal/testing"
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
	buffer := new(bytes.Buffer)
	leEncoder := NewEncoder(buffer, LittleEndian)
	if err := leEncoder.Encode(EncodableString("encodable string")); err != nil {
		t.Fatal(err)
	}
	if err := leEncoder.Encode(EncodableInt(305419896)); err != nil {
		t.Fatal(err)
	}
	if err := fitTesting.AssertEqual([]byte{
		'e', 'n', 'c', 'o', 'd', 'a', 'b', 'l', 'e',
		' ', 's', 't', 'r', 'i', 'n', 'g', 0x78, 0x56, 0x34, 0x12,
	}, buffer.Bytes()); err != nil {
		t.Fatalf("Invalid encoding with %s: %v", leEncoder, err)
	}

	buffer = new(bytes.Buffer)
	beEncoder := NewEncoder(buffer, BigEndian)
	if err := beEncoder.Encode(EncodableString("encodable string")); err != nil {
		t.Fatal(err)
	}
	if err := beEncoder.Encode(EncodableInt(305419896)); err != nil {
		t.Fatal(err)
	}
	if err := fitTesting.AssertEqual([]byte{
		'e', 'n', 'c', 'o', 'd', 'a', 'b', 'l', 'e',
		' ', 's', 't', 'r', 'i', 'n', 'g', 0x12, 0x34, 0x56, 0x78,
	}, buffer.Bytes()); err != nil {
		t.Fatalf("Invalid encoding with %s: %v", beEncoder, err)
	}
}
