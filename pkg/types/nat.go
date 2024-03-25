package types

import (
	"net/netip"
)

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

// NATAddress is the address config block for a NAT rule
type NATAddress struct {
	Address netip.Addr `yaml:"address" edge:"address"` // @TODO I think this can also be 1.2.3.4-5.6.7.8 style, so maybe needs to be string or custom type
	Port    uint16     `yaml:"port" edge:"port,omitempty"`
}
