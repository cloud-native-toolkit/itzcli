## itz list reservations

Displays a list of your current reservations.

### Synopsis


Displays a list of your current IBM Technology Zone reservations.

By default, the CLI limits the reservations listed to those in "Pending",
"Provisioning", or "Ready" status. To view reservations in "Deleted" or
"Expired" status, use --all ("-a") to list all of your reservations.

The default output format is text in a tabular list. For scripting or
programmatic interaction, specify the --json flag to the command to view the
output in JSON format.

Examples:

    itz list reservations --json
    itz list reservations --all


```
itz list reservations [flags]
```

### Options

```
  -a, --all    If true, list all reservations (including expired)
  -h, --help   help for reservations
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

