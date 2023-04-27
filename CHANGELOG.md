# itz CHANGELOG 

## v0.1.24
* Added support for synchronizing podman machine date with host date.
* Added IBM SSO login with `itz auth login` command.

## v0.1.23
* Added support for `itz reservation get --reservation-id <reservation>` command
* Added support for getting solutions from a Backstage catalog.
* Updated automation for build, test, and release to GitHub actions.
* Updated bug in release tagging that accidentally used v1.0.0 versioning instead of v0.1.0 versioning.

## v0.1.21, v0.1.22
* Completely rewrote deployment backend to use container-based workflow instead
of container-based services (daemons) (https://github.ibm.com/skol/itzcli/pull/42).
* Changed `itz solution list` to default to listing public solutions (https://github.ibm.com/skol/itzcli/pull/45).
* Added `itz solution list -c` flag to list created (`-c`) solutions.
* Changed `itz reservation list` to include reservations in _Pending_ and _Scheduled_
states; changed `itz reservation list --all` to included _Deleted_ (https://github.ibm.com/skol/itzcli/pull/44).
* Added functionality to use refresh tokens for builder API after first installation.
* Updated podman path in configuration (see https://github.ibm.com/skol/itzcli/pull/37)
* Added binary and ZIP for Windows verison of the CLI.

## v0.1.20
* Added default configuration for ocp installer workspace (see https://github.ibm.com/skol/itzcli/pull/35)
* Improved configuration for directory mapping (see https://github.ibm.com/skol/itzcli/pull/34)

## v0.1.19
* Added functionality to start workspaces (see https://github.ibm.com/skol/itzcli/pull/33)

## v0.1.18
* Added this CHANGELOG into the release and project.
* Added wait for background service (daemon) to start
* Added check for failure to start with ":Z" mount option and correct itself in
event of failure. Mount option can be configured in the `~/.itz/cli-config.yaml`
file and also by supplying the `ITZ_CI_MOUNTOPTS` environment variable.
* Verified functionality of the `ITZ_PODMAN_PATH` environment variable and 
`podman.path` configuration value to set podman or docker location.
* Fixed problem with wrong IP address being used by default in --auto-fix when
the system is using a remote connection.
* Changed some configure commands to new `itz cluster <command>` where they
made sense, such as `itz cluster list` and `itz cluster import` commands.
