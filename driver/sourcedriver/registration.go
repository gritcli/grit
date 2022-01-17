package sourcedriver

// Registration encapsulates the information required to register a source
// driver implementation with Grit.
type Registration struct {
	// Name is the (preferred) name of the driver, as referenced within
	// configuration files.
	Name string

	// Description is a short human-readable description of the driver.
	Description string

	// NewConfigSchema returns a zero-value ConfigSchema for this driver.
	NewConfigSchema func() ConfigSchema

	// ImplicitSources is a set of sources that should be added to the
	// configuration automatically.
	ImplicitSources map[string]func() ConfigSchema
}
