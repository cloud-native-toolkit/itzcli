/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listVariablesCmd represents the listVariables command
var listVariablesCmd = &cobra.Command{
	Use:   "list-variables",
	Short: "Lists the variables required by the project",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list-variables called")
	},
}

func init() {
	projectCmd.AddCommand(listVariablesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listVariablesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listVariablesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
