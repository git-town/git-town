# Git Town Development

Welcome to the Git Town developer guide. Hacking on Git Town is easy and a great
way to learn to work on a medium-sized codebases. These guidelines will help you
get started hacking on Git Town.

- [work on the source code](development.md)
- [work on the website](website.md)
- [make a release](release.md)

<code textrun="verify-make-command">make test</code> shows all available make
tasks.

## Architecture

- [overview](architecture.md)
- [branch hierarchy](branch_hierarchy.md): how Git Town sees branches
- [repository drivers](drivers.md): third-party specific functionality
- [step lists](steps_list.md): how most of the Git Town commands implement
  automatic, bullet-proof `undo` and `continue` commands
- [test architecture](test-architecture.md)
