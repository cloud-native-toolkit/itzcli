package pkg

import (
	"bytes"
	"fmt"
	"github.com/cloud-native-toolkit/atkmod"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/mitchellh/mapstructure"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/url"
	"os"
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

// WriteMessage writes the given message to the output writer.
func WriteMessage(msg string, w io.Writer) {
	_, err := fmt.Fprintf(w, "%s\n", msg)
	if err != nil {
		logger.Errorf("error trying to write message: \"%s\" (%v)", msg, err)
	}
}

// VariableGetter is a function that returns the value of a variable specified
// by key.
type VariableGetter func(key string) (string, bool)

func NewCollectionGetter(source []atkmod.EventDataVarInfo) VariableGetter {
	return func(key string) (string, bool) {
		for _, v := range source {
			if v.Name == key {
				return v.Value, true
			}
		}
		return "", false
	}
}

// NewEnvvarGetter uses the OS environment to get the value of the variable.
func NewEnvvarGetter() VariableGetter {
	return func(key string) (string, bool) {
		return os.LookupEnv(key)
	}
}

// NewPromptGetter creates a getter that uses the prompt answers for a source.
func NewPromptGetter(source prompt.Prompt) VariableGetter {
	return func(key string) (string, bool) {
		val, ok := source.VarMap()[key]
		return val, ok
	}
}

// NewStructGetter creates a getter that will look up the variables in the given
// source.
func NewStructGetter(source interface{}) VariableGetter {
	envals, err := ResolveVars(source, nil)
	if err != nil {
		return func(key string) (string, bool) {
			return "", false
		}
	}
	vars := NewEventDataVarInfoSlice(envals)
	return NewCollectionGetter(vars)
}

// NewEventDataVarInfoSlice creates a slice of EventDataVarInfo from the given
// map of key-value pairs.
func NewEventDataVarInfoSlice(envvars map[string]string) []atkmod.EventDataVarInfo {
	var result []atkmod.EventDataVarInfo
	for k, v := range envvars {
		result = append(result, atkmod.EventDataVarInfo{
			Name:  k,
			Value: v,
		})
	}
	return result
}

// VariableResolver is an structure for resolving variables.
type VariableResolver struct {
	requiredVars []atkmod.EventDataVarInfo
	sources      []VariableGetter
}

func (r *VariableResolver) AddSource(source VariableGetter) {
	r.sources = append(r.sources, source)
}

func (r *VariableResolver) GetString(key string) string {
	result, ok := r.LookupString(key)
	if ok {
		return result
	}
	return ""
}

func (r *VariableResolver) LookupString(key string) (string, bool) {
	for _, source := range r.sources {
		val, exists := source(key)
		if exists {
			return val, true
		}
	}
	return "", false
}

// UnresolvedVars looks through the potential sources to see what variables are
// still required.
func (r *VariableResolver) UnresolvedVars() []atkmod.EventDataVarInfo {
	// This is going to actually do the resolution
	var unresolved = make([]atkmod.EventDataVarInfo, 0)
	for _, v := range r.requiredVars {
		val, found := r.LookupString(v.Name)
		// look up the value in each source, and if it is found in any
		// source and found to not be empty, we do not add it to the list
		if !found || len(strings.TrimSpace(val)) == 0 {
			unresolved = append(unresolved, v)
		}
	}
	return unresolved
}

// NewVariableResolver creates a new VariableResolver
func NewVariableResolver(required []atkmod.EventDataVarInfo, sources []VariableGetter) (*VariableResolver, error) {
	return &VariableResolver{
		requiredVars: required,
		sources:      sources,
	}, nil
}

// NewVariablePrompter builds a prompter that will prompt the user for the
// values of the required variables, using the provided q as the initial top
// level question (e.g., "Would you like to continue?".
func NewVariablePrompter(q string, required []atkmod.EventDataVarInfo, includeDefaults bool) (*prompt.Prompt, error) {
	builder := prompt.NewPromptBuilder()

	rootQuestion, err := builder.Path("proceed").
		Text(q).
		WithOptions(prompt.YesNo()).
		Build()

	if err != nil {
		return nil, err
	}

	for _, v := range required {
		logger.Tracef("Building prompt for <%s>", v)
		b := prompt.NewPromptBuilder().
			Path(v.Name).
			Textf("What value would you like to use for '%s'?", v.Name)

		if len(v.Default) > 0 {
			if !includeDefaults {
				continue
			}
			b.WithDefaultValue(v.Default)
		}
		subP, err := b.Build()
		if err != nil {
			return nil, err
		}
		rootQuestion.AddSubPrompt(subP)
	}
	return rootQuestion, nil
}

func addVolToImage(img *atkmod.ImageInfo, dir string) error {
	mountPath := viper.GetString("workspace.dir")
	if len(mountPath) == 0 {
		mountPath = "/workspace"
	}
	mountExists := func(v []atkmod.VolumeInfo, path string) bool {
		for _, m := range v {
			if m.MountPath == path {
				return true
			}
		}
		return false
	}

	if img == nil {
		return nil
	}

	v := img.Volumes
	if len(v) == 0 || !mountExists(v, mountPath) {
		img.Volumes = append(img.Volumes, atkmod.VolumeInfo{
			Name:      dir,
			MountPath: mountPath,
		})
	}

	return nil
}

func appendIfNotNil(slice []error, err error) []error {
	if err != nil {
		return append(slice, err)
	}
	return slice
}

func AddDefaultVolumeMappings(manifest *atkmod.ModuleInfo, dir string) error {
	errs := make([]error, 0)
	errs = appendIfNotNil(errs, addVolToImage(&manifest.Specifications.Hooks.List, dir))
	errs = appendIfNotNil(errs, addVolToImage(&manifest.Specifications.Hooks.GetState, dir))
	errs = appendIfNotNil(errs, addVolToImage(&manifest.Specifications.Hooks.Validate, dir))
	errs = appendIfNotNil(errs, addVolToImage(&manifest.Specifications.Lifecycle.PreDeploy, dir))
	errs = appendIfNotNil(errs, addVolToImage(&manifest.Specifications.Lifecycle.Deploy, dir))
	errs = appendIfNotNil(errs, addVolToImage(&manifest.Specifications.Lifecycle.PostDeploy, dir))
	if len(errs) > 0 {
		return fmt.Errorf("failed to add default volume mappings: %v", errs)
	}
	return nil
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

	ctx := NewRunContext(cmd)

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

// NewRunContext properly creates the atkmod.RunContext for the given
// cobra.Command
func NewRunContext(cmd *cobra.Command) *atkmod.RunContext {
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
