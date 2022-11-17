package dr

import "path/filepath"

// Create the checks for the configuration values that I know I'll need.
const SolutionsListPermissionsError = `
Permissions error while trying to read from your list of solutions. The most
common cause is an expired or bad API token. You can resolve this issue by going
to https://builder.cloudnativetoolkit.dev/ to get your API token, save it in a 
file (e.g., /path/to/token.txt) and use the command:

    $ itz auth login --from-file /path/to/token.txt --service builder

`

// AllConfigChecks checks for configuration values on the system and defines
// defaulters for fixing the missing values if the user specifies --auto-fix.
var AllConfigChecks = []Check{
	// The bifrost configuration values
	NewConfigCheck("bifrost.api.image", "", Static("quay.io/ocpnow/bifrost:latest")),
	NewConfigCheck("bifrost.api.local", "", Static(true)),
	NewConfigCheck("bifrost.api.url", "", ServiceURL("http", 8088)),
	// The builder configuration values
	NewConfigCheck("builder.api.token", "", Messager(SolutionsListPermissionsError)),
	NewConfigCheck("builder.api.url", "", Static("https://ascent-bff-mapper-staging.dev-mapper-ocp-4be51d31de4db1e93b5ed298cdf6bb62-0000.eu-de.containers.appdomain.cloud")),
	NewConfigCheck("builder.api.username", "", Prompter("Please enter your ibm.com email address (for getting your solutions):")),
	// The Jenkins (ci) configuration values
	NewConfigCheck("ci.api.image", "", Static("quay.io/ocpnow/atkci:latest")),
	NewConfigCheck("ci.api.local", "", Static(true)),
	NewConfigCheck("ci.api.password", "", RandomVal(24)),
	NewConfigCheck("ci.api.url", "", ServiceURL("http", 8080)),
	NewConfigCheck("ci.api.user", "", Static("bifrost")),
	NewConfigCheck("ci.buildtoken", "", Static("68b2a8e9ffe5a395d839d9bf87db6800")),
	NewConfigCheck("ci.localdir", "", ConfigDir("build_home")),
	NewConfigCheck("ci.mountOpts", "", Static(":Z")),
	// The reservations configuration values
	NewConfigCheck("reservations.api.token", "", Prompter("There is no token defined for the Reservations API. Please provide one:")),
	NewConfigCheck("reservations.api.url", "", Static("https://api.techzone.ibm.com/api/my/reservations/all")),
	NewConfigCheck("workspaces.ocpinstaller", "", Static(DefaultOCPInstallerConfig)),
}

// FileChecks defines the checks that are done for files on the system.
var FileChecks = []Check{
	NewBinaryFileCheck("podman", "%s was not found on your path", UpdateConfig("podman.path")),
	NewReqConfigDirCheck("build_home"),
	NewReqConfigDirCheck("save"),
	NewFixableConfigFileCheck("cli-config.yaml", EmptyFileCreator),
	NewFixableConfigFileCheck(filepath.Join("build_home", "casc.yaml"), TemplatedFileCreator(CascTemplateString)),
}
