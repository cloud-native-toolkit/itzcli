/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	logger "github.com/sirupsen/logrus"
	"github.ibm.com/Nathan-Good/atkdep"
	"github.ibm.com/Nathan-Good/atkmod"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Configures, generates, and runs the GitOps projects for OPC",
	Long: `This is a "meta-command" that calls other commands--configure,
	generate, and run.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("starting installation...")
		modLoader := atkdep.NewGitHubModuleDownloader()
		path := filepath.Join(atkdep.AtkCliHomeDir(), "staged.yaml")
		selected, err := readFromFile(path)
		if err != nil {
			return err
		}

		selectedModules, err := modLoader.FetchSelected(selected)

		module := &selectedModules[0]

		ctx := context.Background()
		ctx = context.WithValue(ctx, atkmod.BaseDirectory, "/tmp")
		outbuff := new(bytes.Buffer)
		errbuff := new(bytes.Buffer)

		runCtx := &atkmod.RunContext{
			Out: outbuff,
			Err: errbuff,
			Log: *logger.StandardLogger(),
		}

		deployment := atkmod.NewDeployableModule(ctx, runCtx, module)
		// For the test purposes, let us just start out with this ready to pre-deploy
		deployment.Notify(atkmod.Validated)
		// Gets the correct command for the current state
		//for cmd, exists := deployment.Next(); exists == true; {
		//	err = cmd(runCtx, deployment)
		//	if err != nil {
		//		return err
		//	}
		//}

		return nil
	},
}

// TODO: This is really gross. This should come from the library, not defined here.
func readFromFile(filename string) ([]atkdep.IndexInfo, error) {
	var entries []atkdep.IndexInfo = []atkdep.IndexInfo{}
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &entries)

	if err != nil {
		return nil, err
	}

	return entries, nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}
