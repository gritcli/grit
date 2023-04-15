package httpserver

import (
	"net/http"

	"github.com/gritcli/grit/daemon/internal/httpserver/statuspage"
)

// IndexHandler is the HTTP handler for the index page.
type IndexHandler struct{}

func (*IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		statuspage.RenderError(w, r, http.StatusNotFound)
		return
	}

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
				"To authenticate, run the `grit source sign-in` command from your terminal.",
			},
		},
	)
}
