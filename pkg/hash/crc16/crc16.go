// Copyright 2021 Artem Mikheev

// Package crc16 implements the 16-bit cyclic redundancy check, CRC-16,
// with polynomial parameters as defined by the Garmin FIT protocol.
// See https://developer.garmin.com/fit/protocol/#crc
package crc16

import (
	"encoding/binary"
	"fmt"
	"hash"
)

type crc16 struct {
	crc uint16
}

// crc_table as specified in fit protocol docs
var crcTable = [...]uint16{
	0x0000, 0xCC01, 0xD801, 0x1400, 0xF001, 0x3C00, 0x2800, 0xE401,
	0xA001, 0x6C00, 0x7800, 0xB401, 0x5000, 0x9C01, 0x8801, 0x4400,
}

// iterate does a single iteration of the crc calculation,
// also as specified in fit protocol docs
func (c *crc16) iterate(b byte) {
	var (
		tmp uint16
		crc = c.crc
	)
	tmp = crcTable[crc&0xF]
	crc = (crc >> 4) & 0x0FFF
	crc = crc ^ tmp ^ crcTable[b&0xF]

	// now compute checksum of upper four bits of byte
	tmp = crcTable[crc&0xF]
	crc = (crc >> 4) & 0x0FFF
	crc = crc ^ tmp ^ crcTable[(b>>4)&0xF]

	c.crc = crc
}

func (c *crc16) Write(bytes []byte) (n int, err error) {
	for i, b := range bytes {
		bf := c.crc
		c.iterate(b)
		if c.crc == 0 && bf != 0 {
			fmt.Println(i)
			c.crc = 0
		}
	}
	return len(bytes), nil
}

func (c *crc16) Sum(b []byte) []byte {
	crcBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcBytes, c.crc)
	return append(b, crcBytes...)
}

func (c *crc16) Size() int {
	return 2
}

func (c *crc16) Reset() {
	c.crc = 0
}

func (c *crc16) BlockSize() int {
	return 1
}

// New creates a hash.Hash computing the CRC-16/ARC hash as defined by the
// parameters poly=0x8005, init=0x0000, refIn=true, refOut=true, xorOut=0x0000.
func New() hash.Hash {
	return &crc16{crc: 0}
}

// Checksum returns the CRC-16 checksum of data
func Checksum(data []byte) uint16 {
	crc := &crc16{crc: 0}
	_, _ = crc.Write(data)
	return crc.crc
}
