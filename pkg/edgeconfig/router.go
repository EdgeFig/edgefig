package edgeconfig

import (
	"net/netip"
)

// Router is the top level config that applies to EdgeRouters
type Router struct {
	Firewall   Firewall             `edge:"firewall"`
	Interfaces map[string]Interface `edge:"interface"`
	Service    RouterServices       `edge:"service"`
}

// Firewall is the firewall config for routers
type Firewall struct {
	AllPing              EnableDisable `edge:"all-ping"`
	BroadcastPing        EnableDisable `edge:"broadcast-ping"`
	IPv6ReceiveRedirects EnableDisable `edge:"ipv6-receive-redirects"`
	IPv6SrcRoute         EnableDisable `edge:"ipv6-src-route"`
	LogMartians          EnableDisable `edge:"log-martians"`
	ReceiveRedirects     EnableDisable `edge:"receive-redirects"`
	SendRedirects        EnableDisable `edge:"send-redirects"`
	SourceValidation     EnableDisable `edge:"source-validation"`
	SynCookies           EnableDisable `edge:"syn-cookies"`
}

// Interface is a single interface on the router
type Interface struct {
	Enable  EnableDisable  `edge:""`
	Address []netip.Prefix `edge:"address"` // @TODO this might be a weird tag - this just gets repeated
	Duplex  string         `edge:"duplex"`  // @TODO this is likely uint w/ default string auto
	Speed   string         `edge:"speed"`   // @TODO this is likely uint w/ default string auto
	MTU     uint16         `edge:"mtu"`     // @TODO this is likely uint w/ default string auto
}

// RouterServices Available services on the router
type RouterServices struct {
	DHCPServer DHCPServer `edge:"dhcp-server"`
	GUI        GUI        `edge:"gui"`
}

// DHCPServer information about enabled DHCP servers
type DHCPServer struct {
	Disabled       bool          `edge:"disabled"`
	HostfileUpdate EnableDisable `edge:"hostfile-update"`
	StaticARP      EnableDisable `edge:"static-arp"`
	UseDNSMASQ     EnableDisable `edge:"use-dnsmasq"`
	Networks       []DHCPNetwork `edge:"shared-network-name {{ .Name }}"`
}

// DHCPNetwork is a single network managed by the DHCP server
type DHCPNetwork struct {
	Name          string
	Authoritative EnableDisable `edge:"authoritative"`
	Subnets       []DHCPSubnet  `edge:"subnet {{ .Subnet }}"`
}

// DHCPSubnet is a subnet within a DHCP Network
type DHCPSubnet struct {
	Subnet    netip.Prefix
	Router    netip.Addr    `edge:"default-router"`
	Lease     uint64        `edge:"lease"`
	DNS       []netip.Addr  `edge:"dns-server"`
	StartStop DHCPStartStop `edge:"start {{ .Start }}"`
}

// DHCPStartStop is the start/stop range for DHCP Servers
type DHCPStartStop struct {
	Start netip.Addr
	Stop  netip.Addr `edge:"stop"`
}

// GUI Settings
type GUI struct {
	HTTPPort     uint16        `edge:"http-port"`
	HTTPSPort    uint16        `edge:"https-port"`
	OlderCiphers EnableDisable `edge:"older-ciphers"`
}
