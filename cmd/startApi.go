/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Starting that api...")
		for _, cmd := range cmd.Parent().Commands() {
			// What we want is `itz solution list --list-all` becomes
			// <url>/api/itz/solution/list&list-all=true
			r.Add(fmt.Sprintf("%s", cmd.Name()), CliHanderToRESTHandler(cmd.Run))
		}
	},
}

func CliHanderToRESTHandler(run func(cmd *cobra.Command, args []string)) func(c *gin.Context) {
	return func(c *gin.Context) {
		run(startCmd, args)
	}
}

func init() {
	apiCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
