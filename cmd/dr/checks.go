package dr

import "path/filepath"

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

// AllConfigChecks checks for configuration values on the system and defines
// defaulters for fixing the missing values if the user specifies --auto-fix.
var AllConfigChecks = []Check{
	NewConfigCheck("backstage.api.url", "", Static("https://atk-backstage.us7aimk1t6d.us-south.codeengine.appdomain.cloud")),
	// The reservations configuration values
	NewConfigCheck("reservations.api.token", "", Static("There is no token defined for the Reservations API. Please run the auth login command.")),
	NewConfigCheck("reservations.api.url", "", Static("https://api.techzone.ibm.com/api/my/reservations/all")),
	NewConfigCheck("reservation.api.url", "", Static("https://api.techzone.ibm.com/api/reservation/ibmcloud-2/")),
	NewConfigCheck("itz.workspace.ocpinstaller", "", Static(DefaultOCPInstallerConfig)),
	NewConfigCheck("solution.deploy.getcode", "", Static(DefaultSolutionDeployGetCode)),
	NewConfigCheck("solution.deploy.listparams", "", Static(DefaultSolutionDeployListParams)),
	NewConfigCheck("solution.deploy.setparams", "", Static(DefaultSolutionDeploySetParams)),
	NewConfigCheck("solution.deploy.applyall", "", Static(DefaultSolutionDeployApplyAll)),
}

// FileChecks defines the checks that are done for files on the system.
var FileChecks = []Check{
	NewResourceFileCheck(OneExistsOnPath("podman","docker"), "%s was not found on your path", UpdateConfig("podman.path")),
	NewReqConfigDirCheck("build_home"),
	NewReqConfigDirCheck("save"),
	NewFixableConfigFileCheck("cli-config.yaml", EmptyFileCreator),
	NewFixableConfigFileCheck(filepath.Join("build_home", "casc.yaml"), TemplatedFileCreator(CascTemplateString)),
}
