# IBM Technology Zone (ITZ) Command Line Quickstart

1. List basic usage using `--help` with any of the commands.

   ```
   $ itz --help
   IBM Technology Zone (ITZ) Command Line Interface (CLI)

   Usage:
   itz [command]

   Available Commands:
   auth        Manage tokens and authentication to APIs.
   completion  Generate the autocompletion script for the specified shell
   configure   Configures the itz command
   help        Help about any command
   reservation List and get IBM Technology Zone reservations.
   solution    Lists metadata, builds, and deploys solutions

   Flags:
   --config string   config file (default is $HOME/.itz/cli-config.yaml)
   -X, --debug           Prints trace messaging for debugging
   -h, --help            help for itz
   -v, --verbose         Prints verbose messages

   Use "itz [command] --help" for more information about a command.   
   ```
   
1. When you run the CLI for the first time, the CLI will create the `~/.itz`
folder. At any point during using the CLI, you can use the `itz doctor` command
to print information about the required programs (e.g., `podman`) and configuration
variables. An example of the `itz doctor` output is shown here:

   ```
   INFO[0000] Performing 19 checks...
   INFO[0000] podman...  OK
   INFO[0000] build_home...  OK
   INFO[0000] cli-config.yaml...  OK
   INFO[0000] build_home/casc.yaml...  OK
   INFO[0000] bifrost.api.image... OK
   INFO[0000] bifrost.api.local... OK
   INFO[0000] bifrost.api.url... OK
   INFO[0000] builder.api.refresh_token... OK
   INFO[0000] builder.api.url... OK
   INFO[0000] builder.api.username... OK
   INFO[0000] ci.api.image... OK
   INFO[0000] ci.api.local... OK
   INFO[0000] ci.api.password... OK
   INFO[0000] ci.api.url... OK
   INFO[0000] ci.api.user... OK
   INFO[0000] ci.buildtoken... OK
   INFO[0000] ci.localdir... OK
   INFO[0000] reservations.api.token... OK
   INFO[0000] reservations.api.url... OK
   INFO[0000] Done.
   ```

1. If you are running `itz` for the first time or need to fix missing configuration
values, you can try using the `--auto-fix` option. The `itz doctor --auto-fix`
command will do its best to default certain values, such as your local IP address,
to reasonable values but the `~/.itz/cli-config.yaml` may need some tweaking.

1. After the first run, you may need to use the `itz auth login` command to 
authenticate against the IBM Technology Zone APIs.

   ```
   itz auth login
   ```

1. Now that you have authenticated, you can list your current IBM Technology Zone 
reservations:

    ```
   $ itz reservation list
   - Redhat 8.5 Base Image with RDP (Fyre-2) (request id: 857b2bf8-cca8-4910-8fda-261229f84e90)
   ```

1. List the available solutions from the IBM Technology Zone catalog:

   ```
   $ itz solutions list
   - Composite Solution with IBM Maximo (id: 567514d3-ca27-4975-aa5b-d0450f9e779d)
   - TurboDemo (id: 8fc2e31d-bb6f-4534-8644-06c2a717ab5e)
   - Data Fabric for AWS, Azure and IBM Cloud (id: automation-datafabric)
   - Data Foundation for AWS, Azure and IBM Cloud (id: automation-datafoundation)
   - IBM Cloud z/OS Development Reference Architecture (id: automation-zos-dev)
   - AWS Quick Start OCP ROSA (id: aws-quickstart)
   - Azure Quick Start OCP IPI (id: azure-quickstart)
   - IBM Cloud for Financial Services with OpenShift (id: fs-cloud-szr-ocp)
   - IBM Cloud common Infrastructure Reference Architectures (id: ibmcloud-infrastructure)
   - IBM Cloud Quick Start OCP ROKS (id: ibmcloud-quickstart)
   - Integration Platform for AWS, Azure and IBM Cloud (id: integration-multicloud)
   - Maximo Application Suite for AWS, Azure and IBM Cloud (id: maximo-multicloud)
   - Turbonomic for AWS, Azure and IBM Cloud (id: turbonomic-multicloud)   
   ```
   
1. List the configuration from `ocpnow` to 

   ```
   $ itz configure list
   Project "my-project"

   Clusters:

   Name: some-cluster (deployed)
   URL: https://api.dlsdemo092622.activation-assets.com:6443

   Name: some-cluster-2 (deployed)
   URL: https://c111-e.us-east.containers.cloud.ibm.com:31012
   ```

   > *Note: It may be necessary to import your configuration from ocpnow by using
   > the `itz configure import --from-ocpnow-project /path/to/project1.yaml`
   > command.*
   
   > **Important: Direct integration with ocpnow is not complete but on the
   > roadmap.**

1. Select a solution to deploy from the list and deploy it at a customer site:

   ```
   $ itz solution deploy --solution automation-module-integration --cluster-name some-cluster-2
   ```

1. Alternatively, deploy the same solution in IBM Technology Zone using the web site.
