package edgeconfig

import (
	"bytes"
	"fmt"
	"net/netip"
	"reflect"
	"strings"

	"github.com/cmmarslender/edgefig/pkg/types"
)

// Router is the top level config that applies to EdgeRouters
type Router struct {
	Firewall   Firewall        `edge:"firewall"`
	Interfaces Interfaces      `edge:"interfaces"`
	Protocols  RouterProtocols `edge:"protocols,omitempty"`
	Service    RouterServices  `edge:"service"`
	System     RouterSystem    `edge:"system"`
}

// Firewall is the firewall config for routers
type Firewall struct {
	AllPing              types.EnableDisable `edge:"all-ping"`
	BroadcastPing        types.EnableDisable `edge:"broadcast-ping"`
	IPv6ReceiveRedirects types.EnableDisable `edge:"ipv6-receive-redirects"`
	IPv6SrcRoute         types.EnableDisable `edge:"ipv6-src-route"`
	IPSrcRoute           types.EnableDisable `edge:"ip-src-route"`
	LogMartians          types.EnableDisable `edge:"log-martians"`
	ReceiveRedirects     types.EnableDisable `edge:"receive-redirects"`
	SendRedirects        types.EnableDisable `edge:"send-redirects"`
	SourceValidation     types.EnableDisable `edge:"source-validation"`
	SynCookies           types.EnableDisable `edge:"syn-cookies"`
	Group                FirewallGroup       `edge:"group"`
	Zones                []FirewallZone      `edge:"{{ .NamePrefix }}name {{ .Name }}"`
}

// FirewallGroup groups for the firewall
type FirewallGroup struct{}

// FirewallZone is a specific zone for the firewall
type FirewallZone struct {
	NamePrefix    string
	Name          string
	DefaultAction string         `edge:"default-action"`
	Description   string         `edge:"description"`
	Rules         []FirewallRule `edge:"rule {{ .Count }}"`
}

// FirewallRule is a single rule within a firewall zone
type FirewallRule struct {
	Action      string              `edge:"action"`
	Description string              `edge:"description,omitempty"`
	Destination types.NetworkPort   `edge:"destination,omitempty"`
	Log         types.EnableDisable `edge:"log"`
	Protocol    types.Protocol      `edge:"protocol,omitempty"`
	State       FirewallRuleState   `edge:"state,omitempty"`
}

// FirewallRuleState connection state settings
type FirewallRuleState struct {
	Established types.EnableDisable `edge:"established"`
	Invalid     types.EnableDisable `edge:"invalid"`
	New         types.EnableDisable `edge:"new"`
	Related     types.EnableDisable `edge:"related"`
}

// Interfaces wraps all interfaces in the config
type Interfaces struct {
	Interfaces []Interface `edge:"ethernet {{ .Name }}"`
	Loopback   Loopback    `edge:"loopback lo"`
}

// Interface is a single interface on the router
type Interface struct {
	Name        string
	State       types.DisableProp           `edge:".,omitempty"`
	Address     []netip.Prefix              `edge:"address,omitempty"`
	Description string                      `edge:"description,omitempty"`
	AddressDHCP string                      `edge:"address,omitempty"`
	Duplex      AutoString                  `edge:"duplex"`
	IPv6        InterfaceIPv6Settings       `edge:"ipv6,omitempty"`
	MTU         uint16                      `edge:"mtu,omitempty"`
	Firewall    InterfaceFirewallAssignment `edge:"firewall,omitempty"`
	Speed       AutoUint32                  `edge:"speed"`
	VLANs       []VLAN                      `edge:"vif {{ .ID }}"`
}

// InterfaceIPv6Settings controls ipv6 settings/networks for the interface
type InterfaceIPv6Settings struct {
	DupAddrDetectTransmits uint16           `edge:"dup-addr-detect-transmits"`
	RouterAdvert           IPv6RouterAdvert `edge:"router-advert"`
}

// IPv6RouterAdvert Advertisement settings for ipv6
type IPv6RouterAdvert struct {
	CurHopLimit     uint16                    `edge:"cur-hop-limit"`
	LinkMTU         uint16                    `edge:"link-mtu"`
	ManagedFlag     bool                      `edge:"managed-flag"`
	MaxInterval     uint16                    `edge:"max-interval"`
	NameServer      netip.Addr                `edge:"name-server"`
	OtherConfigFlag bool                      `edge:"other-config-flag"`
	Prefixes        []IPv6PrefixAdvertisement `edge:"prefix {{ .Prefix }}"`
	ReachableTime   uint16                    `edge:"reachable-time"`
	RetransTimer    uint16                    `edge:"retrans-timer"`
	SendAdvert      bool                      `edge:"send-advert"`
}

