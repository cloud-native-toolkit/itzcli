package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkmod"
	"github.ibm.com/skol/itzcli/internal/prompt"
	"github.ibm.com/skol/itzcli/pkg"
	"io"
	"os"
)

var fn string
var sol string
var cluster string
var rez string
var useCached bool

// deploySolutionCmd represents the deployProject command
var deploySolutionCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the specified solution.",
	Long: `Use this command to deploy the specified solution
locally in your own environment. You can specify the environment by using
either --cluster-name or --reservation as a target.

    --cluster-name requires the name of a cluster that has been deployed
using ocpnow. To see the clusters that are configured, use the "itz configure 
list" command to list the available clusters. If you have none, you may need to
import the ocpnow configuration using the "itz configure import" command. See
the help for those commands for more information.

    --reservation requires the id of a reservation in the IBM Technology Zone system. Use
the "itz reservation list" command to list the available reservations.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		SetLoggingLevel(cmd, args)
		if len(fn) == 0 && len(sol) == 0 {
			return fmt.Errorf("either \"--solution\" or \"--file\" must be specified")
		}
		if len(cluster) == 0 && len(rez) == 0 {
			return fmt.Errorf("either \"--cluster-name\" or \"--reservation\" must be specified")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Deploying solution \"%s\"...", sol)
		return DeploySolution(cmd, args)
	},
}

func init() {
	solutionCmd.AddCommand(deploySolutionCmd)
	deploySolutionCmd.Flags().StringVarP(&fn, "file", "f", "", "The full path to the solution file to be deployed.")
	deploySolutionCmd.Flags().StringVarP(&sol, "solution", "s", "", "The name of the solution to be deployed.")
	deploySolutionCmd.Flags().StringVarP(&cluster, "cluster-name", "c", "", "The name of the cluster created by ocpnow to target.")
	deploySolutionCmd.Flags().StringVarP(&rez, "reservation", "r", "", "The id of the reservation to target.")
	// TODO: Change this from true to false by default
	deploySolutionCmd.Flags().BoolVarP(&useCached, "use-cache", "u", false, "If true, uses a cached solution file instead of downloading from target.")
	deploySolutionCmd.MarkFlagsMutuallyExclusive("file", "solution")
	deploySolutionCmd.MarkFlagsMutuallyExclusive("reservation", "cluster-name")
}

// DeploySolution deploys the solution by handing it off to the bifrost
// API
func DeploySolution(cmd *cobra.Command, args []string) error {
	// TODO: Eventually, it would be really neat to have some way of making
	// this be configurable, too. Or maybe this is just moved from here to a
	// container...
	var vars = make([]pkg.JobParam, 0)
	var rootQuestion *prompt.Prompt

	prompterHandler := func(buf *bytes.Buffer) error {
		data, err := io.ReadAll(buf)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &vars)
		if err != nil {
			return err
		}

		ocpCfg := viper.GetStringSlice("ocpnow.configFiles")
		if len(ocpCfg) == 0 {
			return fmt.Errorf("no OPC configuration found")
		}

		project, err := pkg.LoadProject(ocpCfg[0])
		if err != nil {
			return err
		}
		resolver, err := pkg.NewBuildParamResolver(project, cluster, vars)
		rootQuestion, err = resolver.BuildPrompter(sol)

		nextPrompter := rootQuestion.Itr()

		for p := nextPrompter(); p != nil; p = nextPrompter() {
			logger.Tracef("Asking <%s>", p.String())
			err = prompt.Ask(p, os.Stdout, os.Stdin)
			if err != nil {
				return err
			}
		}

		return nil
	}

	paramsHandler := func(buf *bytes.Buffer) error {
		var tfvars = make([]pkg.JobParam, 0)
		for k, v := range rootQuestion.VarMap() {
			if len(v) > 0 {
				tfvars = append(tfvars, pkg.JobParam{Name: k, Value: v})
			}
		}
		json, err := json.Marshal(tfvars)
		if err != nil {
			return err
		}
		buf.Write(json)
		return nil
	}

	setEnvs := func(cli *atkmod.CliModuleRunner) error {
		cli.WithEnvvar("ITZ_SOLUTION_PLATFORM", "ibm")
		cli.WithEnvvar("ITZ_SOLUTION_STORAGE", "odf")
		cli.WithEnvvar("ITZ_SOLUTION_LAYERS", "280")
		return nil
	}

	pkg.DoContainerizedStep(cmd, "getcode", nil, nil, nil)
	pkg.DoContainerizedStep(cmd, "listparams", nil, prompterHandler, nil)
	pkg.DoContainerizedStep(cmd, "setparams", paramsHandler, nil, nil)
	pkg.DoContainerizedStep(cmd, "applyall", nil, nil, setEnvs)

	return nil
}
