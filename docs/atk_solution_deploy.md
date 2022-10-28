## atk solution deploy

Deploys the specified solution.

### Synopsis

Use this command to deploy the specified solution
locally in your own environment. You can specify the environment by using
either --cluster-name or --reservation as a target.

    --cluster-name requires the name of a cluster that has been deployed
using ocpnow. To see the clusters that are configured, use the "atk configure 
list" command to list the available clusters. If you have none, you may need to
import the ocpnow configuration using the "atk configure import" command. See
the help for those commands for more information.

    --reservation requires the id of a reservation in the TechZone system. Use
the "atk reservation list" command to list the available reservations.

```
atk solution deploy [flags]
```

### Options

```
  -c, --cluster-name string   The name of the cluster created by ocpnow to target.
  -f, --file string           The full path to the solution file to be deployed.
  -h, --help                  help for deploy
  -r, --reservation string    The id of the reservation to target.
  -s, --solution string       The name of the solution to be deployed.
  -u, --use-cache             If true, uses a cached solution file instead of downloading from target.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.atk.yaml)
  -X, --debug           Prints trace messaging for debugging
  -n, --name string     The name of the solution
  -v, --verbose         Prints verbose messages
```

### SEE ALSO

* [atk solution](atk_solution.md)	 - Lists metadata, builds, and deploys solutions

