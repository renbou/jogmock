package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// assertEncoded asserts that the value actually encoded into what we expected,
// otherwise it prints a pretty informative message. For proper debug output,
// make sure you explicitly specify value's int type (ex: int8(-4))
func assertEncoded(t *testing.T, encoderName string, endianness Endianness,
	value interface{}, actual []byte, expected []byte) {

	var endianString string
	switch endianness {
	case LittleEndian:
		endianString = "LittleEndian"
	case BigEndian:
		endianString = "BigEndian"
	default:
		endianString = "NotEndianLol"
	}

	var valueRepr string
	if strings.HasPrefix(reflect.TypeOf(value).Name(), "uint") {
		valueRepr = fmt.Sprintf("0x%x", value)
	} else {
		valueRepr = fmt.Sprintf("%v", value)
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("%s(%s, %s) failed, expected %s, but got %s",
			encoderName, endianString, valueRepr,
			hex.EncodeToString(expected), hex.EncodeToString(actual))
	}
}

func TestEncoding(t *testing.T) {
	var (
		b                 []byte
		currentEncoder    string
		currentEndianness Endianness
	)

	// helper for assertions with better debug output
	assert := func(value interface{}, expected []byte) {
		assertEncoded(t, currentEncoder, currentEndianness, value, b, expected)
	}

	// uint8
	currentEncoder = "EncodeUint8"
	currentEndianness = LittleEndian
	b = EncodeUint8(nil, LittleEndian, 0x69)
	assert(uint8(0x69), []byte{0x69})

	currentEndianness = BigEndian
	b = EncodeUint8(nil, BigEndian, 0x42)
	assert(uint8(0x42), []byte{0x42})

	// sint8
	currentEncoder = "EncodeSint8"
	currentEndianness = LittleEndian
	b = EncodeSint8(nil, LittleEndian, 100)
	assert(int8(100), []byte{0x64})

	b = EncodeSint8(nil, LittleEndian, -100)
	assert(int8(-100), []byte{0x9c})

	currentEndianness = BigEndian
	b = EncodeSint8(nil, BigEndian, 34)
	assert(int8(34), []byte{0x22})

	b = EncodeSint8(nil, BigEndian, -34)
	assert(int8(-34), []byte{0xde})

	// enum
	currentEncoder = "EncodeEnum"
	currentEndianness = LittleEndian
	b = EncodeEnum(nil, LittleEndian, 0)
	assert(Enum(0), []byte{0})

	currentEndianness = BigEndian
	b = EncodeUint8(nil, BigEndian, 19)
	assert(Enum(19), []byte{0x13})

	// uint16
	currentEncoder = "EncodeUint16"
	currentEndianness = LittleEndian
	b = EncodeUint16(nil, LittleEndian, 0xA932)
	assert(uint16(0xA932), []byte{0x32, 0xA9})

	currentEndianness = BigEndian
	b = EncodeUint16(nil, BigEndian, 0xF243)
	assert(uint16(0xF243), []byte{0xF2, 0x43})

	// sint16
	currentEncoder = "EncodeSint16"
	currentEndianness = LittleEndian
	b = EncodeSint16(nil, LittleEndian, 31234)
	assert(int16(31234), []byte{0x02, 0x7a})

	b = EncodeSint16(nil, LittleEndian, -25043)
	assert(int16(-25043), []byte{0x2d, 0x9e})

	currentEndianness = BigEndian
	b = EncodeSint16(nil, BigEndian, 23984)
	assert(int16(23984), []byte{0x5d, 0xb0})

	b = EncodeSint16(nil, BigEndian, -12398)
	assert(int16(-12398), []byte{0xcf, 0x92})

	// uint32
	currentEncoder = "EncodeUint32"
	currentEndianness = LittleEndian
	b = EncodeUint32(nil, LittleEndian, 0xEF73DAB8)
	assert(uint32(0xEF73DAB8), []byte{0xB8, 0xDA, 0x73, 0xEF})

	currentEndianness = BigEndian
	b = EncodeUint32(nil, BigEndian, 0x18B3EF73)
	assert(uint32(0x18B3EF73), []byte{0x18, 0xB3, 0xEF, 0x73})

	// sint32
	currentEncoder = "EncodeSint32"
	currentEndianness = LittleEndian
	b = EncodeSint32(nil, LittleEndian, 234763279)
	assert(int32(234763279), []byte{0x0f, 0x34, 0xfe, 0x0d})

	b = EncodeSint32(nil, LittleEndian, -32776438)
	assert(int32(-32776438), []byte{0x0a, 0xdf, 0x0b, 0xfe})

	currentEndianness = BigEndian
	b = EncodeSint32(nil, BigEndian, 974326637)
	assert(int32(974326637), []byte{0x3a, 0x13, 0x0b, 0x6d})

	b = EncodeSint32(nil, BigEndian, -1283723832)
	assert(int32(-1283723832), []byte{0xb3, 0x7b, 0xed, 0xc8})

	// uint64
	currentEncoder = "EncodeUint64"
	currentEndianness = LittleEndian
	b = EncodeUint64(nil, LittleEndian, 0xEF73DAB818B3EF73)
	assert(uint64(0xEF73DAB818B3EF73), []byte{0x73, 0xEF, 0xB3, 0x18, 0xB8, 0xDA, 0x73, 0xEF})

	currentEndianness = BigEndian
	b = EncodeUint64(nil, BigEndian, 0x18B3EF73EF73DAB8)
	assert(uint64(0x18B3EF73EF73DAB8), []byte{0x18, 0xB3, 0xEF, 0x73, 0xEF, 0x73, 0xDA, 0xB8})

	// sint64
	currentEncoder = "EncodeSint64"
	currentEndianness = LittleEndian
	b = EncodeSint64(nil, LittleEndian, 2389749234789324833)
	assert(int64(2389749234789324833), []byte{0x21, 0x18, 0xe1, 0x81, 0x44, 0x18, 0x2a, 0x21})

	b = EncodeSint64(nil, LittleEndian, -984593745873495847)
	assert(int64(-984593745873495847), []byte{0xd9, 0x98, 0x23, 0x69, 0x34, 0x05, 0x56, 0xf2})

	currentEndianness = BigEndian
	b = EncodeSint64(nil, BigEndian, 1827367672434837284)
	assert(int64(1827367672434837284), []byte{0x19, 0x5c, 0x1d, 0x37, 0x5d, 0x82, 0xc7, 0x24})

	b = EncodeSint64(nil, BigEndian, -3746765475475764754)
	assert(int64(-3746765475475764754), []byte{0xcc, 0x00, 0xd0, 0xa6, 0xb9, 0x91, 0xc9, 0xee})

	// string
	currentEncoder = "EncodeString"
	currentEndianness = -1
	b = EncodeString(nil, "aboba")
	assert("aboba", []byte{'a', 'b', 'o', 'b', 'a', 0})

	b = EncodeString(nil, "")
	assert("", []byte{0})

	// fit type
	currentEncoder = "EncodeFitType"
	currentEndianness = LittleEndian
	b = EncodeFitType(nil, LittleEndian, FIT_TYPE_FLOAT64)
	assert(FIT_TYPE_FLOAT64, []byte{uint8(FIT_TYPE_FLOAT64)})

	currentEndianness = BigEndian
	b = EncodeFitType(nil, BigEndian, FIT_TYPE_UINT64)
	assert(FIT_TYPE_UINT64, []byte{uint8(FIT_TYPE_UINT64)})
}
