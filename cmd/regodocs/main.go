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

			var regexps []*regexp.Regexp
			for _, p := range patterns {
				r, err := regexp.Compile(p)
				if err != nil {
					return fmt.Errorf("failed to compile pattern %s: %w", p, err)
				}
				regexps = append(regexps, r)
			}
			g := regodocs.NewGenerator(&regodocs.GeneratorConfig{
				OutputPath: outputPath,
				Patterns:   regexps,
			})
			return g.Generate(args)
		},
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path for generated documentation")
	cmd.Flags().StringSliceVarP(&patterns, "pattern", "p", []string{"deny.*", "violation.*", "warn.*"}, "pattern to filter files")

	return cmd
}

func main() {
	cmd := RootCmd()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
