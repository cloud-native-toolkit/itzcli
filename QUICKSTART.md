# ATK Command Line Quickstart

1. List basic usage using `--help` with any of the commands.

   ```
   ./atk --help
   Activation ToolKit (ATK) Command Line Interface (CLI)

   Usage:
   atk [command]

   Available Commands:
   auth        Manage tokens and authentication to APIs.
   completion  Generate the autocompletion script for the specified shell
   configure   Configures the atk command
   help        Help about any command
   reservation List and get TechZone reservations.
   solution    Lists metadata, builds, and deploys solutions

   Flags:
   --config string   config file (default is $HOME/.atk.yaml)
   -X, --debug           Prints trace messaging for debugging
   -h, --help            help for atk
   -v, --verbose         Prints verbose messages

   Use "atk [command] --help" for more information about a command.   
   ```
   
1. When you run the CLI for the first time, the CLI will create the `~/.atk`
folder. At any point during using the CLI, you can use the `atk doctor` command
to print information about the required programs (e.g., `podman`) and configuration
variables. An example of the `atk doctor` output is shown here:

   ```
   INFO[0000] Performing 19 checks...
   INFO[0000] podman...  OK
   INFO[0000] build_home...  OK
   INFO[0000] cli-config.yaml...  OK
   INFO[0000] build_home/casc.yaml...  OK
   INFO[0000] bifrost.api.image... OK
   INFO[0000] bifrost.api.local... OK
   INFO[0000] bifrost.api.url... OK
   INFO[0000] builder.api.token... OK
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

1. If you are running `atk` for the first time or need to fix missing configuration
values, you can try using the `--auto-fix` option. The `atk doctor --auto-fix`
command will do its best to default certain values, such as your local IP address,
to reasonable values but the `~/.atk/cli-config.yaml` may need some tweaking.

1. After the first run, you may need to use the `atk auth login` command to 
save your API keys so that ATK can access the solutions and reservations APIs.
See "Getting API keys" for more information about obtaining API keys. Once you
have them, use the following commands to save them in ATK's configuration:

   ```
   ./atk auth login --from-file /tmp/token.txt --service-name builder
   ./atk auth login --from-file /tmp/token.txt --service-name reservations
   ```

2. Now that you have saved the API keys, you can list your current TechZone 
reservations:

    ```
   ./atk reservation list
   - Redhat 8.5 Base Image with RDP (Fyre-2) (request id: 857b2bf8-cca8-4910-8fda-261229f84e90)
   ```

3. List your solutions from the **IBM Technology Zone Accelerator Toolkit** > **Solutions** site:

    ```
   ./atk solutions list
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
   
   > Note: it may be necessary to update your token for the Solution Builder API,
   > which you can do with the command `./atk auth login --from-file /tmp/token.txt --service-name builder `,
   > after you have visited the Solution Builder website and saved the API token into a file
   > called `/tmp/token.txt`

4. List the configuration from `ocpnow` to 

   ```
   ./atk configure list
   Project "my-project"

   Clusters:

   Name: some-cluster (deployed)
   URL: https://api.dlsdemo092622.activation-assets.com:6443

   Name: some-cluster-2 (deployed)
   URL: https://c111-e.us-east.containers.cloud.ibm.com:31012
   ```

   > *Note: It may be necessary to import your configuration from ocpnow by using
   > the `./atk configure import --from-ocpnow-project /path/to/project1.yaml`
   > command.*
   
   > **Important: Direct integration with ocpnow is not complete but on the
   > roadmap.**

5. Select a solution to deploy from the list and deploy it at a customer site:

   ```
   ./atk solution deploy --solution automation-module-integration --cluster-name some-cluster-2
   ```

6. Alternatively, deploy the same solution in TechZone using the web site.
