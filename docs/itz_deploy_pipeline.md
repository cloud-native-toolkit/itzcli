## itz deploy pipeline

Deploys the given pipeline to the specified cluster

### Synopsis


Deploys the given pipeline to the cluster specified by --cluster-api-url ("-a").
The pipeline is identified by a UUID and can be found by executing the command:

    itz list pipelines

To view the current pipelines. With the pipeline ID, you can deploy the pipeline
to a cluster with the given API endpoint ("--cluster-api-url" or "-a"), and a 
username/password of a user with permissions to create Pipelines and PipelineRuns.

Example:

    itz deploy pipeline -p c567d9bd-5f0f-4254-bce1-c40ef1fedc0c \
      -a http://cluster.api.example.com \
      -u clusteruser \
      -P mysecretpassword



```
itz deploy pipeline [flags]
```

### Options

```
  -d, --accept-defaults           Accept defaults for pipeline parameters without asking
  -a, --cluster-api-url string    The URL of the target cluster
  -P, --cluster-password string   A password to login to the target cluster
  -u, --cluster-username string   A username to login to the target cluster
  -h, --help                      help for pipeline
  -p, --pipeline-id string        ID of the build in the catalog
  -c, --use-container             If true, the commands run in a container
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.itz/cli-config.yaml)
  -X, --debug           Prints trace messaging for debugging
      --json            Changes output to JSON
  -v, --verbose         Prints verbose messages
```

### SEE ALSO

* [itz deploy](itz_deploy.md)	 - Deploys a build in a cluster

