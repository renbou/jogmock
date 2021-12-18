// Copyright 2021 Artem Mikheev

package encoding

import "errors"

type Endianness int

//go:generate stringer -type Endianness
const (
	LittleEndian Endianness = iota
	BigEndian
	NonEndian
)

var ErrUnknownEndianness error = errors.New("unknown endianness")
