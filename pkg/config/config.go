package config

import (
	"net/netip"

	"github.com/cmmarslender/edgefig/pkg/types"
)

// Config is the top level config container
type Config struct {
	Routers []Router `yaml:"routers"`
}

// Connection common details for connecting to devices
type Connection struct {
	IP       netip.Addr `yaml:"ip"`
	Port     uint16     `yaml:"port"`
	Username string     `yaml:"username"`
	Password string     `yaml:"password"`
}

// User defines a common struct that represents a user across routers, switches, etc
type User struct {
	Username string          `yaml:"username"`
	Password string          `yaml:"password"`
	Role     types.UserLevel `yaml:"role"`
}

// Router is the top level config for a single router
type Router struct {
	Name       string                     `yaml:"name"`
	Connection
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
	Router        netip.Addr   `yaml:"router"`
	Start         netip.Addr   `yaml:"start"`
	Stop          netip.Addr   `yaml:"stop"`
	Lease         uint64       `yaml:"lease"`
	DNS           []netip.Addr `yaml:"dns"`
}

// NAT configures NAT rules in a router
type NAT struct {
	Name              string           `yaml:"name"`
	Type              types.NATType    `yaml:"type"`
	InboundInterface  string           `yaml:"inbound_interface"`
	OutboundInterface string           `yaml:"outbound_interface"`
	Protocol          types.Protocol   `yaml:"protocol"`
	Log               bool             `yaml:"log"`
	InsideAddress     types.NATAddress `yaml:"inside_address"`
	OutsideAddress    types.NATAddress `yaml:"outside_address"`
}

// Switch is the top level config for a single switch
type Switch struct{}
