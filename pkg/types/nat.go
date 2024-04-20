package types

// NATType is the known types of NAT that can be set in the config
type NATType string

const (
	// NATTypeDestination Destination NAT rule
	NATTypeDestination NATType = "destination"
	// NATTypeSource Source NAT rule
	NATTypeSource NATType = "source"
	// NATTypeMasquerade Masquerade rule
	NATTypeMasquerade NATType = "masquerade"
)
