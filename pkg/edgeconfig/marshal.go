package edgeconfig

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

// EdgeMarshaller is an interface that indicates customized support to marshal to the edgeconfig format
type EdgeMarshaller interface {
	// MarshalEdge does the actual marshalling to edgeconfig format
	MarshalEdge() ([]byte, error)
}

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
			tag, omitEmpty := parseEdgeTag(structField.Tag.Get("edge"))

			if tag == " " {
				continue // Skip fields without 'edge' tag
			}

			// Handle recursive types
			switch field.Kind() {
			case reflect.Struct:
				var err error
				tag, err = templateTagValues(tag, field, 0)
				if err != nil {
					return err
				}
				fallthrough
			case reflect.Map:
				specificType := field.Type().String()
				switch specificType {
				case "netip.Addr":
					val, err := formatValue(field, omitEmpty)
					if err != nil {
						return err
					}
					buffer.WriteString(fmt.Sprintf("%s%s%s\n", strings.Repeat(" ", depth), tag, val))
				default:
					if omitEmpty && field.IsZero() {
						return nil
					}
					buffer.WriteString(strings.Repeat(" ", depth) + tag + "{\n")
					err := marshalValue(buffer, field, depth+4)
					if err != nil {
						return err
					}
					buffer.WriteString(strings.Repeat(" ", depth) + "}\n")
				}
			case reflect.Slice:
				for i := 0; i < field.Len(); i++ {
					sliceElement := field.Index(i)

					tag, err := templateTagValues(tag, sliceElement, i)
					if err != nil {
						return err
					}
					buffer.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", depth), tag))

					typeStr := sliceElement.Type().String()
					switch typeStr {
					case "netip.Prefix", "netip.Addr":
						val, err := formatValue(sliceElement, omitEmpty)
						if err != nil {
							return err
						}
						buffer.WriteString(val)
					default:
						buffer.WriteString("{\n")
						err = marshalValue(buffer, sliceElement, depth+4)
						buffer.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", depth), "}"))
						if err != nil {
							return err
						}
					}

					buffer.WriteString("\n")
				}
			default:
				// Directly marshal field with value.
				fieldValue, err := formatValue(field, omitEmpty)
				if err != nil {
					return err
				}
				if !omitEmpty || (fieldValue != "") {
					buffer.WriteString(fmt.Sprintf("%s%s%s\n", strings.Repeat(" ", depth), tag, fieldValue))
				}
			}
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			keyValue, err := formatValue(key, false) // @TODO omit empty support here
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
				valueValue, err := formatValue(value, false) // @TODO omitempty support here
				if err != nil {
					return err
				}
				buffer.WriteString(fmt.Sprintf("%s%s%s\n", strings.Repeat(" ", depth), keyValue, valueValue))
			}
		}
	}

	return nil
}

// formatValue converts field values into their EdgeOS string representation.
func formatValue(val reflect.Value, omitEmpty bool) (string, error) {
	// Check if the type implements the CustomMarshaller interface.
	// This approach works if you have a known interface.
	if marshaller, ok := val.Interface().(EdgeMarshaller); ok {
		data, err := marshaller.MarshalEdge()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	switch val.Kind() {
	case reflect.String:
		strval := val.String()
		if strings.Contains(strval, " ") {
			// Needs quotes
			strval = fmt.Sprintf("\"%s\"", val.String())
		}
		return strval, nil
	case reflect.Bool:
		specificBool := val.Type().Name()
		switch specificBool {
		case "EnableDisable":
			if val.Bool() {
				return "enable", nil
			}
			return "disable", nil
		case "DisableProp":
			if val.Bool() {
				return "disable", nil
			}
			return "", nil
		default:
			if val.Bool() {
				return "true", nil
			}
			return "false", nil
		}
	case reflect.Slice:
		if omitEmpty && val.IsZero() {
			return "", nil
		}
		specificSlice := val.Type().String()
		return "", fmt.Errorf("should not be formatting slices directly: %s", specificSlice)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.Float32, reflect.Float64:
		if omitEmpty && val.IsZero() {
			return "", nil
		}
		fallthrough
	default:
		return fmt.Sprintf("%v", val.Interface()), nil
	}
}

// Tags a tag and parses/fills any go templates
func templateTagValues(tag string, element reflect.Value, index int) (string, error) {
	tmpl, err := template.New("tag").Parse(tag)
	if err != nil {
		return "", err
	}
	var executedTag bytes.Buffer

	data := make(map[string]interface{})
	for i := 0; i < element.NumField(); i++ {
		field := element.Type().Field(i)
		// If field is public
		if field.PkgPath == "" {
			data[field.Name] = element.Field(i).Interface()
		}
	}
	data["Index"] = index

	err = tmpl.Execute(&executedTag, data)
	if err != nil {
		return "", err
	}

	return executedTag.String(), nil
}

func parseEdgeTag(tag string) (string, bool) {
	omitEmpty := false
	splits := strings.Split(tag, ",")
	tag = splits[0]

	if len(splits) > 1 {
		for _, value := range splits[1:] {
			if value == "omitempty" {
				omitEmpty = true
			}
		}
	}

	if tag == "." {
		tag = ""
	} else {
		tag = fmt.Sprintf("%s ", tag)
	}

	return tag, omitEmpty
}
