# IBM Technology Zone Command Line Interface (itz)

![build status](https://github.com/cloud-native-toolkit/itzcli/actions/workflows/build-go.yml/badge.svg) ![release status](https://github.com/cloud-native-toolkit/itzcli/actions/workflows/release-cli.yml/badge.svg)

The `itz` command line interface is a command line interface that provides CLI access to IBM Technology Zone.

## Introduction

Using `itz`, you can:

* List your existing reservations and get their status.
* List the available products that you can install in TechZone.
* Install or deploy products outside of TechZone using infrastructure as code.

## Quickstart

See the [QUICKSTART](QUICKSTART.md).

For usage documentation, see the documentation [in the docs folder](docs/itz.md).

## Architecture

For an architecture overview, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Installing `itz`

Release packages for your OS can be found at https://github.com/cloud-native-toolkit/itzcli/releases.

### Installing on Mac

> **_Note: if you have version 1.24 and installed itz with `brew`, you must
> use brew to uninstall itz and then re-install it._**

#### If you have itz already installed

If you have itz already installed and `itz version` outputs _1.24_, you must 
follow these steps first:

1. Use brew to uninstall the existing itz.
    ```bash
   $ brew uninstall itz
   ```
1. Untap the existing repository.
    ```bash
   $ brew untap cloud-native-toolkit/techzone
   ```
   
Once you have uninstalled itz, you can proceed to 
"[Installing itz using brew](#installing-itz-using-brew)".

#### Installing itz using brew

To install `itz` using [Homebrew](), follow these steps:

1. Tap the cask.
   ```bash
   $ brew tap cloud-native-toolkit/homebrew-techzone
   ```
2. Install ITZ with brew.
   ```bash
   $ brew install itz
   ```

### Signing on to IBM Technology Zone

Version v0.1.245 and higher of `itz` supports IBM's Single Sign On (SSO) to
authenticate against the TechZone APIs. To log in, type the following:

```bash
$ itz auth login --sso
```

This command will automatically open a browser. You can log into IBM
Verify with your IBM ID. When you are done, you can close the browser
window.
