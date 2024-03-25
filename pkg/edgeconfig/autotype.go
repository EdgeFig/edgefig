package edgeconfig

import (
	"fmt"
)

// AutoString is either the string value, or "auto"
type AutoString string

// AutoUint16 is either the string "auto" or the value of the uint
type AutoUint16 uint16

// AutoType returns "auto" if the type is empty, or else the value
func AutoType[T string | uint16](val T) string {
	if isEmpty(val) {
		return "auto"
	}

	// Convert the non-empty value to a string using fmt.Sprintf
	return fmt.Sprintf("%v", val)
}

// isEmpty checks if the given value of type T (string or uint16) is "empty".
func isEmpty[T string | uint16](val T) bool {
	// Type switch to handle different types differently
	switch v := any(val).(type) {
	case string:
		return v == ""
	case uint16:
		return v == 0
	default:
		// This should never happen due to the type constraint on T
		return false
	}
}

// MarshalEdge handles marshalling down to the edgeconfig format
func (as AutoString) MarshalEdge() ([]byte, error) {
	val := AutoType[string](string(as))
	return []byte(val), nil
}
