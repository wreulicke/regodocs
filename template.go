package regodocs

import (
	"text/template"

	"github.com/open-policy-agent/opa/ast"
)

var docTemplate = `
# {{.Package}}

## Rules

{{range .Rules -}}
### {{.Name}}

{{.Description}}

{{end}}
`

type Rule struct {
	Name        string
	Description string
}

func newTemplateData(pkg string, rules []*ast.Rule) *templateData {
	td := &templateData{
		Package: pkg,
	}
	for _, rule := range rules {
		r := newRule(rule)
		if r != nil {
			td.Rules = append(td.Rules, r)
		}
	}
	return td
}

func newRule(rule *ast.Rule) *Rule {
	if len(rule.Annotations) == 0 {
		return nil
	}

	var r Rule
	r.Name = rule.Head.Name.String()
	if rule.Annotations[0].Title != "" {
		r.Name = rule.Annotations[0].Title
	}
	if rule.Annotations[0].Description != "" {
		r.Description = rule.Annotations[0].Description
	}
	return &r
}

type templateData struct {
	Package string
	Rules   []*Rule
}

func newTemplate() *template.Template {
	return template.Must(template.New("doc").Parse(docTemplate))
}
