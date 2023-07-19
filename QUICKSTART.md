# IBM Technology Zone (ITZ) Command Line Quickstart

1. On Mac and Linux operating systems, you can install `itz` by using the 
shell script:
    ```bash
    curl https://raw.githubusercontent.com/cloud-native-toolkit/itzcli/main/scripts/install.sh | bash -
    ```
    Or, on Mac you can also use `brew` to install `itz`:
    ```bash
   brew update
   brew tap cloud-native-toolkit/homebrew-techzone
   brew install itz
    ```

1. List basic usage using `--help` with any of the commands.

    ```
    $ itz --help
    IBM Technology Zone (ITZ) Command Line Interface (CLI)
    
    Usage:
      itz [command]
    
    Available Commands:
      completion  Generate the autocompletion script for the specified shell
      deploy      Deploys a build in a cluster
      doctor      Checks the environment and configuration
      execute     Executes workspaces or pipelines in a cluster
      help        Help about any command
      list        Lists the summaries of the requested objects
      login       Uses your browser to authenticate with TechZone.
      show        Shows the details of the requested single object
      version     Prints the current version and exits
    
    Flags:
          --config string   config file (default is $HOME/.itz/cli-config.yaml)
      -X, --debug           Prints trace messaging for debugging
      -h, --help            help for itz
          --json            Changes output to JSON
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

1. After the first run, you may need to use the `itz login` command to 
authenticate against the IBM Technology Zone APIs. This command will open a browser
through which you can log into http://techzone.ibm.com using your IBM id.

    ```
    itz login
    ```
    You can also log in without using the browser, which is useful if you are trying
    to run `itz` on a headless VM or script. To log in without opening a browser, use
    the `--from-file` flag to load the API token from a file:
    ```bash
    $ echo "thisismyapitokenigotfrommyechzoneprofile" > /tmp/token.txt
    $ itz login --from-file /tmp/token.txt
    ```

1. Now that you have authenticated, you can list your current IBM Technology Zone 
reservations:

    ```
    $ itz list reservations
    ```

1. List the available solutions from the IBM Technology Zone catalog:

   ```
   $ itz list pipelines
    NAME                                                                ID                                    NAMESPACE                     
    Deployer Pipeline for Maximo Application Suite Automation Operator  09f581a8-c4f1-47e0-b66f-41e4795c1ad5  default                       
    Deployer CP4BA Starter 22.x                                         5d473a15-35f0-4346-8826-b96714b3dff3  default                       
    Deployer CP4I 2022.4 Platform UI pipeline                           700bbd44-2d02-4c5f-96ed-44a3e0829862  default                       
    Deployer CP4S 1.10                                                  8a041e60-9e16-485f-8201-7ae3f6325c25  default                       
    Deployer CP4I 2022.4 for IBM API Connect                            9784a15a-3ddc-4aef-994e-fa1284341146  default                       
    Deployer IBM Turbonomic pipeline                                    a5285c9c-8ff7-460a-b63c-f9102145fa7a  default                       
    Deployer for IBM watsonx.data GA versions                           c52a12e0-6511-4de8-9af1-2108eb84266e  default                       
    Deployer CP4I 2022.4 for IBM App Connect                            c658d99d-fcd6-4ae3-8c3f-daf906650c45  default                       
    Deployer for IBM watsonx.data pre-release version                   cd29b8cb-7f82-4192-aebf-27e664bce248  default                       
    Deployer CP4I 2022.4 for IBM MQ                                     f0548ebf-cac7-4194-bd58-f00b5d623ec4  default                       
    Deployer CP4D Cloud Pak Deployer pipeline                           fd65d166-0995-408d-ba35-28dc90471aa4  default                       
   ```

1. Select a solution to deploy from the list and deploy it at a customer site:

    ```
    $ itz deploy pipeline --pipeline-id 8a041e60-9e16-485f-8201-7ae3f6325c25 \
      -a https://mycluster.example.com \
      -u myclusteruser \
      -P mysecretclusterpassword
    ```
    For more information about `itz deploy pipeline`, see the `--help` documentation or
    use `man itz-deploy-pipeline` to view more examples.
