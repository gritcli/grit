package source_test

import (
	. "github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/internal/stubs"
	"github.com/gritcli/grit/logs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Suggester", func() {
	var (
		repoA1, repoA2, repoB1, repoB2 sourcedriver.RemoteRepo
		srcA, srcB                     *stubs.Source
		suggester                      *Suggester
	)

	BeforeEach(func() {
		repoA1 = sourcedriver.RemoteRepo{
			Name: "<repo-a1>",
		}
		repoA2 = sourcedriver.RemoteRepo{
			Name: "<repo-a2>",
		}
		repoB1 = sourcedriver.RemoteRepo{
			Name: "<repo-b1>",
		}
		repoB2 = sourcedriver.RemoteRepo{
			Name: "<repo-b2>",
		}

		srcA = &stubs.Source{
			SuggestFunc: func(w string, log logs.Log) map[string][]sourcedriver.RemoteRepo {
				Expect(w).To(Equal("<word>"))
				return map[string][]sourcedriver.RemoteRepo{
					"<word>": {repoA1, repoA2},
				}
			},
		}
		srcB = &stubs.Source{
			SuggestFunc: func(w string, log logs.Log) map[string][]sourcedriver.RemoteRepo {
				Expect(w).To(Equal("<word>"))
				return map[string][]sourcedriver.RemoteRepo{
					"<word>": {repoB1, repoB2},
				}
			},
		}

		suggester = &Suggester{
			Sources: List{
				{
					Name:   "<source-a>",
					Driver: srcA,
				},
				{
					Name:   "<source-b>",
					Driver: srcB,
				},
			},
		}
	})

	Describe("func Suggest()", func() {
		It("aggregates the suggestions from all sources", func() {
			matches := suggester.Suggest("<word>", true, true)
			Expect(matches).To(Equal(
				map[string][]sourcedriver.RemoteRepo{
					"<word>": {
						repoA1,
						repoA2,
						repoB1,
						repoB2,
					},
				},
			))
		})
	})
})
