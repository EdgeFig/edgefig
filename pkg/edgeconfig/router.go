package edgeconfig

import (
	"net/netip"
)

// Router is the top level config that applies to EdgeRouters
type Router struct {
	Firewall   Firewall             `edge:"firewall"`
	Interfaces map[string]Interface `edge:"interface"`
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
