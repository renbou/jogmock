// Copyright 2021 Artem Mikheev

package encoding

import "io"

// EndianEncoder interface is a generic encodable value
type EndianEncoder interface {
	Encode(io.Writer, Endianness) error
}

// Encoder is a type which incapsulates encoding of multiple
// EndianEncoder objects into a single writer with configured endianness
type Encoder struct {
	endianness Endianness
	wr         io.Writer
}

func NewEncoder(wr io.Writer, endianness Endianness) *Encoder {
	if !endianness.IsKnown() {
		panic(ErrUnknownEndianness)
	}
	return &Encoder{
		endianness: endianness,
		wr:         wr,
	}
}

func (enc *Encoder) String() string {
	return "Encoder(" + enc.endianness.String() + ")"
}

func (enc *Encoder) Encode(value EndianEncoder) error {
	return value.Encode(enc.wr, enc.endianness)
}
