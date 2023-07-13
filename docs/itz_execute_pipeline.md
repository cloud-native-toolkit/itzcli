## itz execute pipeline

Executes the given pipeline

### Synopsis


Executes the given pipeline provided by the --pipeline-url ("p") and
--pipeline-run-url ("r") arguments on a Kubernetes or OpenShift cluster.  The
cluster is identified with the --cluster-api-url ("c") argument. You must also
supply the --cluster-username and --cluster-password arguments, with the a user
and password, respectively, with sufficient privileges to execute the pipeline.

The command will read the parameters from the pipeline. If there are default
values specified in the pipeline, you can accept all of them by using the
--accept-defaults ("d") argument. By accepting defaults, the CLI will only
provide prompts for the parameters without default values specified in the
pipeline parameters.

For non-interactive execution, for scripting or automation, you can provide the
values to parameters two different ways. First, you can supply the parameter
values as environment variables that begin with ITZ_ and then the rest of the
variable in uppercase, with non-number and non-digits replaced by _. For
example, if a variable is called "repo-url", the environment variable is
"ITZ_REPO_URL".

    ITZ_REPO_URL=http://github.com/me/myrepo itz execute pipeline \
      --pipeline-url file://somepipeline.yaml \
	  --pipeline-run-url file://somepipelinerun.yaml \
	  --cluster-api-url http://localhost \
	  --cluster-username myclusteruser \
	  --cluster-password mysecretpassword 

You can also provide the parameters as arguments at the end of the command line.
For example, for the repo-url variable, you could execute the following command:

    itz execute pipeline --pipeline-url file://somepipeline.yaml \
	  --pipeline-run-url file://somepipelinerun.yaml \
	  --cluster-api-url http://localhost \
	  --cluster-username myclusteruser \
	  --cluster-password mysecretpassword \
	  "repo-url=http://github.com/me/myrepo"
	

```
itz execute pipeline [flags]
```

### Options

```
  -d, --accept-defaults           Accept defaults for pipeline parameters without asking
  -a, --cluster-api-url string    The URL of the target cluster
  -P, --cluster-password string   A password to login to the target cluster
  -u, --cluster-username string   A username to login to the target cluster
  -h, --help                      help for pipeline
  -r, --pipeline-run-url string   The URL of the pipeline run as YAML
  -p, --pipeline-url string       The URL of the pipeline as YAML
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

* [itz execute](itz_execute.md)	 - Executes workspaces or pipelines in a cluster

