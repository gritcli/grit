# {{if .HasParent}}{{.Parent.CommandPath}} {{end}}{{.NameAndAliases}}

{{.Short}}

## usage

{{if .HasAvailableSubCommands}}
`{{.CommandPath}} <sub-command> [...]`
{{else}}
`{{.Use}}`
{{end}}

{{if .Long}}
## description

{{.Long}}
{{end}}

{{if .HasExample}}
## example

```
{{.Example}}
```

{{end}}


{{if .HasAvailableSubCommands}}
## sub-commands

{{range .Commands}}
{{if (or .IsAvailableCommand (eq .Name "help"))}}
- `{{.Name}}` &mdash; {{.Short}}
{{end}}
{{end}}
{{end}}

{{if .LocalFlags}}
## options

{{range .LocalFlags}}
{{if ne .Long  "help" }}
- `{{if .Short}}-{{.Short}}, {{end}}--{{.Long}}{{if .Arg}} <{{.Arg}}>{{end}}`{{if .Default}} [default = `{{.Default}}`]{{end}}

{{.Usage}}
{{end}}
{{end}}
{{end}}

{{if .GlobalFlags}}
## global options

{{range .GlobalFlags}}
{{if ne .Long  "help" }}
- `{{if .Short}}-{{.Short}}, {{end}}--{{.Long}}{{if .Arg}} <{{.Arg}}>{{end}}` {{if .Default}} [default = `{{.Default}}`]{{end}}

{{.Usage}}
{{end}}
{{end}}
{{end}}
