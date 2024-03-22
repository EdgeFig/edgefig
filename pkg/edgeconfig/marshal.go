package edgeconfig

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Marshal takes something in and marshals it according to the edge tags
func Marshal(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	err := marshalValue(&buffer, reflect.ValueOf(v), 0)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func marshalValue(buffer *bytes.Buffer, val reflect.Value, depth int) error {
	// Ensure we're dealing with the base type (in case of pointers).
	val = reflect.Indirect(val)
	if !val.IsValid() {
		return nil // Skip invalid fields (e.g., uninitialized pointers)
	}

	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			structField := t.Field(i)
			tag := structField.Tag.Get("edge")

			if tag == "" {
				continue // Skip fields without 'edge' tag
			}

			// Handle recursive types
			switch field.Kind() {
			case reflect.Struct:
				fallthrough
			case reflect.Map:
				buffer.WriteString(strings.Repeat(" ", depth) + tag + " {\n")
				err := marshalValue(buffer, field, depth+4)
				if err != nil {
					return err
				}
				buffer.WriteString(strings.Repeat(" ", depth) + "}\n")
			case reflect.Slice:
				for i := 0; i < field.Len(); i++ {
					sliceElement := field.Index(i)
					val, err := formatValue(sliceElement)
					if err != nil {
						return err
					}
					buffer.WriteString(fmt.Sprintf("%s%s %s\n", strings.Repeat(" ", depth), tag, val))
				}
			default:
				// Directly marshal field with value.
				fieldValue, err := formatValue(field)
				if err != nil {
					return err
				}
				buffer.WriteString(fmt.Sprintf("%s%s %s\n", strings.Repeat(" ", depth), tag, fieldValue))
			}
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			keyValue, err := formatValue(key)
			if err != nil {
				return err
			}
			value := val.MapIndex(key)

			// Handle recursive types
			switch value.Kind() {
			case reflect.Struct:
				fallthrough
			case reflect.Map:
				buffer.WriteString(strings.Repeat(" ", depth) + keyValue + " {\n")
				err := marshalValue(buffer, value, depth+4)
				if err != nil {
					return err
				}
				buffer.WriteString(strings.Repeat(" ", depth) + "}\n")
			default:
				valueValue, err := formatValue(value)
				if err != nil {
					return err
				}
				buffer.WriteString(fmt.Sprintf("%s%s %s\n", strings.Repeat(" ", depth), keyValue, valueValue))
			}
		}
	}

	return nil
}

// formatValue converts field values into their EdgeOS string representation.
func formatValue(val reflect.Value) (string, error) {
	// Handle any specific types here (e.g., custom types like EnableDisable).
	// This is a simplification. Adjust as needed for your actual types.
	switch val.Kind() {
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		if val.Bool() {
			return "enable", nil
		}
		return "disable", nil
	case reflect.Slice:
		return "", fmt.Errorf("should not be formatting slices directly - slices indicate repeated statements")
	default:
		return fmt.Sprintf("%v", val.Interface()), nil
	}
}
