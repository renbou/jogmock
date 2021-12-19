// Copyright 2021 Artem Mikheev

package maybe

import "io"

// MaybeWriter is an implementation of the io.Writer interface which
// writes if no error has yet occurred, otherwise returns the error.
type MaybeWriter struct {
	Writer io.Writer
	err    error
}

func (mwr *MaybeWriter) Write(b []byte) (n int, err error) {
	if mwr.err != nil {
		return 0, mwr.err
	}

	n, err = mwr.Writer.Write(b)
	mwr.err = err
	return
}

// Error returns the error stored in a MaybeWriter
func (mwr *MaybeWriter) Error() error {
	return mwr.err
}
