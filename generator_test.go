package regodocs

import (
	"io"
	"os"
	"testing"

	"github.com/gobwas/glob"
	"github.com/wreulicke/snap"
)

func newDefaultRegexps() []glob.Glob {
	return []glob.Glob{
		glob.MustCompile("{deny*, violation*, warn*}"),
	}
}

func TestGenerator(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	g := NewGenerator(&GeneratorConfig{
		OutputPath:        dir,
		Patterns:          newDefaultRegexps(),
		IgnoreFilePattern: []glob.Glob{glob.MustCompile("*_test.rego")},
	})

	err := g.Generate([]string{"testdata"})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	s := snap.New()
	f, err := os.Open(dir + "/bar.md")
	if os.IsNotExist(err) {
		t.Errorf("expected file to exist, got %v", err)
	}
	io.Copy(s, f)
	s.Assert(t)
}
