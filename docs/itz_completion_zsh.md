## itz completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(itz completion zsh)

To load completions for every new session, execute once:

#### Linux:

	itz completion zsh > "${fpath[1]}/_itz"

#### macOS:

	itz completion zsh > $(brew --prefix)/share/zsh/site-functions/_itz

You will need to start a new shell for this setup to take effect.


```
itz completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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

