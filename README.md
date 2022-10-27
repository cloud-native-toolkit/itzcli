# IBM TechZone Activation ToolKit Command Line Interface (atk)

version 0.1

[![Build Status](https://travis.ibm.com/skol/atkcli.svg?token=wGYsX6PCXyDddvgpBC56&branch=main)](https://travis.ibm.com/skol/atkcli)

The `atk` command line interface is a scriptable command line interface that provides CLI access to IBM TechZone.

## Introduction

Using `atk`, you can:

* List your existing reservations and get their status.
* Log in to your reservations by printing credentials and/or login links.
* Create new GitOps projects on your own infrastructure using your favorite templates and reference architectures in TechZone.
* Install or provision your local projects created from TechZone.
* Submit your own content for use in TechZone.

See the [QUICKSTART](QUICKSTART.md).

For usage documentation, see the documentation [in the docs folder](docs/atk.md).

## Installing `atk`

> **Note: the "Installing" section currently describes desired functionality, but
> not functionality that is currently implemented. For now, to install atkcli on
> your machine, either clone this repo and build it locally or use the releases.
> The releases are currently only compiled for Linux x86.**

The `atk` command line interface can be installed on different operating systems using common packages for each system.

### MacOS

You can use `homebrew` to install the `atk` command line interface:

```bash
$ brew install atk
```

### Linux

`atk` is available on different Linux distributions as a RPM, APT package, or archive.

#### Red Hat, CentOS (rpm)

```bash
$ rpm -i https://static.techzone.ibm.com/packages/atk-0.1.0.rpm
```

#### Debian, Ubuntu (apt)

```bash
$ apt install atk
```

## Getting API keys

The `atk` command makes calls to two different APIs to get reservation and 
builder solution data. To make those calls, the `atk` command requires API 
keys for authorization. The keys are available when you log into
[IBM Technology Zone](https://techzone.ibm.com/my/profile) 
and [TechZone Accelerator Toolkit](https://builder.cloudnativetoolkit.dev/).
See the next sections for more information.

### Getting the IBM Technology Zone API key

To get the API key for the listing your IBM TechZone reservations, follow
these steps:

1. Log into [IBM Technology Zone](https://techzone.ibm.com/).
1. Go to your profile by clicking your picture in the upper right corner or by
   going to [My Profile](https://techzone.ibm.com/my/profile).
1. Find the **API token** field and click the copy icon or copy the value into
your clipboard.
1. Once you have the key, put the key into a temporary file, such as `/tmp/token.txt`.
1. Use the `atk auth login` command with the `--from-file` option to import the
token into your configuration, like this:
    ```
    ./atk auth login --from-file /tmp/token.txt --service-name reservations
    ```

### Getting the TechZone Accelerator Toolkit API key

To get the API key for the listing your solutions that you have created using
the IBM Technology Zone Accelerator Toolkit, follow these steps:

1. Log into [TechZone Accelerator Toolkit](https://builder.cloudnativetoolkit.dev/).
1. Go to your profile by clicking your picture in the upper right corner.
1. Find the **API token** field and click the copy icon to copy the API key into
   your clipboard.
1. Once you have the key, put the key into a temporary file, such as `/tmp/token.txt`.
1. Use the `atk auth login` command with the `--from-file` option to import the
   token into your configuration, like this:
    ```
    ./atk auth login --from-file /tmp/token.txt --service-name builder
    ```