package edgeconfig

// EnableDisable is a simple boolean wrapper that marshals to "enable" or "disable"
type EnableDisable bool

const (
	// Enable is a quick const to set EnableDisable to Enable
	Enable = EnableDisable(true)
	// Disable is a quick const to set EnableDisable to Disable
	Disable = EnableDisable(false)
)

// DisableProp is similar to EnableDisable, but ONLY shows up if Disabled, and shows nothing if Enabled
// Used on interfaces
type DisableProp bool

const (
	// Disabled is the state for "disabled"
	Disabled DisableProp = true
	// Enabled is the state for "enabled" and nothing will show up in the config
	Enabled DisableProp = false
)
