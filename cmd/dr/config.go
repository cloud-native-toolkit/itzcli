package dr

import (
	"github.ibm.com/skol/itzcli/pkg"
	"path/filepath"
	"strings"
)

var DefaultOCPInstallerConfig = &pkg.ServiceConfig{
	Name:  "ocp-installer",
	Local: true,
	Type:  "interactive",
	Image: "quay.io/ibmtz/ocpinstaller:latest",
	Volumes: []string{
		strings.Join([]string{filepath.Join(MustITZHomeDir(), "save"), "/usr/src/ocpnow/save"}, ":"),
	},
}

var DefaultSolutionDeployGetCode = &pkg.ServiceConfig{
	Name:  "solution-download",
	Type:  "background",
	Image: "quay.io/podman/hello",
	Volumes: []string{
		strings.Join([]string{filepath.Join(MustITZHomeDir(), "cache"), "/git"}, ":"),
	},
	Env: []string{
		"ITZ_SOLUTION_ID={{solution}}",
	},
}

var DefaultSolutionDeployListParams = &pkg.ServiceConfig{
	Name:  "variables-get",
	Type:  "inout",
	Image: "itz-plugins/solution-meta",
	Volumes: []string{
		strings.Join([]string{filepath.Join(MustITZHomeDir(), "cache"), "/workspace"}, ":"),
	},
	Env: []string{
		"ITZ_SOLUTION_ID={{solution}}",
		"ITZ_SOLUTION_META_ACTION=list-parameters",
		"ITZ_SOLUTION_CREDENTIALS_FILE=credentials.template",
	},
}

var DefaultSolutionDeploySetParams = &pkg.ServiceConfig{
	Name:  "variables-save",
	Type:  "inout",
	Image: "itz-plugins/solution-meta",
	Volumes: []string{
		strings.Join([]string{filepath.Join(MustITZHomeDir(), "cache"), "/workspace"}, ":"),
	},
	Env: []string{
		"ITZ_SOLUTION_ID={{solution}}",
		"ITZ_SOLUTION_META_ACTION=list-parameters",
		"ITZ_SOLUTION_CREDENTIALS_FILE=credentials.properties",
	},
}

var DefaultSolutionDeployApplyAll = &pkg.ServiceConfig{
	Name:  "deployer",
	Type:  "interactive",
	Image: "itz-plugins/solution-deployer",
	Volumes: []string{
		strings.Join([]string{filepath.Join(MustITZHomeDir(), "cache"), "/techzone"}, ":"),
		strings.Join([]string{filepath.Join(MustITZHomeDir(), "workspace"), "/workspace"}, ":"),
	},
	Env: []string{
		"ITZ_SOLUTION_ID={{solution}}",
	},
}
