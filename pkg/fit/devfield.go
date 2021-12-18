// Copyright 2021 Artem Mikheev

package fit

import (
	"io"

	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit/types"
)

type DevFieldDefinition struct {
	FieldNum types.FitUint8
	Size     types.FitUint8
	DevIndex types.FitUint8
}

func (devFieldDef *DevFieldDefinition) Encode(wr io.Writer, endianness encoding.Endianness) error {
	if err := devFieldDef.FieldNum.Encode(wr, endianness); err != nil {
		return err
	}
	if err := devFieldDef.Size.Encode(wr, endianness); err != nil {
		return err
	}
	if err := devFieldDef.DevIndex.Encode(wr, endianness); err != nil {
		return err
	}
	return nil
}
