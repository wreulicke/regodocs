package regodocs

import (
	"strings"
	"text/template"

	"github.com/open-policy-agent/opa/ast"
)

var docTemplate = `
# {{.Package.Name}}

{{.Package.Description}}
## Rules

{{range .Rules -}}
- [{{.Name}}](#{{.Name | anchor}})
{{- end}}

{{range .Rules -}}
### {{.Name}}

{{.Description}}
{{- end}}
`

type Package struct {
	Name        string
	Description string
}

type Rule struct {
	Name        string
	Description string
}

type templateData struct {
	Package *Package
	Rules   []*Rule
}

func newTemplateData(pkg *Package, rules []*ast.Rule) *templateData {
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

func funcMap() template.FuncMap {
	return template.FuncMap{
		"anchor": anchor,
	}
}

func anchor(heading string) string {
	return strings.ReplaceAll(strings.ToLower(heading), " ", "-")
}

func newTemplate() *template.Template {
	return template.Must(template.New("doc").Funcs(funcMap()).Parse(docTemplate))
}
