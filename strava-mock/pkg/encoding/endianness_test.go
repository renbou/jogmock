package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndianness(t *testing.T) {
	a := assert.New(t)

	fakeEndian := Endianness(123)

	a.True(LittleEndian.IsKnown())
	a.True(BigEndian.IsKnown())
	a.False(fakeEndian.IsKnown())

	a.Equal(LittleEndian.String(), "LittleEndian")
	a.Equal(BigEndian.String(), "BigEndian")
	a.Equal(fakeEndian.String(), "Endianness(123)")
}
