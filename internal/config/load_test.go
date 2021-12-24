package config_test

import (
	. "github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load()", func() {
	DescribeTable(
		"it returns the expected configuration",
		func(dir string, expect Config) {
			cfg, err := Load(dir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfg).To(Equal(expect))
		},
		Entry(
			"default configuration",
			"testdata/valid/default",
			DefaultConfig,
		),
		Entry(
			"empty configuration file (should be the same as the default)",
			"testdata/valid/empty-file",
			DefaultConfig,
		),
		Entry(
			"empty configuration directory (should be the same as the default)",
			"testdata/valid/empty-dir",
			DefaultConfig,
		),
		Entry(
			`not existent directory (should be the same as the default)`,
			`testdata/valid/non-existent`,
			DefaultConfig,
		),
		Entry(
			"ignores non-HCL files, directories and files beginning with underscores",
			"testdata/valid/ignore",
			DefaultConfig,
		),
		// Entry(
		// 	"implicit github source disabled",
		// 	"testdata/valid/github-disabled.conf",
		// 	Config{Dir: "~/grit"},
		// ),
		// Entry(
		// 	"implicit github source overridden",
		// 	"testdata/valid/github-overridden.conf",
		// 	Config{
		// 		Dir: "~/grit",
		// 		Sources: map[string]Source{
		// 			"github": GitHubSource{
		// 				SourceName: "github",
		// 				API: &url.URL{
		// 					Scheme: "https",
		// 					Host:   "github.example.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// ),
		// Entry(
		// 	"implicit github source augmented with token",
		// 	"testdata/valid/github-augmented.conf",
		// 	Config{
		// 		Dir: "~/grit",
		// 		Sources: map[string]Source{
		// 			"github": GitHubSource{
		// 				SourceName: "github",
		// 				API: &url.URL{
		// 					Scheme: "https",
		// 					Host:   "api.github.com",
		// 				},
		// 				Token: "<token>",
		// 			},
		// 		},
		// 	},
		// ),
		// Entry(
		// 	"custom git source defined",
		// 	"testdata/valid/git-custom.conf",
		// 	Config{
		// 		Dir: "~/grit",
		// 		Sources: map[string]Source{
		// 			"github": DefaultConfig.Sources["github"],
		// 			"my-company": GitSource{
		// 				SourceName: "my-company",
		// 				Endpoint: &transport.Endpoint{
		// 					Protocol: "ssh",
		// 					User:     "git",
		// 					Host:     "git.example.com",
		// 					Port:     22,
		// 					Path:     "{repo}.git",
		// 				},
		// 			},
		// 		},
		// 	},
		// ),
		// Entry(
		// 	"custom github source defined",
		// 	"testdata/valid/github-custom.conf",
		// 	Config{
		// 		Dir: "~/grit",
		// 		Sources: map[string]Source{
		// 			"github": DefaultConfig.Sources["github"],
		// 			"my-company": GitHubSource{
		// 				SourceName: "my-company",
		// 				API: &url.URL{
		// 					Scheme: "https",
		// 					Host:   "github.example.com",
		// 				},
		// 				Token: "<token>",
		// 			},
		// 		},
		// 	},
		// ),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		func(dir string, expect string) {
			_, err := Load(dir)
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expect), err.Error())
		},
		Entry(
			`syntax error`,
			`testdata/invalid/syntax-error`,
			`testdata/invalid/syntax-error/grit.hcl:1,1-2: Argument or block definition required; An argument or block definition is required here.`,
		),
		Entry(
			`unrecognized file`,
			`testdata/invalid/unrecognized-file`,
			`testdata/invalid/unrecognized-file/unrecognized.hcl: unrecognized configuration file`,
		),
	)
})
