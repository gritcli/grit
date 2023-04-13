package statuspage

import (
	"fmt"
	html "html/template"
	"io"
	"net/http"
	"regexp"

	accept "github.com/timewasted/go-accept-headers"
)

// Render writes a status page to w.
func Render(
	w http.ResponseWriter,
	r *http.Request,
	code int,
	values TemplateValues,
) {
	contentType, _ := accept.Negotiate(
		r.Header.Get("Accept"),
		"text/plain",
		"text/html",
		"application/xhtml+xml",
	)

	var template interface {
		Execute(w io.Writer, data any) error
	}

	switch contentType {
	case "text/html", "application/xhtml+xml":
		template = htmlTemplate

		backticks := regexp.MustCompile("`[^`]*`")
		for i, p := range values.Paragraphs {
			if p, ok := p.(string); ok {
				values.Paragraphs[i] = html.HTML(
					backticks.ReplaceAllStringFunc(
						p,
						func(s string) string {
							return "<code>" + s[1:len(s)-1] + "</code>"
						},
					),
				)
			}
		}

	default:
		contentType = "text/plain"
		template = markdownTemplate
	}

	w.Header().Add("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(code)
	template.Execute(w, values)
}

// RenderError writes a status page to w.
func RenderError(
	w http.ResponseWriter,
	r *http.Request,
	code int,
) {
	text := http.StatusText(code)
	explanation := StatusExplanation(code)

	Render(
		w,
		r,
		code,
		TemplateValues{
			Title:      fmt.Sprintf("%d %s", code, text),
			Heading:    fmt.Sprintf("%d", code),
			SubHeading: text,
			Paragraphs: []any{explanation, "another paragraph."},
		},
	)
}
