# ITZCLI Architecture

Before reading, see the documentation for
[atkmod](https://github.com/cloud-native-toolkit/atkmod), on which this CLI
heavily relies. It is primarily use for interfacing with the docker and podman
CLIs, reading YAML configuration files, and a simple state machine that handles
lifecycle.

Other main libraries are used are:

* [Cobra](https://github.com/spf13/cobra) - for CLI command scaffolding and
handling.
* [Viper](https://github.com/spf13/viper) - for flexible, Java Spring-like 
configuration.
* [logrus](https://github.com/sirupsen/logrus) - for structured and leveled 
logging.

## Background

From the very first version of this CLI, the goal was to create a lightweight of
a CLI as possible and defer much, if not most, of the logic to a longer-running
service that could run locally in the background or even remotely. The local
service was inspired a bit by the Gradle "daemon" for caching and handling
long-running tasks such as deployments (builds not being altogether different).
The background services were Docker-based, meaning they could run locally or
could be installed in OpenShift for remote execution in any infrastructure.

The first services started as a "mediator" service written using Java, Spring
Boot, and Apache Camel; and a build service that was Jenkins running in a
container and mounting a workspace in the `itz` home directory. This worked well
for some basic demonstrations and proofs of concepts, but weren't the long-term
goal and quickly became more difficult to maintain than temporary work should
become.

December 2022 saw a large refactoring of the ITZ CLI code, taking most of the Go
logic that had become part of the CLI code and moving toward the end-goal of the
CLI, which had always been to act as a lightweight orchestrator of containers
through an opinionated and small lifecycle for deployment.

The atkmod documentation contains more background, but the summary is that a
container-based "plugin" approach for stages in the lifecycle will allow
software modules to use different tools structures while allowing the core `itz`
code to stay (roughly) the same over time.

Currently, and for the short term, the configuration file for the CLI
(`~/.itz/cli-config.yml`) contains a single lifecycle and plugins made to
generically handle IBM Technology Zone (TechZone) builder solutions in a
predictable format. This is an intermediate solution until the complete
atkmod solution can be implemented.

See the architecture diagram ([(full size here](assets/itzcli-arch-overview.png)):

![ITZ CLI architecture overview](assets/itzcli-arch-overview.png)

## Basic philosophy and goals

The architectural goals for the ITZ CLI are:

1. To keep the codebase as lightweight and small as possible, using containers
images to do the "heavy lifting".
1. To have a predictable command structure that enables simple, yet powerful,
interaction with TechZone to deploy solutions in hybrid cloud infrastructure.
1. To also provide a web API (via `itz api start`) that surfaces the command
structure over HTTP for integration with UIs and bots.

## Use of cobra

[Cobra](https://github.com/spf13/cobra) is a framework for making command-line
interfaces. You can read more documentation also about the
[Cobra Generator](https://github.com/spf13/cobra-cli/blob/main/README.md), used
to generate the Go file scaffolding for the commands.

## Use of viper

[Viper](https://github.com/spf13/viper) is used for configuration and offers
quite a bit of flexibility, including loading YAML structures directly into
objects (see [this
example](https://github.com/cloud-native-toolkit/itzcli/blob/e6fda5cf93fc513f5a81ec3bb9a0473531843beb/pkg/executils.go#L190))
as well as providing the ability to override configuration using environment
variables. That means that any value that is configurable in the default
configuration file (`~/.itz/cli-config.yml`), such as `podman.path` can be
overridden using an environment variable, such as `ITZ_PODMAN_PATH`. See [the
Viper documentation](https://github.com/spf13/viper#putting-values-into-viper)
for more information about how to use Viper for configuration.

## Roadmap

For roadmap items, see [features in the Issues log](https://github.com/cloud-native-toolkit/itzcli/issues?q=is%3Aopen+is%3Aissue+label%3Afeature)
