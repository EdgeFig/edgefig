package config

import (
	"net"
	"net/netip"
)

// Config is the top level config container
type Config struct {
	Routers []Router `yaml:"routers"`
}

// User defines a common struct that represents a user across routers, switches, etc
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Role     string `yaml:"role"`
}

// Router is the top level config for a single router
type Router struct {
	Name       string                     `yaml:"name"`
	Interfaces map[string]RouterInterface `yaml:"interfaces"`
	DHCP       []DHCP                     `yaml:"dhcp"`
	NAT        []NAT                      `yaml:"nat"`
	Users      []User                     `yaml:"users"`
}

// RouterInterface is a single physical interface on a router
type RouterInterface struct {
	Name      string         `yaml:"name"`
	Addresses []netip.Prefix `yaml:"addresses"`
	MTU       uint16         `yaml:"mtu"`
}

// DHCP is a single DHCP config for a single subnet
type DHCP struct {
	Name          string       `yaml:"name"`
	Authoritative bool         `yaml:"authoritative"`
	Subnet        netip.Prefix `yaml:"subnet"`
	Router        net.IP       `yaml:"router"`
	Start         net.IP       `yaml:"start"`
	Stop          net.IP       `yaml:"stop"`
	Lease         uint64       `yaml:"lease"`
	DNS           []net.IP     `yaml:"dns"`
}

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

// Protocol defines the known protocol types
type Protocol string

const (
	// ProtocolAll for all protocols
	ProtocolAll Protocol = "all"
	// ProtocolTCP for tcp only
	ProtocolTCP Protocol = "tcp"
	// ProtocolUDP for udp only
	ProtocolUDP Protocol = "udp"
)

// NAT configures NAT rules in a router
type NAT struct {
	Name              string     `yaml:"name"`
	Type              NATType    `yaml:"type"`
	InboundInterface  string     `yaml:"inbound_interface"`
	OutboundInterface string     `yaml:"outbound_interface"`
	Protocol          Protocol   `yaml:"protocol"`
	Log               bool       `yaml:"log"`
	InsideAddress     NATAddress `yaml:"inside_address"`
	OutsideAddress    NATAddress `yaml:"outside_address"`
}

// NATAddress is the address config block for a NAT rule
type NATAddress struct {
	Address net.IP `yaml:"address"` // @TODO I think this can also be 1.2.3.4-5.6.7.8 style, so maybe needs to be string or custom type
	Port    uint16 `yaml:"port"`
}

// Switch is the top level config for a single switch
type Switch struct{}
