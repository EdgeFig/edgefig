package config

import (
	"fmt"
	"net/netip"

	"github.com/cmmarslender/edgefig/pkg/types"
)

// Config is the top level config container
type Config struct {
	Routers []Router `yaml:"routers"`
	VLANs   []VLAN   `yaml:"vlans"`
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
	Name string `yaml:"name"`
	Connection
	Interfaces map[string]RouterInterface `yaml:"interfaces"`
	Firewall   Firewall                   `yaml:"firewall"`
	BGP        []BGP                      `yaml:"bgp"`
	Routes     []StaticRoute              `yaml:"routes"`
	DHCP       []DHCP                     `yaml:"dhcp"`
	DNS        DNS                        `yaml:"dns"`
	NAT        []NAT                      `yaml:"nat"`
	Users      []User                     `yaml:"users"`
}

// RouterInterface is a single physical interface on a router
type RouterInterface struct {
	Name      string         `yaml:"name"`
	Addresses []netip.Prefix `yaml:"addresses"`
	MTU       uint16         `yaml:"mtu"`
	Speed     uint32         `yaml:"speed"`
	Duplex    string         `yaml:"duplex"`
	IPv6      IPv6Config     `yaml:"ipv6"`
	VLANs     []string       `yaml:"vlans"`
}

// IPv6Config configures ipv6 for this interface
type IPv6Config struct {
	Nameserver netip.Addr   `yaml:"nameserver"`
	Prefixes   []IPv6Prefix `yaml:"prefixes"`
}

// IPv6Prefix a single advertised prefix on this interface
type IPv6Prefix struct {
	Prefix     netip.Prefix `yaml:"prefix"`
	Autonomous bool         `yaml:"autonomous"`
}

// Firewall config for the router firewall
type Firewall struct {
	Groups FirewallGroups `yaml:"groups"`
	Zones  []FirewallZone `yaml:"zones"`
}

// FirewallGroups groups of hosts for use in firewall rules
type FirewallGroups struct {
	AddressGroups []AddressGroup `yaml:"address-groups"`
}

// AddressGroup is an address group for the firewall
type AddressGroup struct {
	// @TODO check if this can accept port and all the types of "address"
	// @TODO Probably should accept mutiple addresses/prefixes/etc? (point of a group right?)
	Name              string `yaml:"name"`
	types.AddressPort `yaml:",inline"`
	Description       string `yaml:"description"`
}

// FirewallZone a single firewall zone
type FirewallZone struct {
	Name          string              `yaml:"name"`
	IPType        types.IPAddressType `yaml:"ip-type"`
	DefaultAction string              `yaml:"default-action"`
	Description   string              `yaml:"description"`
	In            []string            `yaml:"in"`
	Out           []string            `yaml:"out"`
	Local         []string            `yaml:"local"`
	Rules         []FirewallRule      `yaml:"rules"`
}

// FirewallRule is a single rule within a firewall zone
type FirewallRule struct {
	Action      string              `yaml:"action"`
	Description string              `yaml:"description"`
	Destination types.AddressPort   `yaml:"destination"`
	Source      types.AddressPort   `yaml:"source"`
	Log         types.EnableDisable `yaml:"log"`
	Protocol    types.Protocol      `yaml:"protocol"`
	Established types.EnableDisable `yaml:"established"`
	Invalid     types.EnableDisable `yaml:"invalid"`
	New         types.EnableDisable `yaml:"new"`
	Related     types.EnableDisable `yaml:"related"`
}

// BGP Defines a single BGP configuration for an AS
type BGP struct {
	ASN           uint32         `yaml:"asn"`
	RouterID      netip.Addr     `yaml:"router-id"`
	Peers         []BGPPeer      `yaml:"peers"`
	Announcements []netip.Prefix `yaml:"announcements"`
}

// BGPPeer is a peer and its configuration for a given BGP session
type BGPPeer struct {
	IP              netip.Addr `yaml:"ip"`
	SourceIP        netip.Addr `yaml:"source-ip"`
	ASN             uint32     `yaml:"asn"`
	AnnounceDefault bool       `yaml:"announce-default"`
}

// StaticRoute is a statically configured route in the router
type StaticRoute struct {
	Description string       `yaml:"description"`
	Route       netip.Prefix `yaml:"route"`
	NextHop     netip.Addr   `yaml:"next-hop"`
	Distance    uint8        `yaml:"distance"`
	Interface   string       `yaml:"interface"`
}

// DHCP is a single DHCP config for a single subnet
type DHCP struct {
	Name            string            `yaml:"name"`
	Authoritative   bool              `yaml:"authoritative"`
	Subnet          netip.Prefix      `yaml:"subnet"`
	Router          netip.Addr        `yaml:"router"`
	Start           netip.Addr        `yaml:"start"`
	Stop            netip.Addr        `yaml:"stop"`
	Lease           uint64            `yaml:"lease"`
	DNS             []netip.Addr      `yaml:"dns"`
	Domain          string            `yaml:"domain"`
	UnifiController string            `yaml:"unifi-controller"`
	Reservations    []DHCPReservation `yaml:"reservations"`
}

// DHCPReservation is a reserved IP by MAC address for a DHCP server
type DHCPReservation struct {
	Name string     `yaml:"name"`
	MAC  string     `yaml:"mac"`
	IP   netip.Addr `yaml:"ip"`
}

// DNS Config for the router
type DNS struct {
	Forwarding DNSForwarding `yaml:"forwarding"`
}

// DNSForwarding is the settings when using dns forwarding
type DNSForwarding struct {
	CacheSize   uint16       `yaml:"cache-size"`
	ListenOn    []string     `yaml:"listen-on"`
	Nameservers []netip.Addr `yaml:"nameservers"`
}

// NAT configures NAT rules in a router
type NAT struct {
	Name              string            `yaml:"name"`
	Type              types.NATType     `yaml:"type"`
	InboundInterface  string            `yaml:"inbound_interface"`
	OutboundInterface string            `yaml:"outbound_interface"`
	Protocol          types.Protocol    `yaml:"protocol"`
	Log               bool              `yaml:"log"`
	InsideAddress     types.AddressPort `yaml:"inside_address"`
	OutsideAddress    types.AddressPort `yaml:"outside_address"`
}

// Switch is the top level config for a single switch
type Switch struct{}

// VLAN defines a single shared VLAN configuration
type VLAN struct {
	Name    string       `yaml:"name"`
	ID      uint16       `yaml:"id"`
	Address netip.Prefix `yaml:"address"`
	MTU     uint16       `yaml:"mtu"`
}

// GetVLANByName returns a VLAN by its name attribute
func (c *Config) GetVLANByName(name string) (VLAN, error) {
	for _, vlan := range c.VLANs {
		if vlan.Name == name {
			return vlan, nil
		}
	}

	return VLAN{}, fmt.Errorf("could not find requested VLAN %s in config", name)
}
