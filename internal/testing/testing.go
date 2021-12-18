// Copyright 2021 Artem Mikheev

package testing

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
)

// AssertEqual compares two values and returns an error if they aren't equal
func AssertEqual(expected, actual interface{}) error {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		fmt.Errorf("non-comparable types %s (typeOf expected) and %s (typeOf actual)",
			reflect.TypeOf(expected), reflect.TypeOf(actual))
	}
	switch expected.(type) {
	case []byte:
		expectedBytes := expected.([]byte)
		actualBytes := actual.([]byte)
		if !bytes.Equal(expectedBytes, actualBytes) {
			return fmt.Errorf("expected bytes=%s but actual bytes=%s",
				hex.EncodeToString(expectedBytes), hex.EncodeToString(actualBytes))
		}
	}
	return nil
}
