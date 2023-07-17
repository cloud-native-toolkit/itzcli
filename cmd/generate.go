package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var generateCmd = &cobra.Command{
	Use:    "generate",
	Short:  "Tool for generating the documentation and other helpers",
	Long:   "Tool for generating the documentation and other helpers",
	PreRun: SetLoggingLevel,
	Hidden: true,
}

var generateManCmd = &cobra.Command{
	Use:    "man",
	Short:  "Generates the man documentation",
	Long:   "Generates the man documentation",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Root().DisableAutoGenTag = true
		mh := &doc.GenManHeader{
			Title:   "ITZ",
			Section: "1",
		}
		err := doc.GenManTree(cmd.Root(), mh, "./contrib/manpages")
		if err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

var generateHTMLCmd = &cobra.Command{
	Use:    "html",
	Short:  "Generates the html documentation",
	Long:   "Generates the html documentation",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Root().DisableAutoGenTag = true
		err := doc.GenMarkdownTree(cmd.Root(), "./docs")
		if err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateManCmd)
	generateCmd.AddCommand(generateHTMLCmd)

	rootCmd.AddCommand(generateCmd)
}
