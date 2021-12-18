package types

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	fitTesting "github.com/renbou/strava-keker/internal/testing"
	"github.com/renbou/strava-keker/pkg/encoding"
)

// assertEncoded asserts that the value actually encoded into what we expected,
// otherwise it prints a pretty informative message. For proper debug output,
// make sure you explicitly specify value's int type (ex: int8(-4))
func assertEncoded(t *testing.T, encoder *encoding.Encoder, value interface{}, expected, actual []byte) {
	var valueRepr string
	if strings.HasPrefix(reflect.TypeOf(value).Name(), "uint") {
		valueRepr = fmt.Sprintf("0x%x", value)
	} else {
		valueRepr = fmt.Sprintf("%v", value)
	}

	if err := fitTesting.AssertEqual(expected, actual); err != nil {
		t.Fatalf("Invalid encoding of %s with %s: %v",
			valueRepr, encoder, err)
	}
}

func TestEncoding(t *testing.T) {
	var (
		buffer       *bytes.Buffer
		encoder      *encoding.Encoder
		currentValue interface{}
	)

	// helper for assertions with better debug output
	assert := func(expected []byte) {
		assertEncoded(t, encoder, currentValue, expected, buffer.Bytes())
	}

	resetEncoder := func(endianness encoding.Endianness) {
		buffer = new(bytes.Buffer)
		encoder = encoding.NewEncoder(buffer, endianness)
	}

	encode := func(value encoding.EndianEncoder, endianness encoding.Endianness) {
		resetEncoder(endianness)
		currentValue = value
		err := encoder.Encode(value)
		if err != nil {
			t.Fatalf("Unexpected error during encoding: %v", err)
		}
	}

	// uint8
	encode(FitUint8(0x69), encoding.LittleEndian)
	assert([]byte{0x69})

	resetEncoder(encoding.BigEndian)
	encode(FitUint8(0x42), encoding.BigEndian)
	assert([]byte{0x42})

	// sint8
	encode(FitSint8(100), encoding.LittleEndian)
	assert([]byte{0x64})

	encode(FitSint8(-100), encoding.LittleEndian)
	assert([]byte{0x9c})

	encode(FitSint8(34), encoding.BigEndian)
	assert([]byte{0x22})

	encode(FitSint8(-34), encoding.BigEndian)
	assert([]byte{0xde})

	// enum
	encode(FitEnum(0), encoding.LittleEndian)
	assert([]byte{0})

	encode(FitEnum(19), encoding.BigEndian)
	assert([]byte{0x13})

	// uint16
	encode(FitUint16(0xA932), encoding.LittleEndian)
	assert([]byte{0x32, 0xA9})

	encode(FitUint16(0xF243), encoding.BigEndian)
	assert([]byte{0xF2, 0x43})

	// sint16
	encode(FitSint16(31234), encoding.LittleEndian)
	assert([]byte{0x02, 0x7a})

	encode(FitSint16(-25043), encoding.LittleEndian)
	assert([]byte{0x2d, 0x9e})

	encode(FitSint16(23984), encoding.BigEndian)
	assert([]byte{0x5d, 0xb0})

	encode(FitSint16(-12398), encoding.BigEndian)
	assert([]byte{0xcf, 0x92})

	// uint32
	encode(FitUint32(0xEF73DAB8), encoding.LittleEndian)
	assert([]byte{0xB8, 0xDA, 0x73, 0xEF})

	encode(FitUint32(0x18B3EF73), encoding.BigEndian)
	assert([]byte{0x18, 0xB3, 0xEF, 0x73})

	// sint32
	encode(FitSint32(234763279), encoding.LittleEndian)
	assert([]byte{0x0f, 0x34, 0xfe, 0x0d})

	encode(FitSint32(-32776438), encoding.LittleEndian)
	assert([]byte{0x0a, 0xdf, 0x0b, 0xfe})

	encode(FitSint32(974326637), encoding.BigEndian)
	assert([]byte{0x3a, 0x13, 0x0b, 0x6d})

	encode(FitSint32(-1283723832), encoding.BigEndian)
	assert([]byte{0xb3, 0x7b, 0xed, 0xc8})

	// uint64
	encode(FitUint64(0xEF73DAB818B3EF73), encoding.LittleEndian)
	assert([]byte{0x73, 0xEF, 0xB3, 0x18, 0xB8, 0xDA, 0x73, 0xEF})

	encode(FitUint64(0x18B3EF73EF73DAB8), encoding.BigEndian)
	assert([]byte{0x18, 0xB3, 0xEF, 0x73, 0xEF, 0x73, 0xDA, 0xB8})

	// sint64
	encode(FitSint64(2389749234789324833), encoding.LittleEndian)
	assert([]byte{0x21, 0x18, 0xe1, 0x81, 0x44, 0x18, 0x2a, 0x21})

	encode(FitSint64(-984593745873495847), encoding.LittleEndian)
	assert([]byte{0xd9, 0x98, 0x23, 0x69, 0x34, 0x05, 0x56, 0xf2})

	encode(FitSint64(1827367672434837284), encoding.BigEndian)
	assert([]byte{0x19, 0x5c, 0x1d, 0x37, 0x5d, 0x82, 0xc7, 0x24})

	encode(FitSint64(-3746765475475764754), encoding.BigEndian)
	assert([]byte{0xcc, 0x00, 0xd0, 0xa6, 0xb9, 0x91, 0xc9, 0xee})

	// string
	encode(&FitEncodableString{"aboba", 6}, encoding.LittleEndian)
	assert([]byte{'a', 'b', 'o', 'b', 'a', 0})

	encode(&FitEncodableString{"", 3}, encoding.BigEndian)
	assert([]byte{0, 0, 0})

	// fit type
	encode(FitBaseType(FIT_TYPE_FLOAT64), encoding.LittleEndian)
	assert([]byte{uint8(FIT_TYPE_FLOAT64)})

	encode(FitBaseType(FIT_TYPE_UINT64), encoding.BigEndian)
	assert([]byte{uint8(FIT_TYPE_UINT64)})
}
