package maybeio

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaybeWriter(t *testing.T) {
	r := require.New(t)

	buffer := new(bytes.Buffer)
	mwr := NewWriter(buffer)

	n, err := mwr.Write([]byte("TEST STRING"))
	r.Equal(n, len("TEST STRING"))
	r.NoError(err)
	r.NoError(mwr.Error())
	r.NoError(mwr.Error(), "no error should happen without writes")
}

var errFailedWrite = errors.New("failed to write")

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) {
	return 0, errFailedWrite
}

func TestMaybeWriterError(t *testing.T) {
	r := require.New(t)

	mwr := NewWriter(failWriter{})

	n, err := mwr.Write([]byte("TEST STRING"))
	r.Equal(n, 0)
	r.ErrorIs(err, errFailedWrite)
	r.ErrorIs(mwr.Error(), errFailedWrite)
	r.ErrorIs(mwr.Error(), errFailedWrite, "error should be retrievable multiple times")

	n, err = mwr.Write([]byte("ANOTHER TEST STRING"))
	r.Equal(n, 0)
	r.ErrorIs(err, errFailedWrite, "error shouldn't change after failing once")
	r.ErrorIs(mwr.Error(), errFailedWrite)
}
