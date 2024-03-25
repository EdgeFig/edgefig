package edgeconfig

// EnableDisable is a simple boolean wrapper that marshals to "enable" or "disable"
type EnableDisable bool

const (
	// Enable is a quick const to set EnableDisable to Enable
	Enable = EnableDisable(true)
	// Disable is a quick const to set EnableDisable to Disable
	Disable = EnableDisable(false)
)
