# Using the Activation ToolKit (ATK) Installer

# Working with the modular framework

There is a .now directory in each git project that contains manifest information
about the project. The information contains, but is not limited to:

* The variables required by the project, such as input variables for Terraform,
Ansible, etc.  
* How to build the project.
* How to deploy the project.

This file should be maintained by the module maintainer, as they would know best
how each one of these should work.

## Example `.now/manifest.yml`

Here is an example of the `.now/manifest.yml` file, which contains the information
about how to query for information, prepare the project, build, and deploy the
project.

```yaml
# This is the name of the project
name: My Base Project
# The version of this file spec
version: 0.1
# The URL to the .git repo that is the template for this project.
template_url: https://github.com/someorg/someproject

facets:
  - terraform
  - travisci

dependencies:
  - None

spec:

  params:

    list:
      cmd: docker run hello-world

    validate:
      cmd: docker run hello-world

  prepare:

    commands:
      - docker run hello-world

  build:

    commands:
      - docker run hello-world

  deploy:

    commands:
      - docker run hello-world
```

## Using the CLI as an API to the file

The main concept behind the combination of the CLI and the `manifest.yml`
file specification is that each maintainer of a project who
wishes to expose their installer/product in a GitOps project can do so in 
a modular, consistent way.

## The parameter spec

One of the features of the metadata and CLI is to obtain required parameters
for a specific project. This allows wizards, UIs, or other mechanisms to build
intelligent prompts dynamically. Since the parameters metadata is maintained by
the project maintainers, it does not need to be maintained in multiple places.

The `params/list/cmd` key is used to declare a command that is suitable for
providing the output--in JSON format, by default--of the variables used in
the project.

```yaml
vars:
  - name: cloud-provider
    order: 1
    prompt: Which cloud provider(s) would you like to use?
    type: multi-choice
    options:
      - AWS
      - Azure
      - IBM Cloud
      - VMWare
    default: IBM Cloud
    required: yes

  - name: vpc
    order: 2
    prompt: Would you like to use an existing VPC or create a new one?
    type: choice
    options:
      - Existing
      - New
    default: New
    required: yes

  - name: vpc-ip
    order: 2.1
    prompt: What is the IP address range of the VPC?
    type: text
    default: 10.10.0.0/16
    validationExpr: ^[0-9./]+$
    askWhen: vpc == New

  - name: subnet-ip
    order: 2.2
    prompt: What is the IP address range of the subnet?
    type: text
    default: 10.10.10.0/24
    validationExpr: ^[0-9./]+$
    askWhen: vpc == New
```

* Ability to call an API and get a number of packages
* To make a package available, just put a .now folder in it and have the file
correctly formatted--in other words, this process "dog foods" by using GitOps
for itself.
* The package can be downloaded or pulled from a git repository and stored somewhere
locally while the configuration is being collected.


