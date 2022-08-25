/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.ibm.com/Nathan-Good/atkcli/internal/outputters"
	"github.ibm.com/Nathan-Good/atkdep"
)

func lookup(from []atkdep.IndexInfo, with string) (*atkdep.IndexInfo, bool) {
	for _, e := range from {
		if e.Id == with {
			return &e, true
		}
	}
	return nil, false
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info <MODULE>",
	Short: "Gets more information about the specified module",
	Args:  cobra.ExactArgs(1),
	Long: `
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("getting information about module: ", args[0])

		idxReader := atkdep.NewAtkGitHubIndexReader("test")
		entries, err := idxReader.List()
		if err != nil {
			return err
		}

		entry, exists := lookup(entries, args[0])

		if !exists {
			return fmt.Errorf("no module found with name %s", args[0])
		}

		modInfo, err := idxReader.Info(*entry)

		if err != nil {
			return err
		}

		outputters.WriteModuleInfo(modInfo, modulesCmd.OutOrStdout())

		return nil
	},
}

func init() {
	modulesCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
