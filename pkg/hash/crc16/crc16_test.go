// Copyright 2021 Artem Mikheev
package crc16

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash"
	"testing"
)

func assertCrc16(b []byte, expectedSum uint16) error {
	var (
		crc       hash.Hash
		actualSum uint16
		err       error
	)

	crc = New()
	n, err := crc.Write(b)

	if err != nil {
		return err
	} else if n != len(b) {
		return fmt.Errorf(
			"crc calculation failed, written %v out of %v bytes", n, len(b))
	}

	buffer := crc.Sum(nil)
	err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &actualSum)
	if err != nil {
		return err
	}
	if actualSum != expectedSum {
		return fmt.Errorf(
			"crc actual sum (0x%x) not equal to expected sum (0x%x)",
			actualSum, expectedSum)
	}

	if checksum := Checksum(b); checksum != actualSum {
		return fmt.Errorf("Checksum (0x%x) not equal to calculated actual sum (0x%x)",
			checksum, actualSum,
		)
	}

	return nil
}

func TestValidateCrc16Constants(t *testing.T) {
	crc := New()

	if blockSize := crc.BlockSize(); blockSize != 1 {
		t.Fatalf("crc.BlockSize() (= %v) != %v", blockSize, 1)
	}

	if sumSize := crc.Size(); sumSize != 2 {
		t.Fatalf("crc.Size() (= %v) != %v", sumSize, 2)
	}
}

func TestCrc16(t *testing.T) {
	// simple tests
	if err := assertCrc16([]byte("aboba"), 0xE57A); err != nil {
		t.Error(err)
	}

	if err := assertCrc16([]byte("123456789"), 0xBB3D); err != nil {
		t.Error(err)
	}

	// no bytes test
	if err := assertCrc16([]byte{}, 0x0000); err != nil {
		t.Error(err)
	}

	// long test
	lotsOfAAA := make([]byte, 9000)
	for i := range lotsOfAAA {
		lotsOfAAA[i] = 'A'
	}
	if err := assertCrc16(lotsOfAAA, 0x43EA); err != nil {
		t.Error(err)
	}
}
