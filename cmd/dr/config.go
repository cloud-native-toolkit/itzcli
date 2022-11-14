package dr

import (
	"github.ibm.com/skol/itzcli/pkg"
	"path/filepath"
	"strings"
)

var DefaultOCPInstallerConfig = &pkg.ServiceConfig{
	Local: true,
	Type:  "interactive",
	Image: "quay.io/ibmtz/ocpinstaller:latest",
	Volumes: []string{
		strings.Join([]string{filepath.Join(MustITZHomeDir(), "save"), "/usr/src/ocpnow/save"}, ":"),
	},
}
