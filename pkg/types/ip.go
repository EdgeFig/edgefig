package types

import (
	"fmt"
	"net/netip"
	"strings"

	"gopkg.in/yaml.v3"
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
	Address netip.Addr   `yaml:"address" edge:"address,omitempty"`
	Prefix  netip.Prefix `yaml:"prefix" edge:"address,omitempty"`
	Range   AddressRange `yaml:"range" edge:"address,omitempty"`
	Group   AddressGroup `yaml:",inline" edge:"group,omitempty"`
	Port    uint16       `yaml:"port" edge:"port,omitempty"`
}

// AddressRange enables address ranges like 10.0.0.1-10.0.0.5
type AddressRange struct {
	Start netip.Addr
	End   netip.Addr
}

// PortGroup used to specify an IP by group
type AddressGroup struct {
	AddressGroup string `yaml:"address-group" edge:"address-group"`
}

// UnmarshalYAML unmarshals the range 10.0.0.1-10.0.0.5 to the struct representing the Start/End
func (a *AddressRange) UnmarshalYAML(value *yaml.Node) error {
	var v string
	if err := value.Decode(&v); err != nil {
		return err
	}

	parts := strings.Split(v, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range format: expected '<start_ip>-<end_ip>'")
	}

	startIP, err := netip.ParseAddr(parts[0])
	if err != nil {
		return fmt.Errorf("invalid start IP address: %v", err)
	}

	endIP, err := netip.ParseAddr(parts[1])
	if err != nil {
		return fmt.Errorf("invalid end IP address: %v", err)
	}

	a.Start = startIP
	a.End = endIP
	return nil
}
