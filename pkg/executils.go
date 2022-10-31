package pkg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkmod"
	"io"
	"net/url"
	"regexp"
	"strings"
)

// ImgHandler handler for doing something with an image. It returns a true or
// false, which changes meaning slightly on how it's being used (PreStart, Start).
type ImgHandler func(svc *Service, runCtx *atkmod.RunContext, runner *atkmod.CliModuleRunner) bool

type VolumeMap map[string]string
type PortMap map[string]string
type Envvars map[string]string

// Service is a background service that is really a container that run
type Service struct {
	// DisplayName is the name that is displayed in the log messages and other
	// output, so should map to the business purpose of the container, such
	// as "builder" or "integration".
	DisplayName string
	// ImgName is the name of the image, such as "localhost/bifrost:latest"
	// It may include the label or not. If it includes the label, the exact
	// image will be used, giving us a means of pinning versions.
	ImgName string
	// IsLocal if the service is meant to be running locally. The rest of the
	// CLI is designed to use the service either way.
	IsLocal bool
	// URL is the URL of the image
	URL *url.URL
	// PreStart can be used to do anything before running the image, such as
	// checking the status of the service. If PreStart returns a false here,
	// then the service is not started. Errors and other messages will be in
	// the RunContext supplied to the handler.
	PreStart ImgHandler
	// Start is used to start the actual service (run the image). It returns
	// a true if the images was started successfully and a false if the image
	// was not started. A return of false does not necessarily mean the service
	// errored, though. Check the handler's RunContext for that.
	Start ImgHandler
	// PostStart, if defined, is an opportunity to do any additional setup
	// needed before starting the next image. In sequential execution (which is
	// currently the only one supported), this function will block before
	// continuing to next service.
	PostStart ImgHandler
	// Volumes are the local volumes to container volume mappings that are used
	// by the container and will be parsed. The key to the map is the local
	// volume name, the value is the container's volume name.
	Volumes VolumeMap
	// Ports are a map just like the volume mapping.
	Ports PortMap
	// Envvars are the environment variables for the container.
	Envvars Envvars
	// Flags are any additional flags required to start the container.
	Flags []string
	// MapToUID allows you to map the current user to a UID on the container.
	MapToUID int
}

// StatusHandler handles the default means of getting the status of a container
// using the given command line runner and the run context.
func StatusHandler(svc *Service, runCtx *atkmod.RunContext, runner *atkmod.CliModuleRunner) bool {
	runCtx.Log.Debugf("Checking to see if %s service is running...", svc.DisplayName)
	out := new(bytes.Buffer)
	localCtx := &atkmod.RunContext{
		Out: out,
	}
	err := runner.Run(localCtx)
	if err != nil || localCtx.IsErrored() {
		return false
	}
	runCtx.Log.Debugf("Found running services: %v", out)

	return ImageFound(out, svc.ImgName)
}

// StartHandler handles the default means of starting a container using the given
// command line runner and the run context.
func StartHandler(svc *Service, runCtx *atkmod.RunContext, runner *atkmod.CliModuleRunner) bool {
	runCtx.Log.Infof("Starting %s service...", svc.DisplayName)
	out := new(bytes.Buffer)
	localCtx := &atkmod.RunContext{
		Out: out,
	}
	cmdStr, _ := runner.Build()
	runCtx.Log.Tracef("Using command <%s> to start %s service...", cmdStr, svc.DisplayName)
	err := runner.Run(localCtx)
	if err != nil || localCtx.IsErrored() {
		return false
	}
	return true
}

type ServiceHandlingPolicy string

const (
	Sequential ServiceHandlingPolicy = "sequential"
	Parallel   ServiceHandlingPolicy = "parallel"
)

// StartupServices handles the status and starting of the necessary services.
// It takes a ServiceHandlingPolicy
func StartupServices(ctx *atkmod.RunContext, svcs []Service, policy ServiceHandlingPolicy) error {
	if policy == Sequential {
		for _, svc := range svcs {
			// First, check to see if the service is already started...
			isStarted := svc.PreStart(&svc, ctx, createStatusRunner())
			if !isStarted {
				ctx.Log.Warnf("%s service not found; starting...", svc.DisplayName)
				ok := svc.Start(&svc, ctx, createStartRunner(svc))
				if !ok || ctx.IsErrored() {
					return fmt.Errorf("error while trying to start service %s: %v", svc.DisplayName, ctx.Errors)
				}
				if svc.PostStart != nil {
					ok = svc.PostStart(&svc, ctx, nil)
					if !ok || ctx.IsErrored() {
						return fmt.Errorf("error handling post start for service: %s", svc.DisplayName)
					}
				}
			} else {
				ctx.Log.Infof("Found %s service; using service <%s> on port: %s", svc.DisplayName, svc.ImgName, getPort(svc.URL))
			}
		}
	} else {
		return errors.New("only sequential execution is supported")
	}
	return nil
}

func getPort(uri *url.URL) string {
	return uri.Port()
}

func createStatusRunner() *atkmod.CliModuleRunner {
	cfg := &atkmod.CliParts{
		Path: viper.GetString("podman.path"),
		Cmd:  "ps --format \"{{.Image}}\"",
	}
	cmd := atkmod.NewPodmanCliCommandBuilder(cfg)
	return &atkmod.CliModuleRunner{PodmanCliCommandBuilder: *cmd}
}

func createStartRunner(svc Service) *atkmod.CliModuleRunner {
	cfg := &atkmod.CliParts{
		Path: viper.GetString("podman.path"),
		// in service (daemon) mode...
		Flags: svc.Flags,
	}
	localPort := getPort(svc.URL)
	cmd := atkmod.NewPodmanCliCommandBuilder(cfg).
		WithImage(svc.ImgName)

	if len(localPort) > 0 {
		// HACK: this should probably be configurable, but for now we know that
		// both services (containers) expose their stuff on port 8080, but
		// that needs to be mapped to the port the I expect from the configuration.
		cmd.WithPort(localPort, "8080")
	}

	for key, val := range svc.Volumes {
		cmd.WithVolume(key, val)
	}

	for key, val := range svc.Envvars {
		cmd.WithEnvvar(key, val)
	}

	if svc.MapToUID > 0 {
		cmd.WithUserMap(0, svc.MapToUID, 1)
		cmd.WithUserMap(1, 0, svc.MapToUID)
	}

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
func WriteMessage(msg string, w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s\n", msg)
	return err
}
