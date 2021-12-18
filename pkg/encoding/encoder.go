// Copyright 2021 Artem Mikheev

package encoding

import "io"

type EndianEncoder interface {
	Encode(io.Writer, Endianness) error
}

type Encoder struct {
	endianness Endianness
	wr         io.Writer
}

func NewEncoder(wr io.Writer, endianness Endianness) *Encoder {
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
