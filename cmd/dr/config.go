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

var DefaultOCPInstallerLinuxConfig = &pkg.ServiceConfig{
	Name:  "ocp-installer",
	Local: true,
	Type:  "interactive",
	Image: "quay.io/ibmtz/ocpinstaller:stable",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "save"), "/usr/src/ocpnow/save", "Z"}, ":"),
	},
}

var DefaultPipelineCreatorConfig = &pkg.ServiceConfig{
	Name:  "pipeline-creator",
	Local: true,
	Type:  "interactive",
	Image: "quay.io/ibmtz/pipeline-creator:main",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "projects"), "/workspace"}, ":"),
	},
}

var DefaultPipelineConfigLinuxConfig = &pkg.ServiceConfig{
	Name:  "pipeline-creator",
	Local: true,
	Type:  "interactive",
	Image: "quay.io/ibmtz/pipeline-creator:main",
	Volumes: []string{
		strings.Join([]string{filepath.Join(pkg.MustITZHomeDir(), "projects"), "/workspace", "Z"}, ":"),
	},
}
