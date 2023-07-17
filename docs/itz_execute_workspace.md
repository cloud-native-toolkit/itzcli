## itz execute workspace

Executes the given workspace

### Synopsis


Executes the given workspace specified by the first arg. A "workspace" is a
containerized environment that can be used to run commmands without having to
install all the prerequisites. An example workspace that is provided by default
is the OCP (OpenShift Container Platform) Installer ("ocp-installer") workspace,
which can be used to install OCP in airgapped environments and on different
cloud environments such as AWS (Amazon Web Services) and Azure.

Using workspaces requires either Podman (see
https://podman.io/docs/installation) or Docker (see
https://docs.docker.com/engine/install/). During first-run of the CLI, the path
to either of these is configured automatically in the ~/.itz/cli-config.yaml
configuration file. Podman is preferred, so if you have the podman binary
installed, your configuration file should look like this:

    podman:
	    path: /usr/local/bin/podman

where the "path" is the full path to the binary, provided it is found on your
system. If podman is not installed, this will be set to the full path to your
docker binary (if installed).

The workspace itself is configured in the same file (~/.itz/cli-config.yaml) as
shown here:

    execute:
        workspace:
			ocpinstaller:
                image: quay.io/ibmtz/ocpinstaller:stable
                local: true
                name: ocp-installer
                type: interactive
                volumes:
                    - /Users/myuser/.itz/save:/usr/src/ocpnow/save

When you execute the "itz execute workspace ocpinstaller" command, the CLI looks
up the image information in the configuration file at the configuration key
"execute.workspace.[name]" where [name] is the value supplied on the command
line. For example:

   itz execute workspace ocpinstaller

Will execute the workspace shown in the above configuration.

While not officially supported, you can configure your own workspaces.


```
itz execute workspace [flags]
```

### Options

```
  -h, --help   help for workspace
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.itz/cli-config.yaml)
  -X, --debug           Prints trace messaging for debugging
      --json            Changes output to JSON
  -v, --verbose         Prints verbose messages
```

### SEE ALSO

* [itz execute](itz_execute.md)	 - Executes workspaces or pipelines in a cluster

