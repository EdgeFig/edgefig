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
	Name   string       `yaml:"name"`
	Subnet netip.Prefix `yaml:"subnet"`
	Router net.IP       `yaml:"router"`
	Start  net.IP       `yaml:"start"`
	End    net.IP       `yaml:"end"`
	Lease  uint64       `yaml:"lease"`
	DNS    []net.IP     `yaml:"dns"`
}

// Switch is the top level config for a single switch
type Switch struct{}
