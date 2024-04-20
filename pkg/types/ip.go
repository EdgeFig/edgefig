package types

import (
	"net/netip"
)

// IPAddressType differentiates between iv4 and ipv6 addresses
type IPAddressType string

const (
	// IPAddressTypeV4 ipv4 addresses
	IPAddressTypeV4 IPAddressType = "ipv4"
	// IPAddressTypeV6 ipv6 addresses
	IPAddressTypeV6 IPAddressType = "ipv6"
)

// AddressPort is the address config block for a NAT rule
type AddressPort struct {
	Address netip.Addr `yaml:"address" edge:"address"` // @TODO I think this can also be 1.2.3.4-5.6.7.8 style, so maybe needs to be string or custom type (at least for NAT)
	Port    uint16     `yaml:"port" edge:"port,omitempty"`
}

// NetworkPort is the address config block that accepts networks, not just a single ip address
type NetworkPort struct {
	Address netip.Prefix `yaml:"address" edge:"address"` // @TODO I think this can also be 1.2.3.4-5.6.7.8 style, so maybe needs to be string or custom type (at least for NAT)
	Port    uint16       `yaml:"port" edge:"port,omitempty"`
}
