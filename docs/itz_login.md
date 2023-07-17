## itz login

Uses your browser to authenticate with TechZone.

### Synopsis


Opens a browser window for you to authenticate with IBM Technology Zone using
your IBMid. 

Upon successful login, the CLI updates the configuration with an authentication
token that will be used to access the IBM Technology Zone API as well as the 
IBM Technology Zone Catalog API.


```
itz login [flags]
```

### Options

```
  -f, --from-file string   The name of the file that contains the token.
  -h, --help               help for login
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.itz/cli-config.yaml)
  -X, --debug           Prints trace messaging for debugging
      --json            Changes output to JSON
  -v, --verbose         Prints verbose messages
```

### SEE ALSO

* [itz](itz.md)	 - IBM Technology Zone (ITZ) Command Line Interface (CLI), version No Version Provided

