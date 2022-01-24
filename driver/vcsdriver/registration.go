package vcsdriver

// Registration encapsulates the information required to register a source
// driver implementation with Grit.
type Registration struct {
	// Name is the (preferred) name of the driver, as referenced within
	// configuration files.
	Name string

	// Description is a short human-readable description of the driver.
	Description string

	// ConfigNormalizer is the normalizer used to produce configuration values
	// for this driver.
	ConfigNormalizer ConfigNormalizer
}
