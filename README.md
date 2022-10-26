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

## Searching for solutions

To list builder solutions, use the following command

```bash
$ atk solution list
```

This will list only _your_ solutions by default.

## Listing Your Reservations



