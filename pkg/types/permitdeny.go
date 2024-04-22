package types

// PermitDeny is a wrapper for known string values for permit and deny
type PermitDeny string

const (
	// Permit is the "permit" setting
	Permit PermitDeny = "permit"

	// Deny is the "deny" setting
	Deny PermitDeny = "deny"
)
