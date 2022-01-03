package config

import (
	"reflect"

	"github.com/hashicorp/hcl/v2"
)

// configFile is HCL schema for a configuration file.
type configFile struct {
	DaemonBlock  *daemonBlock  `hcl:"daemon,block"`
	GitBlock     *gitBlock     `hcl:"git,block"`
	SourceBlocks []sourceBlock `hcl:"source,block"`
}

// daemonBlock is the HCL schema for a "daemon" block.
type daemonBlock struct {
	Socket string `hcl:"socket,optional"`
}

// gitBlock is the HCL schema for a "git" block
type gitBlock struct {
	PrivateKey *string `hcl:"private_key"`
	PreferHTTP *bool   `hcl:"prefer_http"`
}

// sourceBlock is the HCL schema for a "source" block.
type sourceBlock struct {
	Name    string   `hcl:",label"`
	Impl    string   `hcl:",label"`
	Enabled *bool    `hcl:"enabled"`
	Body    hcl.Body `hcl:",remain"` // see sourceBlockBody
}

// sourceBlockBody is an interface for implementation-specific HCL schema within
// the body of a source block.
type sourceBlockBody interface {
	resolve(filename string, cfg Config) (SourceConfig, error)
}

// sourceSchemaByImpl is a map of a source implementation name to the type of
// its sourceBlockBody implementation.
var sourceSchemaByImpl = map[string]reflect.Type{}

// registerSourceSchema registers a source implementation, allowing its
// configuration to be parsed.
//
// impl is the name of the implementation, as specified in "source" blocks
// within the configuration file.
func registerSourceSchema(
	impl string,
	schema sourceBlockBody,
	defaultSources ...Source,
) {
	if _, ok := sourceSchemaByImpl[impl]; ok {
		panic("source name already registered")
	}

	sourceSchemaByImpl[impl] = reflect.TypeOf(schema)

	for _, s := range defaultSources {
		DefaultConfig.Sources[s.Name] = s
	}
}
