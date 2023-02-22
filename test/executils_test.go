package test

import (
	"github.com/cloud-native-toolkit/atkmod"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequiredVarsWithEnv(t *testing.T) {
	t.Setenv("TF_VAR_itz_test_1", "this variable has a value")
	t.Setenv("TF_VAR_itz_test_2", "")
	requiredVars := []atkmod.EventDataVarInfo{
		{
			Name:    "TF_VAR_itz_test_1",
			Default: "This has a value",
		},
		{
			Name: "TF_VAR_itz_test_2",
		},
		{
			Name: "TF_VAR_itz_test_3",
		},
	}
	sources := []pkg.VariableGetter{
		pkg.NewEnvvarGetter(),
	}
	resolver, err := pkg.NewVariableResolver(requiredVars, sources)
	assert.NoError(t, err)
	unresolved := resolver.UnresolvedVars()
	assert.Len(t, unresolved, 2)
	assert.Contains(t, unresolved, requiredVars[1])
	assert.Contains(t, unresolved, requiredVars[2])
}

func TestRequiredVarsWithEnvAndCollection(t *testing.T) {
	t.Setenv("TF_VAR_itz_test_4", "this variable has a value")
	collectionVars := []atkmod.EventDataVarInfo{
		{
			Name:  "TF_VAR_itz_test_6",
			Value: "some value",
		},
	}

	requiredVars := []atkmod.EventDataVarInfo{
		{
			Name:    "TF_VAR_itz_test_4",
			Default: "This has a value",
		},
		{
			Name: "TF_VAR_itz_test_5",
		},
		{
			Name: "TF_VAR_itz_test_6",
		},
	}
	sources := []pkg.VariableGetter{
		pkg.NewCollectionGetter(collectionVars),
		pkg.NewEnvvarGetter(),
	}
	resolver, err := pkg.NewVariableResolver(requiredVars, sources)
	assert.NoError(t, err)
	unresolved := resolver.UnresolvedVars()
	assert.Len(t, unresolved, 1)
	assert.Contains(t, unresolved, requiredVars[1])
}

func TestRequiredVarsWithObject(t *testing.T) {
	// Set up an example object with some values
	t.Setenv("TF_VAR_itz_test_8", "this variable has a value")
	credInfo := &pkg.CredInfo{
		Name:   "somenamehere",
		ApiKey: "thisisanapikeyexpected",
	}
	requiredVars := []atkmod.EventDataVarInfo{
		{
			Name: "TF_VAR_ibmcloud_api_key",
		},
		{
			Name: "TF_VAR_itz_test_8",
		},
	}
	sources := []pkg.VariableGetter{
		pkg.NewStructGetter(credInfo),
		pkg.NewEnvvarGetter(),
	}
	resolver, err := pkg.NewVariableResolver(requiredVars, sources)
	assert.NoError(t, err)
	unresolved := resolver.UnresolvedVars()
	assert.Len(t, unresolved, 0)
}

func TestPromptBuilderWithDefaults(t *testing.T) {
	requiredVars := []atkmod.EventDataVarInfo{
		{
			Name:    "TF_VAR_itz_test_1",
			Default: "This has a value",
		},
		{
			Name: "TF_VAR_itz_test_2",
		},
		{
			Name: "TF_VAR_itz_test_3",
		},
	}
	prompter, err := pkg.NewVariablePrompter("Would you like to see how cool this is?", requiredVars, true)
	assert.NoError(t, err)
	assert.Equal(t, "Would you like to see how cool this is?", prompter.String())
	itr := prompter.Itr()
	p := itr()
	assert.Equal(t, "Would you like to see how cool this is?", prompter.String())
	p.Record("Yes")

	p = itr()
	assert.NotNil(t, p)
	assert.Equal(t, "What value would you like to use for 'TF_VAR_itz_test_1'?", p.String())
	p.Record("my answer for test 1")

	p = itr()
	assert.NotNil(t, p)
	assert.Equal(t, "What value would you like to use for 'TF_VAR_itz_test_2'?", p.String())
	p.Record("my answer for test 2")

	p = itr()
	assert.NotNil(t, p)
	assert.Equal(t, "What value would you like to use for 'TF_VAR_itz_test_3'?", p.String())
	p.Record("my answer for test 3")

	p = itr()
	// We've run through all the questions, so there should be no more from the iterator
	assert.Nil(t, p)
}

func TestResolveInterpolation(t *testing.T) {
	cmd := createDeploymentCmd(t)
	expected := "ITZ_SOLUTION_ID=hello-world"
	actual, err := pkg.ResolveInterpolation(cmd, "ITZ_SOLUTION_ID={{solution}}")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestResolveInterpolationWithNone(t *testing.T) {
	cmd := createDeploymentCmd(t)
	expected := "ITZ_SOLUTION_ID=hello-world"
	actual, err := pkg.ResolveInterpolation(cmd, expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddDefaultVolumeMappings(t *testing.T) {
	mod := &atkmod.ModuleInfo{
		ApiVersion: "itzcli/v1alpha1",
		Kind:       "InstallManifest",
		Metadata: atkmod.MetadataInfo{
			Name: "test",
		},
		Specifications: atkmod.SpecInfo{
			Hooks: atkmod.HookInfo{
				List: atkmod.ImageInfo{
					Image: "itz-hook-tf-list",
				},
				GetState: atkmod.ImageInfo{
					Image: "itz-hook-tf-get-state",
					Volumes: []atkmod.VolumeInfo{
						{
							MountPath: "/app",
							Name:      "/var",
						},
					},
				},
			},
		},
	}
	assert.NotNil(t, mod)
	assert.Equalf(t, 0, len(mod.Specifications.Hooks.List.Volumes), "expected no volume mapping")
	assert.Equalf(t, 1, len(mod.Specifications.Hooks.GetState.Volumes), "expected exactly one volume mapping")
	// Now add the default volume mappings
	err := pkg.AddDefaultVolumeMappings(mod, "/tmp")
	assert.NoError(t, err)
	assert.Equalf(t, 1, len(mod.Specifications.Hooks.List.Volumes), "expected exactly one volume mapping")
	assert.Equalf(t, 2, len(mod.Specifications.Hooks.GetState.Volumes), "expected exactly two volume mappings")
	mapping := mod.Specifications.Hooks.List.Volumes[0]
	assert.Equalf(t, "/workspace", mapping.MountPath, "expected container path to be /workspace")
	assert.Equalf(t, "/tmp", mapping.Name, "expected container path to be /workspace")

	// So, I should be able to do this again, and because now there is already a mapping with the /workspace path, the
	// code should not add another one.
	err = pkg.AddDefaultVolumeMappings(mod, "/moo")

	assert.NoError(t, err)
	assert.Equalf(t, 1, len(mod.Specifications.Hooks.List.Volumes), "expected exactly one volume mapping")
	assert.Equalf(t, "/workspace", mapping.MountPath, "expected container path to be /workspace")
	assert.Equalf(t, "/tmp", mapping.Name, "expected container path to be /workspace")
}

func TestCreateCliRunner(t *testing.T) {
	expParts := &atkmod.CliParts{
		Path:  viper.GetString("podman.path"),
		Flags: []string{"-d"},
	}
	expBldr := atkmod.NewPodmanCliCommandBuilder(expParts)
	expBldr.WithImage("localhost/bifrost:latest")

	expRunner := &atkmod.CliModuleRunner{PodmanCliCommandBuilder: *expBldr}
	cmd := createDeploymentCmd(t)
	cfg := &pkg.ServiceConfig{
		Image: "localhost/bifrost:latest",
		Type:  pkg.Background,
	}
	actual, err := pkg.CreateCliRunner(cmd, cfg)
	assert.NoError(t, err)
	assert.Equal(t, expRunner, actual)
}

func TestCreateCliRunnerWithVols(t *testing.T) {
	expParts := &atkmod.CliParts{
		Path:  viper.GetString("podman.path"),
		Image: "localhost/bifrost:latest",
		Flags: []string{"-i"},
	}
	expBldr := atkmod.NewPodmanCliCommandBuilder(expParts)
	expBldr.WithImage("localhost/bifrost:latest")
	expBldr.WithVolume("/home/test/workspace", "/workspace")

	expRunner := &atkmod.CliModuleRunner{PodmanCliCommandBuilder: *expBldr}
	cmd := createDeploymentCmd(t)
	cfg := &pkg.ServiceConfig{
		Image: "localhost/bifrost:latest",
		Type:  pkg.InOut,
		Volumes: []string{
			"/home/test/workspace:/workspace",
		},
	}
	actual, err := pkg.CreateCliRunner(cmd, cfg)
	assert.NoError(t, err)
	assert.Equal(t, expRunner, actual)
}

func TestResolveInterpolationWithEmptyString(t *testing.T) {
	cmd := createDeploymentCmd(t)
	expected := ""
	actual, err := pkg.ResolveInterpolation(cmd, "")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func createDeploymentCmd(t *testing.T) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "root",
	}

	deploySolutionCmd := &cobra.Command{
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
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	var fn string
	var sol string
	var cluster string
	var rez string
	var useCached bool

	deploySolutionCmd.Flags().StringVarP(&fn, "file", "f", "", "The full path to the solution file to be deployed.")
	deploySolutionCmd.Flags().StringVarP(&sol, "solution", "s", "", "The name of the solution to be deployed.")
	deploySolutionCmd.Flags().StringVarP(&cluster, "cluster-name", "c", "", "The name of the cluster created by ocpnow to target.")
	deploySolutionCmd.Flags().StringVarP(&rez, "reservation", "r", "", "The id of the reservation to target.")
	// TODO: Change this from true to false by default
	deploySolutionCmd.Flags().BoolVarP(&useCached, "use-cache", "u", false, "If true, uses a cached solution file instead of downloading from target.")
	deploySolutionCmd.MarkFlagsMutuallyExclusive("file", "solution")
	deploySolutionCmd.MarkFlagsMutuallyExclusive("reservation", "cluster-name")

	err := deploySolutionCmd.ParseFlags([]string{"--solution", "hello-world", "--cluster-name", "mock"})
	assert.NoError(t, err)

	rootCmd.AddCommand(deploySolutionCmd)

	return deploySolutionCmd
}
