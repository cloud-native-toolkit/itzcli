/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// modulesCmd represents the modules command
var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Command for working with modules",
	Long: `This command allows you to add, update, and list modules that are
used by the Activation TookKit (ATK) on your system.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("modules called")
	},
}

func init() {
	rootCmd.AddCommand(modulesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modulesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modulesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
