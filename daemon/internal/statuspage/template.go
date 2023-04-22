package statuspage

import (
	_ "embed"
	html "html/template"
	text "text/template"
)

var (
	//go:embed template.html
	htmlContent  string
	htmlTemplate *html.Template

	//go:embed template.md
	markdownContent  string
	markdownTemplate *text.Template
)

func init() {
	htmlTemplate = html.Must(html.New("statuspage").Parse(htmlContent))
	markdownTemplate = text.Must(text.New("statuspage").Parse(markdownContent))
}

// TemplateValues contains the values that are passed to the status page
// template.
type TemplateValues struct {
	Title      string
	Heading    string
	SubHeading string
	Paragraphs []any
}
