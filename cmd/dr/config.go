package dr

import (
	"path/filepath"
	"strings"

	"github.com/cloud-native-toolkit/itzcli/pkg"
)

var DefaultOCPInstallerConfig = &pkg.ServiceConfig{
	Name:  "ocp-installer",
	Local: true,
	Type:  "interactive",
	Image: "quay.io/ibmtz/ocpinstaller:stable",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "save"), "/usr/src/ocpnow/save"}, ":"),
	},
}

var DefaultSolutionDeployGetCode = &pkg.ServiceConfig{
	Name:  "download",
	Type:  "interactive",
	Image: "quay.io/ibmtz/downloader",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "cache"), "/git"}, ":"),
	},
	Env: []string{
		"ITZ_SOLUTION_ID={{solution}}",
	},
}

var DefaultSolutionDeployListParams = &pkg.ServiceConfig{
	Name:  "variables-get",
	Type:  "inout",
	Image: "quay.io/ibmtz/meta",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "cache"), "/workspace"}, ":"),
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
	Image: "quay.io/ibmtz/meta",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "cache"), "/workspace"}, ":"),
	},
	Env: []string{
		"ITZ_SOLUTION_ID={{solution}}",
		"ITZ_SOLUTION_META_ACTION=set-parameters",
		"ITZ_SOLUTION_CREDENTIALS_FILE=credentials.properties",
	},
}

var DefaultSolutionDeployApplyAll = &pkg.ServiceConfig{
	Name:  "deploy",
	Type:  "interactive",
	Image: "quay.io/ibmtz/deployer",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "cache"), "/techzone"}, ":"),
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "workspace"), "/workspace"}, ":"),
	},
	Env: []string{
		"ITZ_SOLUTION_ID={{solution}}",
	},
}
