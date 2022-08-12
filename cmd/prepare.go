/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// prepareCmd represents the prepare command
var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Prepares the specified project for building",
	Long: `Prepares the project by applying the given variable
values against the template.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prepare called")
	},
}

func init() {
	projectCmd.AddCommand(prepareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prepareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prepareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
