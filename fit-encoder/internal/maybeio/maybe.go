// Copyright 2021 Artem Mikheev

package maybeio

import "io"

// Writer is an implementation of the io.Writer interface which
// writes if no error has yet occurred, otherwise returns the error.
type Writer struct {
	io.Writer
	err error
}

// NewWriter returns a new maybe.Writer with the given underlying writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w}
}

// Write either writes the byte slice to the underlying
// writer, or returns the stored error
func (mwr *Writer) Write(b []byte) (n int, err error) {
	if mwr.err != nil {
		return 0, mwr.err
	}

	n, err = mwr.Writer.Write(b)
	mwr.err = err
	return
}

// Error returns the error stored in a MaybeWriter
func (mwr *Writer) Error() error {
	return mwr.err
}
