import { deepEqual } from "node:assert/strict"
import { suite, test } from "node:test"
import { HelpOutput } from "./help-output.ts"

suite("HelpOutput", () => {
  suite(".flags", () => {
    test("append command", () => {
      const output = new HelpOutput(`
Create a new feature branch as a child of the current branch.

Consider this stack:

main
 \
* feature-1

We are on the "feature-1" branch,
which is a child of branch "main".
After running "git town append feature-2",
the repository will have these branches:

main
 \
  feature-1
   \
*   feature-2

The new branch "feature-2"
is a child of "feature-1".

If there are no uncommitted changes,
it also syncs all affected branches.

Usage:
  git-town append <branch> [flags]

Flags:
      --auto-resolve     auto-resolve phantom merge conflicts
  -b, --beam             beam some commits from this branch to the new branch
  -c, --commit           commit the stashed changes into the new branch
  -d, --detached         don't update the perennial root branch
      --dry-run          print but do not run the Git commands
  -h, --help             help for append
  -m, --message string   the commit message
      --propose          propose the new branch
  -p, --prototype        create a prototype branch
      --push             push local branches
      --stash            stash uncommitted changes when creating branches
      --sync             sync branches (default true)
  -v, --verbose          display all Git commands run under the hood
`)
      const have = output.flags()
      const want = [
        ["--auto-resolve"],
        ["-b", "--beam"],
        ["-c", "--commit"],
        ["-d", "--detached"],
        ["--dry-run"],
        ["-h", "--help"],
        ["-m", "--message string"],
        ["--propose"],
        ["-p", "--prototype"],
        ["--push"],
        ["--stash"],
        ["--sync"],
        ["-v", "--verbose"],
      ]
      deepEqual(have, want)
    })

    test("branch command", () => {
      const output = new HelpOutput(`
Display the local branch hierarchy and types.

Git Town's equivalent of the "git branch" command.

Usage:
  git-town branch [flags]

Flags:
  -d, --display-types string[="all"]   display the branch types
  -h, --help                           help for branch
  -o, --order string                   sort order for branch list (asc or desc)
  -v, --verbose                        display all Git commands run under the hood
`)
      const have = output.flags()
      const want = [
        ["-d", "--display-types string"],
        ["-h", "--help"],
        ["-o", "--order string"],
        ["-v", "--verbose"],
      ]
      deepEqual(have, want)
    })
  })

  suite(".flagLine", () => {
    const tests = [
      {
        give: "  -b, --beam             description",
        want: [["-b", "--beam"]],
      },
      {
        give: `  -d, --display-types string[="all"]   display the branch types`,
        want: [["-d", "--display-types string"]],
      },
    ]
    for (const { give, want } of tests) {
      test(give, () => {
        const output = new HelpOutput("")
        const have = output.flagLine(give)
        deepEqual(have, want)
      })
    }
  })
})