// IPv6PrefixAdvertisement A single ipv6 prefix to announce
type IPv6PrefixAdvertisement struct {
	Prefix         netip.Prefix
	AutonomousFlag bool   `edge:"autonomous-flag"`
	OnLinkFlag     bool   `edge:"on-link-flag"`
	ValidLifetime  uint32 `edge:"valid-lifetime"`
}

// InterfaceFirewallAssignment Maps named firewall zones to an interface
type InterfaceFirewallAssignment struct {
	In    InterfaceFirewallZone `edge:"in,omitempty"`
	Out   InterfaceFirewallZone `edge:"out,omitempty"`
	Local InterfaceFirewallZone `edge:"local,omitempty"`
}

// InterfaceFirewallZone The name of the zone for the firewall
type InterfaceFirewallZone struct {
	Name   string `edge:"name,omitempty"`
	V6Name string `edge:"ipv6-name,omitempty"`
}

// Loopback Special interface struct for loopback
type Loopback struct {
}

// VLAN is how a vlan is defined in the router interface
type VLAN struct {
	ID          uint16
	Address     netip.Prefix `edge:"address"`
	Description string       `edge:"description"`
	MTU         uint16       `edge:"mtu"`
}

// RouterProtocols is the configuration for protocols on the router (BGP, etc)
type RouterProtocols struct {
	BGP    []BGPConfig    `edge:"bgp {{ .ASN }}"`
	Static StaticProtocol `edge:"static,omitempty"`
}

// BGPConfig is the configuration for a single one of our ASNs
type BGPConfig struct {
	ASN        uint32
	Neighbors  []BGPNeighbor `edge:"neighbor {{ .IP }}"`
	Networks   []BGPNetwork  `edge:"network {{ .Prefix }}"`
	Parameters BGPParameters `edge:"parameters"`
}

// BGPNeighbor Is a connection to a BGP peer for a single ASN
type BGPNeighbor struct {
	IP                  netip.Addr
	ASN                 uint32                 `edge:"remote-as"`
	DefaultOriginate    BGPDefaultOriginate    `edge:"default-originate,omitempty"`
	SoftReconfiguration BGPSoftReconfiguration `edge:"soft-reconfiguration,omitempty"`
}

// BGPNetwork is a single network announced to a peer
type BGPNetwork struct {
	Prefix netip.Prefix
}

// BGPDefaultOriginate struct to help get the proper formatting for this option
type BGPDefaultOriginate struct {
	Originate bool
}

// BGPSoftReconfiguration soft-reconfiguration
type BGPSoftReconfiguration struct {
	Inbound types.KeyWhenEnabled `edge:"inbound,omitempty"`
}

// BGPParameters is the parameters of the BGP connection
type BGPParameters struct {
	RouterID string `edge:"router-id"`
}

// StaticProtocol Wraps all static routes
type StaticProtocol struct {
	Routes []StaticRoute `edge:"route{{ .RouteSuffix }} {{ .Route }}"`
}

// StaticRoute is a static route in edgeconfig format
type StaticRoute struct {
	RouteSuffix string // This is empty for ipv4 and "6" for ipv6, so that we end up with route6 for ipv6
	Route       netip.Prefix
	NextHop     NextHop `edge:"next-hop {{ .NextHop }}"`
}

// NextHop is the next hop for the route
type NextHop struct {
	NextHop     netip.Addr
	Description string `edge:"description"`
	Distance    uint8  `edge:"distance"`
}

// RouterServices Available services on the router
type RouterServices struct {
	DHCPServer DHCPServer  `edge:"dhcp-server"`
	GUI        GUIService  `edge:"gui,omitempty"`
	NAT        NatService  `edge:"nat,omitempty"`
	SSH        SSHService  `edge:"ssh,omitempty"`
	UNMS       UNMSService `edge:"unms"`
}

// DHCPServer information about enabled DHCP servers
type DHCPServer struct {
	Disabled       bool                `edge:"disabled"`
	HostfileUpdate types.EnableDisable `edge:"hostfile-update"`
	Networks       []DHCPNetwork       `edge:"shared-network-name {{ .Name }}"`
	StaticARP      types.EnableDisable `edge:"static-arp"`
	UseDNSMASQ     types.EnableDisable `edge:"use-dnsmasq"`
}

// DHCPNetwork is a single network managed by the DHCP server
type DHCPNetwork struct {
	Name          string
	Authoritative types.EnableDisable `edge:"authoritative"`
	Subnets       []DHCPSubnet        `edge:"subnet {{ .Subnet }}"`
}

// DHCPSubnet is a subnet within a DHCP Network
type DHCPSubnet struct {
	Subnet         netip.Prefix
	Router         netip.Addr          `edge:"default-router"`
	DNS            []netip.Addr        `edge:"dns-server"`
	Lease          uint64              `edge:"lease"`
	StartStop      DHCPStartStop       `edge:"start {{ .Start }}"`
	StaticMappings []DHCPStaticMapping `edge:"static-mapping {{ .Name }}"`
}

