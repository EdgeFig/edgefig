package types

import (
	"fmt"
)

// EnableDisable is a simple boolean wrapper that marshals to "enable" or "disable"
type EnableDisable bool

const (
	// Enable is a quick const to set EnableDisable to Enable
	Enable = EnableDisable(true)
	// Disable is a quick const to set EnableDisable to Disable
	Disable = EnableDisable(false)
)

// UnmarshalYAML customizes the unmarshalling of EnableDisable to accept either true/false or "enable"/"disable"
func (e *EnableDisable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp interface{}
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	switch v := tmp.(type) {
	case bool:
		*e = EnableDisable(v)
	case string:
		switch v {
		case "enable":
			*e = Enable
		case "disable":
			*e = Disable
		default:
			return fmt.Errorf("unknown value for EnableDisable: %s", v)
		}
	default:
		return fmt.Errorf("invalid type for EnableDisable: %T", tmp)
	}

	return nil
}

// DisableProp is similar to EnableDisable, but ONLY shows up if Disabled, and shows nothing if Enabled
// Used on interfaces
type DisableProp bool

const (
	// Disabled is the state for "disabled"
	Disabled DisableProp = true
	// Enabled is the state for "enabled" and nothing will show up in the config
	Enabled DisableProp = false
)

// KeyWhenEnabled shows the key and nothing else only when true
type KeyWhenEnabled bool
