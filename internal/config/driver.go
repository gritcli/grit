package config

// DriverConfig is an interface for configuration that is specific to a
// particular source driver.
type DriverConfig interface {
	// String returns a short, human-readable description of the
	// configuration.
	//
	// It may not include information about all its available configuration
	// directives. It typically would include the most important aspects of the
	// configuration that can be used to disambiguate two sources that use the
	// same driver, or unusual non-default settings.
	String() string

	// acceptVisitor calls the appropriate driver-specific method on v.
	acceptVisitor(s Source, v SourceVisitor)
}

// driverConfigPrototypes is a map of driver to an empty configuration structure
// that can be used to parse source configuration.
var driverConfigPrototypes = map[string]DriverConfig{}

// registerDriver registers a source driver so that its configuration can be
// loaded.
//
// p is an empty "prototype" of the configuration struct used to parse source
// configuration for sources that use this driver.
func registerDriver(
	name string,
	proto DriverConfig,
	defaultSources ...Source,
) {
	driverConfigPrototypes[name] = proto

	for _, s := range defaultSources {
		DefaultConfig.Sources[s.Name] = s
	}
}
