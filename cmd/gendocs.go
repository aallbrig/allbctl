package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var genDocsOutputDir string

var genDocsCmd = &cobra.Command{
	Use:    "gen-docs",
	Short:  "Generate markdown CLI reference documentation",
	Long:   `Generate markdown CLI reference documentation into the Hugo site's reference directory.`,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll(genDocsOutputDir, 0755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}

		// Write a _index.md with Hugo front matter for the reference section
		index := filepath.Join(genDocsOutputDir, "_index.md")
		indexContent := `---
weight: 99
bookFlatSection: true
title: "CLI Reference"
---

# CLI Reference

Auto-generated reference documentation for all allbctl commands and flags.

> This section is generated from the source code via ` + "`make gen-docs`" + `.
> Do not edit these files manually.
`
		if err := os.WriteFile(index, []byte(indexContent), 0644); err != nil {
			return fmt.Errorf("write _index.md: %w", err)
		}

		// Generate one markdown file per command with Hugo front matter prepended
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			name = strings.TrimSuffix(name, ".md")
			return fmt.Sprintf("---\ntitle: \"%s\"\ndate: %s\ndescription: \"auto-generated reference\"\n---\n\n",
				name, time.Now().Format("2006-01-02"))
		}
		linkHandler := func(name string) string { return name }

		if err := doc.GenMarkdownTreeCustom(rootCmd, genDocsOutputDir, prepender, linkHandler); err != nil {
			return fmt.Errorf("generate docs: %w", err)
		}

		fmt.Fprintf(os.Stderr, "docs written to %s\n", genDocsOutputDir)
		return nil
	},
}

func init() {
	genDocsCmd.Flags().StringVar(&genDocsOutputDir, "output", "hugo/site/content/docs/reference", "Directory to write generated docs")
	rootCmd.AddCommand(genDocsCmd)
}
