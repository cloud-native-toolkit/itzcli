## itz completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	itz completion fish | source

To load completions for every new session, execute once:

	itz completion fish > ~/.config/fish/completions/itz.fish

You will need to start a new shell for this setup to take effect.


```
itz completion fish [flags]
```

### Options

```
  -h, --help              help for fish
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

