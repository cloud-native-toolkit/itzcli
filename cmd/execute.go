package cmd

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/api"
	"github.com/cloud-native-toolkit/itzcli/cmd/dr"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"net/url"
)

var executeCmd = &cobra.Command{
	Use:    ExecuteAction,
	Short:  "Executes workspaces or pipelines in a cluster",
	Long:   "Executes workspaces or pipelines in a cluster",
	PreRun: SetLoggingLevel,
}

// TODO: Change to true for the default in release 1.0
const DefaultUseContainer = false

var pipelineURI string
var pipelineRunURI string
var acceptDefaults bool
var useContainer bool
var clusterURL string
var clusterUsername string
var clusterPassword string

var executeApiCmd = &cobra.Command{
	Use:    ApiResource,
	Hidden: true,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		api.StartServer(rootCmd)
		return nil
	},
}

var executeWorkspaceCmd = &cobra.Command{
	Use:    WorkspaceResource,
	Short:  "Executes the given workspace",
	Long:   "Executes the given workspace",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Running command: %s", args[0])
		workspace := args[0]
		key := pkg.Keyify(workspace)
		if viper.Get(pkg.FlattenCommandName(cmd, key)) == nil {
			return fmt.Errorf("workspace does not exist: %s", workspace)
		}
		return pkg.DoContainerizedStep(cmd, key, nil, nil)
	},
}

var executePipelineCmd = &cobra.Command{
	Use:    PipelineResource,
	Short:  "Executes the given pipeline",
	Long:   "Executes the given pipeline",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {

		logger.Tracef("Using args: %v", args)

		if len(pipelineURI) == 0 {
			return fmt.Errorf("you must specify a URL for the pipeline to execute")
		}

		if len(pipelineRunURI) == 0 {
			return fmt.Errorf("you must specify a URL for the pipeline run to use")
		}

		// Try to parse the URL
		_, err := url.ParseRequestURI(pipelineURI)
		if err != nil {
			return fmt.Errorf("\"%s\" is not a valid URL", pipelineURI)
		}

		logger.Debugf("Executing pipeline from URI \"%s\" with pipeline run \"%s\"...", pipelineURI, pipelineRunURI)

		client := &pkg.GitServiceClient{
			BaseDest: "/tmp",
		}
		pipeline, err := client.Get(pipelineURI, pkg.UnmarshalPipeline)
		pl := pipeline.(*v1beta1.Pipeline)
		if err != nil {
			return fmt.Errorf("error trying to get pipeline at \"%s\": %v", pipelineURI, err)
		}
		options := pkg.DefaultParseOptions
		enabled := make([]pkg.ParamResolver, 0)

		if len(args) > 0 {
			options = options | pkg.UseCommandLineArgs
		}

		if acceptDefaults {
			options = options | pkg.UsePipelineDefaults
		}

		if options.Includes(pkg.UseEnvironmentVars) {
			enabled = append(enabled, pkg.NewEnvParamResolver())
		}

		if options.Includes(pkg.UseCommandLineArgs) {
			enabled = append(enabled, pkg.NewArgsParamParser(args))
		}

		pipelineResolver := pkg.NewPipelineResolver(pl)
		if options.Includes(pkg.UsePipelineDefaults) {
			enabled = append(enabled, pipelineResolver)
		}

		chainedResolver := pkg.NewChainedResolver(options, enabled...)

		q, err := pkg.BuildPipelinePrompt(pl.Name, pipelineResolver, chainedResolver)
		if err != nil {
			return err
		}

		nextPrompter := q.Itr()

		for p := nextPrompter(); p != nil; p = nextPrompter() {
			logger.Tracef("Asking <%s>", p.String())
			err = prompt.Ask(p, cmd.OutOrStdout(), cmd.InOrStdin())
			if err != nil {
				return err
			}
		}

		promptResolver := pkg.NewPromptResolver(q)
		// Now that we have all the answers, let's create a new resolve to read them
		// and add it to the chained resolver so the pipeline runner can build out the
		// answers
		chainedResolver.AddResolver(promptResolver)

		pRun, err := client.Get(pipelineRunURI, pkg.UnmarshalPipelineRun)
		if err != nil {
			return err
		}
		pr := pRun.(*v1beta1.PipelineRun)
		// Now the pipeline run is updated with the new context and answers from the user...
		updated, err := pkg.MergePipelineRun(pr, pl, pipelineResolver, chainedResolver)
		if err != nil {
			return err
		}
		updated.GenerateName = ""
		updated.Name = fmt.Sprintf("%s-run", pl.Name)

		return pkg.ExecPipelineRun(pl, updated, dr.RunScript, useContainer, pkg.ClusterInfo{URL: clusterURL}, pkg.CredInfo{Name: clusterUsername, ApiKey: clusterPassword}, cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func init() {
	executePipelineCmd.Flags().StringVarP(&pipelineURI, "pipeline-url", "p", "", "The URL of the pipeline as YAML")
	executePipelineCmd.Flags().StringVarP(&pipelineRunURI, "pipeline-run-url", "r", "", "The URL of the pipeline run as YAML")
	executePipelineCmd.Flags().BoolVarP(&acceptDefaults, "accept-defaults", "d", false, "Accept defaults for pipeline parameters without asking")
	executePipelineCmd.Flags().BoolVarP(&useContainer, "use-container", "c", DefaultUseContainer, "If true, the commands run in a container")

	executePipelineCmd.Flags().StringVarP(&clusterURL, "cluster-api-url", "a", "", "The URL of the target cluster")
	executePipelineCmd.Flags().StringVarP(&clusterUsername, "cluster-username", "u", "", "A username to login to the target cluster")
	executePipelineCmd.Flags().StringVarP(&clusterPassword, "cluster-password", "P", "", "A password to login to the target cluster")

	executeCmd.AddCommand(executeWorkspaceCmd)
	executeCmd.AddCommand(executePipelineCmd)
	executeCmd.AddCommand(executeApiCmd)

	rootCmd.AddCommand(executeCmd)
}
