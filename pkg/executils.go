package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkmod"
	"regexp"
	"strings"
)

// StartSvcImg
func StartSvcImg(imgName string, onPort string) error {
	// first, let's check to see if it's actually running...
	cfg := &atkmod.CliParts{
		Path: viper.GetString("podman.path"),
		Cmd:  "ps --format \"{{.Image}}\"",
	}

	logger.Debug("Checking to see if service is running...")
	cmd := atkmod.NewPodmanCliCommandBuilder(cfg)
	runner := &atkmod.CliModuleRunner{*cmd}
	out := new(bytes.Buffer)
	ctx := &atkmod.RunContext{
		Out: out,
	}
	runner.Run(ctx)

	if ctx.IsErrored() {
		return fmt.Errorf("error while trying to check for running service: %v", ctx.Errors)
	}

	logger.Debugf("Found running services: %v", out.String())

	if ImageFound(out, imgName) {
		logger.Infof("Found service; using service <%s> on port: %s", imgName, onPort)
	} else {
		logger.Warn("Service not found; starting...")

		cfg = &atkmod.CliParts{
			Path:  viper.GetString("podman.path"),
			Flags: []string{"-d", "--rm"},
		}

		if len(onPort) > 0 {
			cfg.Flags = append(cfg.Flags, "-p")
			cfg.Flags = append(cfg.Flags, fmt.Sprintf("%s:%s", onPort, "8080"))
		}

		cmd = atkmod.NewPodmanCliCommandBuilder(cfg).
			WithImage(imgName)

		cli, _ := cmd.Build()
		logger.Tracef("Using <%s> to start local service...", cli)
		runner = &atkmod.CliModuleRunner{*cmd}
		out = new(bytes.Buffer)
		runner.Run(&atkmod.RunContext{Out: out})
		logger.Debug(out)

		if ctx.IsErrored() {
			return fmt.Errorf("error while trying to start local service: %v", ctx.Errors)
		}
	}

	return nil
}

// ImageFound returns true if the name of the image was found in the
// output.
// TODO: we may want to look for the exact image
func ImageFound(out *bytes.Buffer, name string) bool {
	logger.Tracef("searching for image <%s> in output <%s>", name, out.String())
	scanner := bufio.NewScanner(out)
	img := strings.Split(name, ":")[0]

	for scanner.Scan() {
		line := scanner.Text()
		logger.Tracef("checking for image <%s> in <%s>...", name, line)
		matched, _ := regexp.MatchString(`^\s*"?`+img+`(:(latest)|([a-z0-9-]+))?"?\s*`, line)
		if matched {
			logger.Tracef("found for image <%s> in line <%s>", name, line)
			return true
		}
	}

	return false
}
