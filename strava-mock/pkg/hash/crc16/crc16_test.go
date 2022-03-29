// Copyright 2021 Artem Mikheev
package crc16

import (
	"bytes"
	"encoding/binary"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertCrc16(a *assert.Assertions, crc hash.Hash, b []byte, expectedSum uint16) {
	n, err := crc.Write(b)
	if !a.NoError(err) || !a.Len(b, n, "crc written less bytes than slice length") {
		return
	}

	var actualSum uint16
	err = binary.Read(bytes.NewBuffer(crc.Sum(nil)), binary.LittleEndian, &actualSum)
	if !a.NoError(err) {
		return
	}

	a.Equalf(expectedSum, actualSum, "sums differ for %s", string(b))
	a.Equalf(actualSum, Checksum(b), "checksum() calculates invalid sum for %s", string(b))
}

func TestCrc16Constants(t *testing.T) {
	a := assert.New(t)
	crc := New()

	a.Equal(1, crc.BlockSize())
	a.Equal(2, crc.Size())
}

func TestCrc16Sum(t *testing.T) {
	a := assert.New(t)
	crc := New()

	// Simple tests
	tests := []struct {
		data string
		sum  uint16
	}{
		{"aboba", 0xE57A},
		{"123456789", 0xBB3D},
		{"fopsfudhaw90fh-10293", 0x1243},
		{"", 0x0000},
	}
	for _, test := range tests {
		assertCrc16(a, crc, []byte(test.data), test.sum)
		crc.Reset()
	}

	// Long test
	lotsOfAAA := make([]byte, 9000)
	for i := range lotsOfAAA {
		lotsOfAAA[i] = 'A'
	}
	assertCrc16(a, crc, lotsOfAAA, 0x43EA)
}
