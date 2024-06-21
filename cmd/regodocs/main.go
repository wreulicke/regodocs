package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/wreulicke/regodocs"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regodocs",
		Short: "root command for regodocs, just show help",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		GenerateCmd(),
	)

	return cmd
}

func GenerateCmd() *cobra.Command {
	var outputPath string
	var patterns []string
	var ignoreFilePatterns []string
	cmd := &cobra.Command{
		Use:        "generate POLICY_PATH...",
		Short:      "Generate documentation from Rego policy files",
		ArgAliases: []string{"POLICY_PATH..."},
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputPath == "" {
				return cmd.Help()
			}
			if len(args) == 0 {
				return cmd.Help()
			}

			patternRegexps, err := compileRegexps(patterns)
			if err != nil {
				return err
			}
			ignoreFileRegexps, err := compileRegexps(ignoreFilePatterns)
			if err != nil {
				return err
			}
			g := regodocs.NewGenerator(&regodocs.GeneratorConfig{
				OutputPath:        outputPath,
				Patterns:          patternRegexps,
				IgnoreFilePattern: ignoreFileRegexps,
			})
			return g.Generate(args)
		},
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path for generated documentation")
	cmd.Flags().StringSliceVarP(&patterns, "pattern", "p", []string{"deny.*", "violation.*", "warn.*"}, "regexp to filter rules")
	cmd.Flags().StringSliceVarP(&ignoreFilePatterns, "ignore", "i", []string{".*_test.rego"}, "regexp to ignore files")

	return cmd
}

func compileRegexps(regexps []string) ([]*regexp.Regexp, error) {
	var compiled []*regexp.Regexp
	for _, r := range regexps {
		re, err := regexp.Compile(r)
		if err != nil {
			return nil, fmt.Errorf("failed to compile regexp %s: %w", r, err)
		}
		compiled = append(compiled, re)
	}
	return compiled, nil
}

func main() {
	cmd := RootCmd()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
