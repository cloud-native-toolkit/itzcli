package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cloud-native-toolkit/atkmod"
	"io"
	"net/url"
	"regexp"
	"strings"
)

type ServiceType string

const (
	Background  ServiceType = "background"
	Interactive ServiceType = "interactive"
	InOut       ServiceType = "inout"
)

// ServiceConfig represents a configuration record for a service in the CLI's
// configuration file. The reason this struct is repeated compared to just using
// the atkmod.CliParts struct is because the layer of abstraction allows us to
// tweak the YAML structure.
type ServiceConfig struct {
	Env       []string    `yaml:"env,omitempty"`
	Image     string      `yaml:"image"`
	Local     bool        `yaml:"local"`
	MountOpts string      `yaml:"mountopts,omitempty"`
	Name      string      `yaml:"name,omitempty"`
	Type      ServiceType `yaml:"type,omitempty"`
	URL       *url.URL    `yaml:"url,omitempty"`
	Volumes   []string    `yaml:"volumes,omitempty"`
}

func createStatusRunner() *atkmod.CliModuleRunner {
	cfg := &atkmod.CliParts{
		Path: viper.GetString("podman.path"),
		Cmd:  "ps --format \"{{.Image}}\"",
	}
	cmd := atkmod.NewPodmanCliCommandBuilder(cfg)
	return &atkmod.CliModuleRunner{PodmanCliCommandBuilder: *cmd}
}

// ImageFound returns true if the name of the image was found in the
// output.
// TODO: create a different function for finding the exact image, or add a flag here...
func ImageFound(out *bytes.Buffer, name string) bool {
	logger.Tracef("Searching for image <%s> in output <%s>", name, out.String())
	scanner := bufio.NewScanner(out)
	img := strings.Split(name, ":")[0]

	for scanner.Scan() {
		line := scanner.Text()
		logger.Tracef("Checking for image <%s> in <%s>...", name, line)
		matched, _ := regexp.MatchString(`^\s*"?`+img+`(:(latest)|([a-z0-9-]+))?"?\s*`, line)
		if matched {
			logger.Tracef("Found image <%s> in line <%s>", name, line)
			return true
		}
	}

	return false
}

// WriteMessage writes the given message to the output writer.
func WriteMessage(msg string, w io.Writer) {
	_, err := fmt.Fprintf(w, "%s\n", msg)
	if err != nil {
		logger.Errorf("error trying to write message: \"%s\" (%v)", msg, err)
	}
}

// InputHandlerFunc is a handler for working with the input of a command.
type InputHandlerFunc func(in *bytes.Buffer) error

// OutputHandlerFunc is a handler for working with the output of a command.
type OutputHandlerFunc func(out *bytes.Buffer) error

// DoContainerizedStep is a function that acts as a Facade and does several
// operations. First, it loads the configuration for the cmd and the step from
// the CLI configuration file. It configures the CLI further using the ServiceConfig
// structure loaded from the CLI configuration file, then it sets up handler
// functions for handling STDIN and STDOUT in the container, then calls the Run
// method on the CliModuleRunner to run the container.
func DoContainerizedStep(cmd *cobra.Command, step string, inHandler InputHandlerFunc,
	outHandler OutputHandlerFunc) error {

	// The buffers for handling in and out
	out := new(bytes.Buffer)
	in := new(bytes.Buffer)

	logger.Debugf("Doing containerized step %s for command %s", step, cmd.Name())
	cfg, err := LoadServiceConfig(cmd, step)
	if err != nil {
		return err
	}

	runner, err := CreateCliRunner(cmd, cfg)

	if err != nil {
		return err
	}

	ctx := NewRunContext(cfg, cmd)

	if inHandler != nil {
		ctx.In = in
		err = inHandler(in)
		if err != nil {
			return err
		}
	}

	if outHandler != nil {
		ctx.Out = out
	}

	err = runner.Run(ctx)

	if outHandler != nil {
		return outHandler(out)
	}

	return err
}

// NewRunContext propertly creates the atkmod.RunContext for the given ServiceConfig
// and cobra.Command
func NewRunContext(svc *ServiceConfig, cmd *cobra.Command) *atkmod.RunContext {
	return &atkmod.RunContext{
		Out: cmd.OutOrStdout(),
		Err: cmd.ErrOrStderr(),
		In:  cmd.InOrStdin(),
	}
}

// CreateCliRunner creates an instance of a atkmod.CliModuleRunner for the given
// cobra.Command and ServiceConfig. The cobra.Command is used for variable
// substitution in the ServiceConfig. For example, you can use {{solution}} in
// the environment variables and it will substitute the value used for `--solution`
// on the command line.
func CreateCliRunner(cmd *cobra.Command, cfg *ServiceConfig) (*atkmod.CliModuleRunner, error) {

	parts := &atkmod.CliParts{
		Path: viper.GetString("podman.path"),
	}

	// Use the correct flag for the type of service.
	if cfg.Type == InOut {
		parts.Flags = append(parts.Flags, "-i")
	} else if cfg.Type == Interactive {
		parts.Flags = append(parts.Flags, "-it")
	} else if cfg.Type == Background {
		parts.Flags = append(parts.Flags, "-d")
	}

	cli := atkmod.NewPodmanCliCommandBuilder(parts).
		WithImage(cfg.Image)

	for _, val := range cfg.Volumes {
		vols := strings.Split(val, ":")
		cli.WithVolume(vols[0], strings.Join(vols[1:], ":"))
	}

	for _, val := range cfg.Env {
		envs := strings.Split(val, "=")
		resolved, err := ResolveInterpolation(cmd, envs[1])
		if err != nil {
			logger.Warnf("could not resolve variable: %s", envs[0])
		}
		cli.WithEnvvar(envs[0], resolved)
	}

	runner := &atkmod.CliModuleRunner{PodmanCliCommandBuilder: *cli}
	return runner, nil
}

// LoadServiceConfig loads the service configuration for the given command and
// path.
func LoadServiceConfig(cmd *cobra.Command, path string) (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	key := FlattenCommandName(cmd, path)

	err := viper.UnmarshalKey(key, &cfg, configOptions)
	if err != nil {
		return nil, err
	}
	logger.Tracef("Found configuration for key %s: %v", key, cfg)
	return cfg, nil
}

// ResolveInterpolation resolves the tokens in the string using the arguments
// configured in the cmd.
func ResolveInterpolation(cmd *cobra.Command, s string) (string, error) {
	if len(s) == 0 {
		return s, nil
	}
	r := regexp.MustCompile(`{{([^}]+)}}`)
	matches := r.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return s, nil
	}

	interpreted := s
	for _, match := range matches {
		logger.Tracef("looking up value of %s:", match[1])
		v := cmd.Flags().Lookup(match[1]).Value.String()
		logger.Tracef("found value of %s: %v", match[1], v)
		interpreted = strings.Replace(interpreted, match[0], v, 1)
	}
	return interpreted, nil
}

func configOptions(config *mapstructure.DecoderConfig) {
	config.ErrorUnused = false
	config.ErrorUnset = false
	config.IgnoreUntaggedFields = true
}
