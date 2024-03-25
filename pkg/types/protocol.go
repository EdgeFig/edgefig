package types

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
