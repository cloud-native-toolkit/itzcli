package cmd

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/api"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
)

var executeCmd = &cobra.Command{
	Use:    ExecuteAction,
	Short:  "Executes workspaces or pipelines in a cluster",
	Long:   "Executes workspaces or pipelines in a cluster",
	PreRun: SetLoggingLevel,
}

var pipelineURI string
var acceptDefaults bool

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

		// Try to parse the URL
		_, err := url.ParseRequestURI(pipelineURI)
		if err != nil {
			return fmt.Errorf("\"%s\" is not a valid URL", pipelineURI)
		}

		logger.Debugf("Executing pipeline from URI \"%s\"...", pipelineURI)

		client := &pkg.GitServiceClient{
			BaseDest: "/tmp",
		}
		pipeline, err := client.Get(pipelineURI)
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

		pipelineResolver := pkg.NewPipelineResolver(pipeline)
		if options.Includes(pkg.UsePipelineDefaults) {
			enabled = append(enabled, pipelineResolver)
		}

		chainedResolver := pkg.NewChainedResolver(options, enabled...)

		q, err := pkg.BuildPipelinePrompt(pipeline.Name, pipelineResolver, chainedResolver)
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

		return nil
	},
}

func init() {
	executePipelineCmd.Flags().StringVarP(&pipelineURI, "pipeline-url", "u", "", "The URL of the pipeline as YAML")
	executePipelineCmd.Flags().BoolVarP(&acceptDefaults, "accept-defaults", "d", false, "Accept defaults for pipeline parameters without asking")

	executeCmd.AddCommand(executeWorkspaceCmd)
	executeCmd.AddCommand(executePipelineCmd)
	executeCmd.AddCommand(executeApiCmd)

	rootCmd.AddCommand(executeCmd)
}
