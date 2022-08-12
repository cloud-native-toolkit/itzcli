# Activation ToolKit (ATK) Command Line Interface (CLI)

## Listing available Cloud Paks

1. To list the available Cloud Paks, use the command `atk cloud-pak list`.

Go implementation of Activation ToolKit (ATK) installer

# The Question DSL

There is a DSL (Domain-Specific Language) for building questions in code. Once
built, the questions can be handed off to a Prompter, which is a kind of 
[Controller](), that will iterate through the questions and display them as
prompts to the user. The answers can then be persisted to a writer.

## Example

```go
rootQuestion := prompt.NewPromptBuilder()

base.AddQuestion("What cloud provider(s) are you using?")
    .WithOption(Option("AWS"))
    .WithOption(Option("Azure"))
    .WithOption(Option("GPC"))

answers := Prompter.ask(rootQuestion)

cfgFile := writers.NewConfigFileWriter()
cfgFile.Write(answers)

```
