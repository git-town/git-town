<a textrun="command-documentation">
# Alias command

Adds or removes default global aliases

Global aliases allow Git Town commands to be used like native Git commands.
When aliases are set, you can run "git hack" instead of having to run "git town hack".
Example: "git append" becomes equivalent to "git town append".

When adding aliases, no existing aliases will be overwritten.

Note that this can conflict with other tools that also define additional Git commands.

#### Usage

```
git town alias (true | false)
```
</a>

#### DESCRIPTION

Adding the default git aliases removes the need for `town` for the following commands
(append, hack, kill, new-pull-request, prepend, prune-branches, rename-branch, repo, ship, and sync).
Example: `git append` becomes equivalent to `git town append`.

When adding aliases, no existing aliases will be overwritten.
