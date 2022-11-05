# itz CHANGELOG 

## vNext
* Added this CHANGELOG into the release and project.
* Added wait for background service (daemon) to start
* Added check for failure to start with ":Z" mount option and correct itself in
event of failure. Mount option can be configured in the `~/.itz/cli-config.yaml`
file and also by supplying the `ITZ_CI_MOUNTOPTS` environment variable.
* Verified functionality of the `ITZ_PODMAN_PATH` environment variable and 
`podman.path` configuration value to set podman or docker location.