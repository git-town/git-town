# git town rename

> _git town rename [--force] [old-name] &lt;new-name&gt;_

The _rename_ command changes the name of the current branch in the local and
origin repository. It requires the branch to be in sync with its tracking branch
to avoid data loss. It also updates the proposals for the branch being renamed,
as well as proposals of its child branches into the branch being renamed.

Please be aware that most code hosting platforms are unable to update the head
branch (aka source branch) of proposals. If you rename a branch that already has
a proposal, the existing proposal will most likely end up closed and you have to
create a new proposal that supersedes the old one. If that happens, Git Town
will notify you. Updating proposals of child branches usually works.

### Positional arguments

When called with only one argument, the _rename_ command renames the current
branch to the given name.

When called with two arguments, it renames the branch with the given name to the
given name.

### --dry-run

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

### --force / -f

Renaming perennial branches requires confirmation with the `--force` aka `-f`
flag.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
