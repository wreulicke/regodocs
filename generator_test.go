package regodocs

import (
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/wreulicke/snap"
)

func newDefaultRegexps() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile("deny.*"),
		regexp.MustCompile("violation.*"),
		regexp.MustCompile("warn.*"),
	}
}

func TestGenerator(t *testing.T) {
	dir := t.TempDir()
	g := NewGenerator(&GeneratorConfig{
		OutputPath: dir,
		Patterns:   newDefaultRegexps(),
	})

	err := g.Generate([]string{"testdata"})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	s := snap.New()
	f, err := os.Open(dir + "/test.md")
	if os.IsNotExist(err) {
		t.Errorf("expected file to exist, got %v", err)
	}
	io.Copy(s, f)
	s.Assert(t)
}
