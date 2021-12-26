package config

// SourceDriver is an enumeration of the supported source drivers.
//
// Each driver provides an implementation for communicating with a specific type
// of repository source, such as GitHub, BitBucket or a vanilla Git server.
type SourceDriver string

// DriverConfig is an interface for configuration that is specific to a
// particular source driver.
type DriverConfig interface {
	// DescribeConfig returns a short, human-readable description of the
	// configuration.
	//
	// It may not include information about all its available configuration
	// directives. It typically would include the most important aspects of the
	// configuration that can be used to disambiguate two sources that use the
	// same driver, or unusual non-default settings.
	DescribeConfig() string
}

// driverConfigPrototypes is a map of driver to an empty configuration structure
// that can be used to parse source configuration.
var driverConfigPrototypes = map[SourceDriver]DriverConfig{}

// registerDriver registers a source driver so that its configuration can be
// loaded.
//
// p is an empty "prototype" of the configuration struct used to parse source
// configuration for sources that use this driver.
func registerDriver(
	d SourceDriver,
	p DriverConfig,
	defaultSources ...Source,
) {
	driverConfigPrototypes[d] = p

	for _, s := range defaultSources {
		DefaultConfig.Sources[s.Name] = s
	}
}
