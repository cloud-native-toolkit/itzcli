# opcnow-go
Go implementation of opcnow installer

# The Question DSL

rootQuestion := prompt.NewPromptBuilder()


base.AddQuestion("What cloud provider(s) are you using?")
    .WithOption(Option("AWS"))
    .WithOption(Option("Azure"))
    .WithOption(Option("GPC"))

Prompter.ask(rootQuestion)