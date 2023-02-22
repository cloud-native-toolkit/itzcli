package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloud-native-toolkit/atkmod"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
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
		return DeployModularSolution(cmd, args)
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

func DeployModularSolution(cmd *cobra.Command, args []string) error {
	// First, look up the base directory name from the environment, but if
	// it's not there, then default it to the home directory's cache folder.
	cacheDir, err := pkg.GetITZCacheDir()
	if err != nil {
		return err
	}
	solutionDir := filepath.Join(cacheDir, sol)
	logger.Tracef("Using solution directory: %s", solutionDir)

	// Now load up the file.
	loader := atkmod.NewAtkManifestFileLoader()
	manifest, err := loader.Load(filepath.Join(solutionDir, ".itz.yml"))

	if err != nil {
		logger.Warnf("error loading manifest file: %v", err)
		return fmt.Errorf("could not load .itz.yml")
	}
	if !manifest.IsSupported() {
		logger.Warnf("do not recognize manifest version: %s", manifest.ApiVersion)
		return fmt.Errorf("itz.yml version %s is not supported", manifest.ApiVersion)
	}

	logger.Tracef("Loaded manifest for %s; version %s", manifest.Kind, manifest.ApiVersion)
	outbuff := new(bytes.Buffer)
	errbuff := new(bytes.Buffer)
	runCtx := &atkmod.RunContext{
		Context: context.Background(),
		Out:     outbuff,
		Err:     errbuff,
	}
	// This will add the default mappings to /workspace to the current
	// solution directory so that the manifest file does not have to have a
	// bunch of repetitive mappings in it. It can, of course, be overriden but
	// by default the directory in which the manifest file is located should be
	// mounted to /workspace in the container.
	err = pkg.AddDefaultVolumeMappings(manifest, solutionDir)

	module := atkmod.NewDeployableModule(runCtx, manifest)
	logger.Tracef("Getting variables for solution %s", sol)
	hook := module.GetHook(atkmod.ListHook)
	err = hook(runCtx)
	if err != nil {
		return err
	}
	// Otherwise, do the prompter
	data, err := pkg.GetEventData(outbuff.String())
	if err != nil {
		return err
	}
	var sources = []pkg.VariableGetter{
		pkg.NewEnvvarGetter(),
		// TODO: Implement this. :)
		// pkg.NewStructGetter(nil),
	}
	resolver, err := pkg.NewVariableResolver(data.Variables, sources)
	if err != nil {
		return err
	}
	required := resolver.UnresolvedVars()
	prompter, err := pkg.NewVariablePrompter("Would you like to start?", required, true)
	if err != nil {
		return err
	}
	itr := prompter.Itr()
	for p := itr(); p != nil; p = itr() {
		logger.Tracef("Asking <%s>", p.String())
		err = prompt.Ask(p, os.Stdout, os.Stdin)
		if err != nil {
			return err
		}
	}

	// OK, now, so we'll write out the answers now
	pSrc := pkg.NewPromptGetter(*prompter)
	resolver.AddSource(pSrc)

	for i, v := range data.Variables {
		data.Variables[i].Value = resolver.GetString(v.Name)
	}

	event := pkg.NewEventWithData(string(atkmod.DeployLifecycleRequestEvent), "something", "itzcli", data)
	eventbuff := new(bytes.Buffer)
	atkmod.WriteEvent(&event, eventbuff)
	logger.Tracef("event: %s", eventbuff.String())

	return nil
}

// DeploySolution deploys the solution by handing it off to the bifrost
// API
func DeployLegacySolution(cmd *cobra.Command, args []string) error {
	// TODO: Eventually, it would be really neat to have some way of making
	// this be configurable, too. Or maybe this is just moved from here to a
	// container...
	var vars = make([]pkg.JobParam, 0)
	var resolver *pkg.BuildParamResolver

	prompterHandler := func(buf *bytes.Buffer) error {
		data, err := io.ReadAll(buf)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &vars)
		if err != nil {
			return err
		}

		// TODO: We should automagically resolve this instead of looking it up
		ocpCfg := viper.GetStringSlice("ocpnow.configFiles")
		if len(ocpCfg) == 0 {
			return fmt.Errorf("no OPC configuration found")
		}

		project, err := pkg.LoadProject(ocpCfg[0])
		if err != nil {
			return err
		}
		resolver, err = pkg.NewBuildParamResolver(project, cluster, vars)
		rootQuestion, err := resolver.BuildPrompter(sol)

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
		for k, v := range resolver.ResolvedParams() {
			if len(v) > 0 {
				tfvars = append(tfvars, pkg.JobParam{Name: k, Value: v})
			}
		}
		data, err := json.Marshal(tfvars)
		if err != nil {
			return err
		}
		buf.Write(data)
		return nil
	}

	err := pkg.DoContainerizedStep(cmd, "getcode", nil, nil)
	if err != nil {
		return err
	}

	err = pkg.DoContainerizedStep(cmd, "listparams", nil, prompterHandler)
	if err != nil {
		return err
	}

	err = pkg.DoContainerizedStep(cmd, "setparams", paramsHandler, nil)
	if err != nil {
		return err
	}

	return pkg.DoContainerizedStep(cmd, "applyall", nil, nil)

}
