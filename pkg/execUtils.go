package pkg

import (
	"bytes"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkmod"
	"strings"
)

// StartUpBifrost
func StartUpBifrost() error {
	// first, let's check to see if it's actually running...
	cfg := &atkmod.CliParts{
		Path: viper.GetString("podman.path"),
		Cmd:  "ps --format \"{{.Image}}\"",
	}

	logger.Debug("Checking to see if service is running...")
	podman := atkmod.NewPodmanCliCommandBuilder(cfg)
	runner := &atkmod.CliModuleRunner{*podman}
	out := new(bytes.Buffer)
	ctx := &atkmod.RunContext{
		Out: out,
	}
	runner.Run(ctx)

	if ctx.IsErrored() {
		return fmt.Errorf("error while trying to check for running service: %v", ctx.Errors)
	}

	logger.Debugf("Found running services: %v", out.String())

	if strings.TrimSpace(out.String()) == "\"localhost/bifrost:latest\"" {
		logger.Info("Found service; using.")
	} else {
		logger.Warn("Service not found; starting...")

		cfg = &atkmod.CliParts{
			Path:  viper.GetString("podman.path"),
			Flags: []string{"-d", "--rm"},
		}
		podman = atkmod.NewPodmanCliCommandBuilder(cfg).
			WithImage("localhost/bifrost")

		runner = &atkmod.CliModuleRunner{*podman}
		out = new(bytes.Buffer)
		runner.Run(&atkmod.RunContext{Out: out})
		logger.Debug(out)

		if ctx.IsErrored() {
			return fmt.Errorf("error while trying to start local service: %v", ctx.Errors)
		}
	}

	return nil
}
