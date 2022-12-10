# IBM Technology Zone Command Line Interface (itz)

version 0.1

[![Build Status](https://travis.ibm.com/skol/itzcli.svg?token=wGYsX6PCXyDddvgpBC56&branch=main)](https://travis.ibm.com/skol/itzcli)

The `itz` command line interface is a scriptable command line interface that provides CLI access to IBM Technology Zone.

## Introduction

Using `itz`, you can:

* List your existing reservations and get their status.
* Log in to your reservations by printing credentials and/or login links.
* Create new GitOps projects on your own infrastructure using your favorite templates and reference architectures in IBM
Technology Zone.
* Install or provision your local projects created from IBM Technology Zone Accelerator Toolkit.
* Submit your own content for use in IBM Technology Zone.

See the [QUICKSTART](QUICKSTART.md).

For usage documentation, see the documentation [in the docs folder](docs/itz.md).

## Installing `itz`

Release packages for your OS can be found at https://github.com/cloud-native-toolkit/itzcli/releases.

## Getting API keys

The `itz` command makes calls to two different APIs to get reservation and 
builder solution data. To make those calls, the `itz` command requires API 
keys for authorization. The keys are available when you log into
[IBM Technology Zone](https://techzone.ibm.com/my/profile) 
and [IBM Technology Zone Accelerator Toolkit](https://builder.cloudnativetoolkit.dev/).
See the next sections for more information.

### Getting the IBM Technology Zone API key

To get the API key for the listing your IBM Technology Zone reservations, follow
these steps:

1. Log into [IBM Technology Zone](https://techzone.ibm.com/).
1. Go to your profile by clicking your picture in the upper right corner or by
   going to [My Profile](https://techzone.ibm.com/my/profile).
1. Find the **API token** field and click the copy icon or copy the value into
your clipboard.
1. Once you have the key, put the key into a temporary file, such as `/tmp/token.txt`.
1. Use the `itz auth login` command with the `--from-file` option to import the
token into your configuration, like this:
    ```
    ./itz auth login --from-file /tmp/token.txt --service-name reservations
    ```

### Getting the IBM Technology Zone Accelerator Toolkit API key

To get the API key for the listing your solutions that you have created using
the IBM Technology Zone Accelerator Toolkit, follow these steps:

1. Log into [IBM Technology Zone Accelerator Toolkit](https://builder.cloudnativetoolkit.dev/).
1. Go to your profile by clicking your picture in the upper right corner.
1. Find the **API token** field and click the copy icon to copy the API key into
   your clipboard.
1. Once you have the key, put the key into a temporary file, such as `/tmp/token.txt`.
1. Use the `itz auth login` command with the `--from-file` option to import the
   token into your configuration, like this:
    ```
    ./itz auth login --from-file /tmp/token.txt --service-name builder
    ```