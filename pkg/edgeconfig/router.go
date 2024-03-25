package edgeconfig

import (
	"net/netip"

	"github.com/cmmarslender/edgefig/pkg/types"
)

// Router is the top level config that applies to EdgeRouters
type Router struct {
	Firewall   Firewall       `edge:"firewall"`
	Interfaces Interfaces     `edge:"interfaces"`
	Service    RouterServices `edge:"service"`
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

type Interfaces struct {
	Interfaces []Interface `edge:"{{ .Type }} {{ .Name }}"`
}

// InterfaceType lists known types for an interface
type InterfaceType string

const (
	// InterfaceTypeEthernet "ethernet"
	InterfaceTypeEthernet InterfaceType = "ethernet"
	// InterfaceTypeLoopback "loopback"
	InterfaceTypeLoopback InterfaceType = "loopback"
)

// Interface is a single interface on the router
type Interface struct {
	Type        InterfaceType // @TODO add validation that this is always present, or add default
	Name        string
	Enable      EnableDisable  `edge:""`
	Description string         `edge:"description,omitempty"`
	Address     []netip.Prefix `edge:"address"` // @TODO this might be a weird tag - this just gets repeated
	Duplex      string         `edge:"duplex"`  // @TODO this is likely uint w/ default string auto
	Speed       string         `edge:"speed"`   // @TODO this is likely uint w/ default string auto
	MTU         uint16         `edge:"mtu"`     // @TODO this is likely uint w/ default string auto
}

// RouterServices Available services on the router
type RouterServices struct {
	DHCPServer DHCPServer  `edge:"dhcp-server"`
	GUI        GUIService  `edge:"gui"`
	NAT        NatService  `edge:"nat"`
	SSH        SSHService  `edge:"ssh"`
	UNMS       UNMSService `edge:"unms"`
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

// GUIService Settings for GUI
type GUIService struct {
	HTTPPort     uint16        `edge:"http-port"`
	HTTPSPort    uint16        `edge:"https-port"`
	OlderCiphers EnableDisable `edge:"older-ciphers"`
}

// NatService NAT settings
type NatService struct {
	Rules []NatRule `edge:"rule {{ .Index }}"`
}

// NatRule a single NAT rule
type NatRule struct {
	Name              string           `edge:"description"`
	Type              types.NATType    `edge:"type"`
	InboundInterface  string           `edge:"inbound-interface,omitempty"`
	OutboundInterface string           `edge:"outbound-interface,omitempty"`
	Protocol          types.Protocol   `edge:"protocol,omitempty"`
	Log               EnableDisable    `edge:"log"`
	OutsideAddress    types.NATAddress `edge:"destination,omitempty"`
	InsideAddress     types.NATAddress `edge:"inside-address,omitempty"`
}

// SSHService settings for ssh
type SSHService struct {
	Port            uint16 `edge:"port"`
	ProtocolVersion string `edge:"protocol-version"`
}

// UNMSService settings for UNMS (only supported here to keep an empty block in config)
type UNMSService struct{}
