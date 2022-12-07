package test

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.ibm.com/skol/atkmod"
	"github.ibm.com/skol/itzcli/pkg"
	"testing"
)

func TestImageFound(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("localhost/bifrost:latest\n")
	out.WriteString("localhost/atkci:latest\n")
	assert.True(t, pkg.ImageFound(out, "localhost/bifrost"))
}

func TestImageNotFound(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("localhost/bifrost:latest\n")
	out.WriteString("localhost/atkci:latest\n")
	assert.False(t, pkg.ImageFound(out, "localhost/mooshoopork"))
}

func TestImageWithQuotes(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("\"localhost/bifrost:latest\"\n")
	out.WriteString("\"localhost/atkci:latest\"\n")
	assert.True(t, pkg.ImageFound(out, "localhost/bifrost"))
}

func TestImageWithLatestTag(t *testing.T) {
	out := new(bytes.Buffer)
	out.WriteString("\"localhost/bifrost:latest\"\n")
	out.WriteString("\"localhost/atkci:latest\"\n")
	assert.True(t, pkg.ImageFound(out, "localhost/bifrost:latest"))
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
