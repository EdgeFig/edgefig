package edgeconfig

import (
	"fmt"
)

// Validate ensures the config is minimally valid
// primarily concerned with ensuring required fields are set, not that every detail is correct
func (rc *Router) Validate() error {
	// @TODO ensure at least one user

	// Can't have addresses and DHCP set
	for _, interf := range rc.Interfaces.Interfaces {
		if len(interf.Address) > 0 && interf.AddressDHCP != "" {
			return fmt.Errorf("interface %s cannot have AddressDHCP and Address both set", interf.Name)
		}
	}

	return nil
}
