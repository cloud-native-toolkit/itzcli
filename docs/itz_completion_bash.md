## itz completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(itz completion bash)

To load completions for every new session, execute once:

#### Linux:

	itz completion bash > /etc/bash_completion.d/itz

#### macOS:

	itz completion bash > $(brew --prefix)/etc/bash_completion.d/itz

You will need to start a new shell for this setup to take effect.


```
itz completion bash
```

### Options

```
  -h, --help              help for bash
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

