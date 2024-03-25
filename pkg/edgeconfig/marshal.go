package edgeconfig

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
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

	kind := val.Kind()
	switch kind {
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
				var err error
				tag, err = templateTagValues(tag, field)
				if err != nil {
					return err
				}
				fallthrough
			case reflect.Map:
				specificType := field.Type().String()
				switch specificType {
				case "netip.Addr":
					val, err := formatValue(field)
					if err != nil {
						return err
					}
					buffer.WriteString(fmt.Sprintf("%s%s %s\n", strings.Repeat(" ", depth), tag, val))
				default:
					buffer.WriteString(strings.Repeat(" ", depth) + tag + " {\n")
					err := marshalValue(buffer, field, depth+4)
					if err != nil {
						return err
					}
					buffer.WriteString(strings.Repeat(" ", depth) + "}\n")
				}
			case reflect.Slice:
				for i := 0; i < field.Len(); i++ {
					sliceElement := field.Index(i)

					tag, err := templateTagValues(tag, sliceElement)
					if err != nil {
						return err
					}
					buffer.WriteString(fmt.Sprintf("%s%s ", strings.Repeat(" ", depth), tag))

					typeStr := sliceElement.Type().String()
					switch typeStr {
					case "edgeconfig.DHCPNetwork":
						fallthrough
					case "edgeconfig.DHCPSubnet":
						buffer.WriteString("{\n")
						err = marshalValue(buffer, sliceElement, depth+4)
						buffer.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", depth), "}"))
						if err != nil {
							return err
						}
					default:
						val, err := formatValue(sliceElement)
						if err != nil {
							return err
						}
						buffer.WriteString(val)
					}

					buffer.WriteString("\n")
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
	switch val.Kind() {
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		specificBool := val.Type().Name()
		switch specificBool {
		case "EnableDisable":
			if val.Bool() {
				return "enable", nil
			}
			return "disable", nil
		default:
			if val.Bool() {
				return "true", nil
			}
			return "false", nil
		}
	case reflect.Slice:
		specificSlice := val.Type().String()
		return "", fmt.Errorf("should not be formatting slices directly: %s", specificSlice)
	default:
		return fmt.Sprintf("%v", val.Interface()), nil
	}
}

// Tags a tag and parses/fills any go templates
func templateTagValues(tag string, sliceElement reflect.Value) (string, error) {
	tmpl, err := template.New("tag").Parse(tag)
	if err != nil {
		return "", err
	}
	var executedTag bytes.Buffer
	iface := sliceElement.Interface()
	err = tmpl.Execute(&executedTag, iface)
	if err != nil {
		return "", err
	}

	return executedTag.String(), nil
}
