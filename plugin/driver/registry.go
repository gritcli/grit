package driver

// Registration encapsulates information about a source driver.
// TODO: find a better name for this
type Registration struct {
	// Name is the (preferred) name of the driver, as referenced within
	// configuration files, etc.
	Name string

	// Description is a short human-readable description of the driver.
	Description string

	// NewConfigSchema returns a zero-value ConfigSchema for this driver.
	NewConfigSchema func() ConfigSchema

	// DefaultSources is a set of default sources that should be added when this
	// driver is registered.
	DefaultSources map[string]func() ConfigSchema
}
