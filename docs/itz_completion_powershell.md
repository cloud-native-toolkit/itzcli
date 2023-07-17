## itz completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	itz completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
itz completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.itz/cli-config.yaml)
  -X, --debug           Prints trace messaging for debugging
      --json            Changes output to JSON
  -v, --verbose         Prints verbose messages
```

### SEE ALSO

* [itz completion](itz_completion.md)	 - Generate the autocompletion script for the specified shell

