# IBM TechZone Activation ToolKit Command Line Interface (atk)

The `atk` command line interface is a scriptable command line interface that provides CLI access to IBM TechZone.

## Introduction

Using `atk`, you can:

* List your existing reservations and get their status.
* Log in to your reservations by printing credentials and/or login links.
* Create new GitOps projects on your own infrastructure using your favorite templates and reference architectures in TechZone.
* Install or provision your local projects created from TechZone.
* Submit your own content for use in TechZone.

## Installing `atk`

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

## Searching for content

You can use the `atk` command line interface to search for TechZone content. Not all content types available at
[https://techzone.ibm.com](https://techzone.ibm.com) are avaiable for use by `atk`, however. The main purpose of `atk` 
is to easily create your own projects based on what you saw on TechZone, so content available to `atk` is limited to
_Environments_, _Activation Kits_, and _Experiences_.

## Listing Your Reservations



