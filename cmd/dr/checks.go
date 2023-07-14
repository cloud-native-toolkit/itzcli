package dr

import "runtime"

// Create the checks for the configuration values that I know I'll need.

// SolutionsListPermissionsError is the permissions error printed by the doctor
// check.
const SolutionsListPermissionsError = `
Permissions error while trying to read from your list of solutions. The most
common cause is an expired or bad API token. You can resolve this issue by going
to https://builder.cloudnativetoolkit.dev/ to get your API token, save it in a 
file (e.g., /path/to/token.txt) and use the command:

    $ itz auth login --from-file /path/to/token.txt --service builder

`

// Linux returns true if the current system is a Linux system.
func Linux() bool {
	return runtime.GOOS == "linux"
}

// AllConfigChecks checks for configuration values on the system and defines
// defaulters for fixing the missing values if the user specifies --auto-fix.
var AllConfigChecks = []Check{
	NewConfigCheck("backstage.api.url", "", Static("https://catalog.techzone.ibm.com")),
	// The reservations configuration values
	NewConfigCheck("techzone.api.url", "", Static("https://api.techzone.ibm.com/api")),
	NewConfigCheck("reservations.api.path", "", Static("my/reservations/all")),
	NewConfigCheck("reservation.api.path", "", Static("reservation/ibmcloud-2")),
	// To address https://github.com/cloud-native-toolkit/itzcli/issues/25, check to see if the system
	// is a Linux system. If it is, then write the configuration file with the option in the volume mapping
	// that disables SELinux for the installer.
	NewConfigCheck("execute.workspace.ocpinstaller", "", IifStatic(Linux, DefaultOCPInstallerLinuxConfig, DefaultOCPInstallerConfig)),
}

// FileChecks defines the checks that are done for files on the system.
var FileChecks = []Check{
	// TODO: These need to be preserved in this order for now, but the checks
	// especially for files should be more stand-alone. But in order to do that,
	// the test coverage needs to be improved so that false errors aren't returned.
	NewReqConfigDirCheck("save"),
	NewReqConfigDirCheck("build_home"),
	NewReqConfigDirCheck("cache"),
	// This check will error out if the .itz directory does not exist, which is
	// created automatically by the checks above so long as the user specifies
	// the --auto-fix flag.
	NewFixableConfigFileCheck("cli-config.yaml", EmptyFileCreator),
	NewResourceFileCheck(OneExistsOnPath("podman", "docker"), "%s was not found on your path", UpdateConfig("podman.path")),
}

var ActionChecks = []Check{
	// This new action check will run if podman is the CLI tool of choice and if
	// the podman machine exists and is up. If it is, it will then fix the date
	NewCmdActionCheck("setting clock on podman machine", PodmanMachineExists(), UpdatePodmanMachineDate()),
}
