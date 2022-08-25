/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	logger "github.com/sirupsen/logrus"
	"github.ibm.com/Nathan-Good/atkcli/internal/outputters"
	"github.ibm.com/Nathan-Good/atkdep"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <MODULE>",
	Args:  cobra.ExactArgs(1),
	Short: "Adds the module to the installation",
	Long: `
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("adding module: ", args[0])

		idxReader := atkdep.NewAtkGitHubIndexReader("test")
		entries, err := idxReader.List()
		if err != nil {
			return err
		}

		entry, exists := lookup(entries, args[0])

		outputters.WriteToFile(entry, filepath.Join(atkdep.AtkCliHomeDir(), "staged.yaml"))

		if !exists {
			return fmt.Errorf("no module found with name %s", args[0])
		}

		return nil
	},
}

func init() {
	modulesCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
