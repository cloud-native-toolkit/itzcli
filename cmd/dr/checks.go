package dr

// Create the checks for the configuration values that I know I'll need.

// AllConfigChecks checks for configuration values on the system and defines
// defaulters for fixing the missing values if the user specifies --auto-fix.
var AllConfigChecks = []Check{
	// The bifrost configuration values
	NewConfigCheck("bifrost.api.image", "", Static("localhost/bifrost:latest")),
	NewConfigCheck("bifrost.api.local", "", Static(true)),
	NewConfigCheck("bifrost.api.url", "", ServiceURL("http", 8088)),
	// The builder configuration values
	NewConfigCheck("builder.api.token", "", Prompter("There is no token defined for the Builder API. Please provide one:")),
	NewConfigCheck("builder.api.url", "", Static("https://ascent-bff-mapper-staging.dev-mapper-ocp-4be51d31de4db1e93b5ed298cdf6bb62-0000.eu-de.containers.appdomain.cloud")),
	NewConfigCheck("bifrost.api.username", "", Prompter("Please enter your ibm.com email address (for getting your solutions):")),
	// The Jenkins (ci) configuration values
	NewConfigCheck("ci.api.image", "", Static("localhost/atkci:latest")),
	NewConfigCheck("ci.api.local", "", Static(true)),
	NewConfigCheck("ci.api.password", "", RandomVal(24)),
	NewConfigCheck("ci.api.url", "", ServiceURL("http", 8080)),
	NewConfigCheck("ci.api.user", "", Static("bifrost")),
	NewConfigCheck("ci.buildtoken", "", Static("68b2a8e9ffe5a395d839d9bf87db6800")),
	NewConfigCheck("ci.localdir", "", ConfigDir("build_home")),
	// The reservations configuration values
	NewConfigCheck("reservations.api.token", "", Prompter("There is no token defined for the Reservations API. Please provide one:")),
	NewConfigCheck("reservations.api.url", "", Static("https://api.techzone.ibm.com/api/my/reservations/all")),
}

// FileChecks defines the checks that are done for files on the system.
var FileChecks = []Check{
	NewBinaryFileCheck("podman", "%s was not found on your path"),
}
