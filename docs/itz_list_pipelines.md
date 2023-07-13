## itz list pipelines

Displays a list of the available pipelines from the TechZone catalog.

### Synopsis


Displays a list of the available IBM Technology Zone (TechZone) pipelines from
the catalog.

From the TechZone catalog (see https://catalog.techzone.ibm.com/), a pipline is
a deployable component. It must be of kind "Component" and type "pipeline" to be
deployed to a cluster.

Example:

    itz list pipelines


```
itz list pipelines [flags]
```

### Options

```
  -c, --created         If true, limits the pipelines to my (created) pipelines
  -h, --help            help for pipelines
  -n, --name string     The name of the pipeline
  -o, --owner strings   The owner of the pipeline
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.itz/cli-config.yaml)
  -X, --debug           Prints trace messaging for debugging
      --json            Changes output to JSON
  -v, --verbose         Prints verbose messages
```

### SEE ALSO

* [itz list](itz_list.md)	 - Lists the summaries of the requested objects

