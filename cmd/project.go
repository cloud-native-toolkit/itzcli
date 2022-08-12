/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var projectName string

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Lists metadata, builds, and deploys projects",
	Long: `The project command provides a CLI for maintaining
the .now/manifest.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("project called")
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	projectCmd.PersistentFlags().StringVarP(&projectName, "name", "n", "", "The name of the project")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
