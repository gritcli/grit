package sourcedriver

// Registration encapsulates the information required to register a source
// driver implementation with Grit.
type Registration struct {
	// Name is the (preferred) name of the driver, as referenced within
	// configuration files.
	Name string

	// Description is a short human-readable description of the driver.
	Description string

	// ConfigLoader loads configuration for this driver.
	ConfigLoader ConfigLoader
}
