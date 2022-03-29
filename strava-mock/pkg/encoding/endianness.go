// Copyright 2021 Artem Mikheev

package encoding

import "errors"

type Endianness int

//go:generate stringer -type Endianness
const (
	LittleEndian Endianness = iota
	BigEndian
)

var ErrUnknownEndianness error = errors.New("unknown endianness")

func (endianness Endianness) IsKnown() bool {
	if endianness != LittleEndian && endianness != BigEndian {
		return false
	}
	return true
}
