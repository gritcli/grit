package httpserver

import (
	"net/http"

	"github.com/gritcli/grit/daemon/internal/httpserver/statuspage"
)

type IndexHandler struct{}

func (*IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statuspage.Render(
		w,
		r,
		http.StatusOK,
		statuspage.TemplateValues{
			Title:      "Grit",
			Heading:    "Grit",
			SubHeading: "Manage your local VCS clones",
			Paragraphs: []any{
				"This webserver handles interactive authentication requests.",
				"To authenticate, run the `grit source login` command from your terminal.",
			},
		},
	)
}
