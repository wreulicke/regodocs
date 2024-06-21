package regodocs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
)

type Generator struct {
	GeneratorConfig *GeneratorConfig
}

type GeneratorConfig struct {
	OutputPath        string
	Patterns          []glob.Glob
	IgnoreFilePattern []glob.Glob
}

type packageSet struct {
	Package *Package
	Files   []*loader.RegoFile
}

func NewGenerator(c *GeneratorConfig) *Generator {
	return &Generator{GeneratorConfig: c}
}

func (g *Generator) Generate(paths []string) error {
	f, err := loader.NewFileLoader().WithProcessAnnotation(true).
		Filtered(paths, func(abspath string, info fs.FileInfo, _ int) bool {
			return !info.IsDir() && !strings.HasSuffix(abspath, ".rego") && !matchGlobs(g.GeneratorConfig.IgnoreFilePattern, abspath)
		})
	if err != nil {
		return fmt.Errorf("failed to load policy: %w", err)
	}

	packageMap := make(map[string]*packageSet)
	for _, m := range f.Modules {
		packageName := m.Parsed.Package.Path.String()
		packageName = strings.TrimPrefix(packageName, "data.") // i don't know why this is necessary
		if packageMap[packageName] == nil {
			packageMap[packageName] = &packageSet{
				Package: &Package{
					Name: packageName,
				},
				Files: []*loader.RegoFile{m},
			}
		} else {
			packageMap[packageName].Files = append(packageMap[packageName].Files, m)
		}
		if len(m.Parsed.Annotations) > 0 && m.Parsed.Annotations[0].Description != "" {
			if packageMap[packageName].Package.Description != "" {
				return fmt.Errorf("package %s has multiple descriptions", packageName)
			}
			packageMap[packageName].Package.Description = m.Parsed.Annotations[0].Description
		}
	}

	t := newTemplate()

	write := func(data templateData, f *os.File) error {
		defer f.Close()
		return t.Execute(f, data)
	}
	for pkg, packageSet := range packageMap {
		var rules []*ast.Rule
		for _, file := range packageSet.Files {
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
		err = write(*newTemplateData(packageSet.Package, rules), f)
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
	return matchGlobs(g.GeneratorConfig.Patterns, rule.Head.Name.String())
}

func matchGlobs(globs []glob.Glob, path string) bool {
	for _, p := range globs {
		if p.Match(path) {
			return true
		}
	}
	return false
}
