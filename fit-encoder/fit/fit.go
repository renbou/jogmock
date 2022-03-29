// Copyright 2021 Artem Mikheev

package fit

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/renbou/jogmock/fit-encoder/encoding"
	"github.com/renbou/jogmock/fit-encoder/hash/crc16"
)

// FitFile stores sequential definition and data messages
// and can be used to later encode all of them with a valid
// fit header and footer with crc
type FitFile struct {
	messages []encoding.EndianEncoder
}

func (f *FitFile) AddMessage(message encoding.EndianEncoder) error {
	if _, ok := message.(*localDefinitionMessage); ok {
		f.messages = append(f.messages, message)
	} else if _, ok := message.(*localDataMessage); ok {
		f.messages = append(f.messages, message)
	} else {
		return ErrInvalidMessage
	}
	return nil
}

func (f *FitFile) Encode(wr io.Writer, endianness encoding.Endianness) error {
	// first encode all of the messages
	buffer := new(bytes.Buffer)
	bufferEncoder := encoding.NewEncoder(buffer, endianness)
	for _, message := range f.messages {
		bufferEncoder.Encode(message)
	}

	dataSize := buffer.Len()
	dataSizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataSizeBytes, uint32(dataSize))
	dataCrc := crc16.New()
	dataCrc.Write(buffer.Bytes())

	// write header
	header := []byte{0x0E, 0x20, 0x54, 0x08}
	header = append(header, dataSizeBytes...)
	header = append(header, []byte(".FIT")...)
	headerCrc := crc16.New()
	headerCrc.Write(header)
	header = append(header, headerCrc.Sum(nil)...)
	_, err := wr.Write(header)
	if err != nil {
		return err
	}

	// write actual messages
	_, err = wr.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	// write data crc
	_, err = wr.Write(dataCrc.Sum(nil))
	if err != nil {
		return err
	}

	return nil
}