// DHCPStaticMapping is a single DHCP reservation in this DHCP server
type DHCPStaticMapping struct {
	Name       string
	IPAddress  netip.Addr `edge:"ip-address"`
	MACAddress string     `edge:"mac-address"`
}

// DHCPStartStop is the start/stop range for DHCP Servers
type DHCPStartStop struct {
	Start netip.Addr
	Stop  netip.Addr `edge:"stop"`
}

// GUIService Settings for GUI
type GUIService struct {
	HTTPPort     uint16              `edge:"http-port"`
	HTTPSPort    uint16              `edge:"https-port"`
	OlderCiphers types.EnableDisable `edge:"older-ciphers"`
}

// NatService NAT settings
type NatService struct {
	// Must be in distinct numbering blocks, so splitting to make that easier
	Dest []NatRule
	Src  []NatRule
}

// MarshalEdge not used for NatService
func (ns NatService) MarshalEdge() ([]byte, error) {
	return nil, fmt.Errorf("marshaledge not implemented for NatService")
}

// MarshalEdgeWithDepth custom marshaller for NatService to ensure we get the numbers correct
func (ns NatService) MarshalEdgeWithDepth(depth int) ([]byte, error) {
	var buffer bytes.Buffer
	destCount := 1
	srcCount := 5000
	for _, rule := range ns.Dest {
		buffer.WriteString(fmt.Sprintf("%srule %d {\n", strings.Repeat(" ", depth), destCount))
		err := marshalValue(&buffer, reflect.ValueOf(rule), depth+4)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("%s}\n", strings.Repeat(" ", depth)))
		destCount++
	}
	for _, rule := range ns.Src {
		buffer.WriteString(fmt.Sprintf("%srule %d {\n", strings.Repeat(" ", depth), srcCount))
		err := marshalValue(&buffer, reflect.ValueOf(rule), depth+4)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("%s}\n", strings.Repeat(" ", depth)))
		srcCount++
	}
	return buffer.Bytes(), nil
}

// NatRule a single NAT rule
type NatRule struct {
	Name              string              `edge:"description"`
	Destination       types.AddressPort   `edge:"destination,omitempty"`
	InboundInterface  string              `edge:"inbound-interface,omitempty"`
	InsideAddress     types.AddressPort   `edge:"inside-address,omitempty"`
	Log               types.EnableDisable `edge:"log"`
	OutboundInterface string              `edge:"outbound-interface,omitempty"`
	OutsideAddress    types.AddressPort   `edge:"outside-address,omitempty"`
	Protocol          types.Protocol      `edge:"protocol,omitempty"`
	Source            types.AddressPort   `edge:"source,omitempty"`
	Type              types.NATType       `edge:"type"`
}

// SSHService settings for ssh
type SSHService struct {
	Port            uint16 `edge:"port"`
	ProtocolVersion string `edge:"protocol-version"`
}

// UNMSService settings for UNMS (only supported here to keep an empty block in config)
type UNMSService struct{}

// RouterSystem is the system config for the router
type RouterSystem struct {
	AnalyticsHandler AnalyticsHandler `edge:"analytics-handler"`
	CrashHandler     CrashHandler     `edge:"crash-handler"`
	HostName         string           `edge:"host-name"`
	Login            RouterLogin      `edge:"login,omitempty"`
	NTP              NTPServers       `edge:"ntp,omitempty"`
	Syslog           Syslog           `edge:"syslog,omitempty"`
	TimeZone         string           `edge:"time-zone,omitempty"`
}

// AnalyticsHandler settings for analytics
type AnalyticsHandler struct {
	SendAnalyticsreport bool `edge:"send-analytics-report"`
}

// CrashHandler settings for crash handling
type CrashHandler struct {
	SendCrashReport bool `edge:"send-crash-report"`
}

// RouterLogin handles user accounts for the router
type RouterLogin struct {
	Users []User `edge:"user {{ .Username }}"`
}

// User is a single router user
type User struct {
	Username       string
	Authentication Authentication  `edge:"authentication"`
	Level          types.UserLevel `edge:"level"`
}

// Authentication is the auth details for a user
type Authentication struct {
	EncryptedPassword string `edge:"encrypted-password,omitempty"`
	PlaintextPassword string `edge:"plaintext-password,omitempty"`
}

// NTPServers wraps the ntp servers in the config
type NTPServers struct {
	Servers []NTPServer `edge:"server {{ .Hostname }}"`
}

// NTPServer a single NTP server
type NTPServer struct {
	Hostname string
}

// Syslog section of system config
type Syslog struct {
	Global SyslogGlobal `edge:"global,omitempty"`
}

// SyslogGlobal is the global syslog config section
type SyslogGlobal struct {
	Facilities []SyslogFacility `edge:"facility {{ .Name }},omitempty"`
}

// SyslogFacility a single syslog config
type SyslogFacility struct {
	Name  string
	Level string `edge:"level"`
}
