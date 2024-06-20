package regodocs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
)

type Generator struct {
	GeneratorConfig *GeneratorConfig
}

type GeneratorConfig struct {
	OutputPath string
	Patterns   []*regexp.Regexp
}

func NewGenerator(c *GeneratorConfig) *Generator {
	return &Generator{GeneratorConfig: c}
}

func (g *Generator) Generate(paths []string) error {
	f, err := loader.NewFileLoader().WithProcessAnnotation(true).
		Filtered(paths, func(abspath string, info fs.FileInfo, depth int) bool {
			return !info.IsDir() && !strings.HasSuffix(abspath, ".rego")
		})
	if err != nil {
		return fmt.Errorf("failed to load policy: %w", err)
	}

	packageMap := make(map[string][]*loader.RegoFile)
	for _, m := range f.Modules {
		packageName := m.Parsed.Package.Path.String()
		packageName = strings.TrimPrefix(packageName, "data.") // i don't know why this is necessary
		if packageMap[packageName] == nil {
			packageMap[packageName] = []*loader.RegoFile{m}
		} else {
			packageMap[packageName] = append(packageMap[packageName], m)
		}
	}

	t := newTemplate()

	write := func(data templateData, f *os.File) error {
		defer f.Close()
		return t.Execute(f, data)
	}
	for pkg, files := range packageMap {
		var rules []*ast.Rule
		for _, file := range files {
			rules = append(rules, file.Parsed.Rules...)
		}
		rules = g.matchedRules(rules)
		if len(rules) == 0 {
			continue
		}

		p := filepath.Join(g.GeneratorConfig.OutputPath, pkg+".md")
		f, err := os.Create(p)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		err = write(*newTemplateData(pkg, rules), f)
		if err != nil {
			return fmt.Errorf("failed to write template: %w", err)
		}
	}
	return nil
}

func (g *Generator) matchedRules(rules []*ast.Rule) []*ast.Rule {
	var r []*ast.Rule
	for _, rule := range rules {
		if g.matchRule(rule) {
			r = append(r, rule)
		}
	}
	return r
}

func (g *Generator) matchRule(rule *ast.Rule) bool {
	for _, p := range g.GeneratorConfig.Patterns {
		ruleName := rule.Head.Name.String()
		if p.MatchString(ruleName) {
			return true
		}
	}
	return true
}
