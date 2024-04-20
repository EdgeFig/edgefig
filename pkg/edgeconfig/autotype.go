package edgeconfig

import (
	"fmt"
)

// AutoString is either the string value, or "auto"
type AutoString string

// AutoUint16 is either the string "auto" or the value of the uint
type AutoUint16 uint16

// AutoUint32 is either the string "auto" or the value of the uint
type AutoUint32 uint32

// AutoType returns "auto" if the type is empty, or else the value
func AutoType[T string | uint16 | uint32](val T) string {
	if isEmpty(val) {
		return "auto"
	}

	// Convert the non-empty value to a string using fmt.Sprintf
	return fmt.Sprintf("%v", val)
}

// isEmpty checks if the given value of type T (string or uint16) is "empty".
func isEmpty[T string | uint16 | uint32](val T) bool {
	// Type switch to handle different types differently
	switch v := any(val).(type) {
	case string:
		return v == ""
	case uint16:
		return v == 0
	case uint32:
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

// MarshalEdge handles marshalling down to the edgeconfig format
func (au AutoUint16) MarshalEdge() ([]byte, error) {
	val := AutoType[uint16](uint16(au))
	return []byte(val), nil
}

// MarshalEdge handles marshalling down to the edgeconfig format
func (au AutoUint32) MarshalEdge() ([]byte, error) {
	val := AutoType[uint32](uint32(au))
	return []byte(val), nil
}

// MarshalEdgeWithDepth not used for AutoString
func (as AutoString) MarshalEdgeWithDepth(depth int) ([]byte, error) {
	return nil, fmt.Errorf("marshaledgewithdepth not implemented for AutoString")
}

// MarshalEdgeWithDepth not used for AutoString
func (au AutoUint16) MarshalEdgeWithDepth(depth int) ([]byte, error) {
	return nil, fmt.Errorf("marshaledgewithdepth not implemented for AutoUint16")
}

// MarshalEdgeWithDepth not used for AutoString
func (au AutoUint32) MarshalEdgeWithDepth(depth int) ([]byte, error) {
	return nil, fmt.Errorf("marshaledgewithdepth not implemented for AutoUint32")
}
