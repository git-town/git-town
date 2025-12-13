import { deepEqual } from "node:assert/strict"
import { suite, test } from "node:test"
import {
  findGroupWithPositiveFlag,
  FlagLine,
  getPositiveFlagName,
  HelpOutput,
  isNegatedFlagsGroup,
  matchesFlag,
  mergeFlags,
  replaceValueNotation,
} from "./help_output.ts"

const appendHelpOutput = `
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
      --auto-resolve      auto-resolve phantom merge conflicts
  -b, --beam              beam some commits from this branch to the new branch
  -c, --commit            commit the stashed changes into the new branch
  -d, --detached          don't update the perennial root branch
      --dry-run           print but do not run the Git commands
  -h, --help              help for append
  -m, --message string    the commit message
      --no-auto-resolve   don't auto-resolve
      --no-detached       disable detached
      --no-push           don't push branches
      --no-stash          don't stash uncommitted changes
      --no-sync           don't sync branches
      --propose           propose the new branch
  -p, --prototype         create a prototype branch
      --push              push local branches
      --stash             stash uncommitted changes when creating branches
      --sync              sync branches (default true)
  -v, --verbose           display all Git commands run under the hood
`

suite("HelpOutput", () => {
  suite(".flags()", () => {
    test("append command", () => {
      const output = new HelpOutput(appendHelpOutput)
      const have = output.flags()
      const want = [
        ["--auto-resolve", "--no-auto-resolve"],
        ["-b", "--beam"],
        ["-c", "--commit"],
        ["-d", "--detached", "--no-detached"],
        ["--dry-run"],
        ["-h", "--help"],
        ["-m", "--message string"],
        ["--propose"],
        ["-p", "--prototype"],
        ["--push", "--no-push"],
        ["--stash", "--no-stash"],
        ["--sync", "--no-sync"],
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

    test("compress command", () => {
      const output = new HelpOutput(`
Squash all commits on the current branch down to a single commit.

Compress is a more convenient way of running "git rebase --interactive"
and choosing to fixup all commits.
Branches must be in sync to compress them, run "git sync" as needed.

Provide the --stack switch to compress all branches in the stack.

The compressed commit uses the commit message of the first commit in the branch.
You can provide a custom commit message with the -m switch.

Assuming you have a feature branch with these commits:

$ git log --format='%s'
commit 1
commit 2
commit 3

Let's compress these three commits into a single commit:

$ git town compress

Now your branch has a single commit with the name of the first commit but
containing the changes of all three commits that existed on the branch before:

$ git log --format='%s'
commit 1

Usage:
  git-town compress [flags]

Flags:
      --dry-run          print but do not run the Git commands
  -h, --help             help for compress
  -m, --message string   customize the commit message
      --no-verify        do not run pre-commit hooks
  -s, --stack            Compress the entire stack
  -v, --verbose          display all Git commands run under the hood
`)
      const have = output.flags()
      const want = [
        ["--dry-run"],
        ["-h", "--help"],
        ["-m", "--message string"],
        ["--no-verify"],
        ["-s", "--stack"],
        ["-v", "--verbose"],
      ]
      deepEqual(have, want)
    })
  })
})

suite("Lines", () => {
  suite(".flagLines()", () => {
    test("append command", () => {
      const output = new HelpOutput(appendHelpOutput)
      const have = output.lines().flagLines()
      const want = [
        new FlagLine("      --auto-resolve      auto-resolve phantom merge conflicts"),
        new FlagLine("  -b, --beam              beam some commits from this branch to the new branch"),
        new FlagLine("  -c, --commit            commit the stashed changes into the new branch"),
        new FlagLine("  -d, --detached          don't update the perennial root branch"),
        new FlagLine("      --dry-run           print but do not run the Git commands"),
        new FlagLine("  -h, --help              help for append"),
        new FlagLine("  -m, --message string    the commit message"),
        new FlagLine("      --no-auto-resolve   don't auto-resolve"),
        new FlagLine("      --no-detached       disable detached"),
        new FlagLine("      --no-push           don't push branches"),
        new FlagLine("      --no-stash          don't stash uncommitted changes"),
        new FlagLine("      --no-sync           don't sync branches"),
        new FlagLine("      --propose           propose the new branch"),
        new FlagLine("  -p, --prototype         create a prototype branch"),
        new FlagLine("      --push              push local branches"),
        new FlagLine("      --stash             stash uncommitted changes when creating branches"),
        new FlagLine("      --sync              sync branches (default true)"),
        new FlagLine("  -v, --verbose           display all Git commands run under the hood"),
      ]
      deepEqual(have, want)
    })
  })
})

suite("FlagLine", () => {
  suite(".flags()", () => {
    const tests = {
      "  -b, --beam             description": ["-b", "--beam"],
      "  -d, --display-types string[=\"all\"]   display the branch types": ["-d", "--display-types string"],
    }
    for (const [give, want] of Object.entries(tests)) {
      test(give, () => {
        const flagLine = new FlagLine(give)
        const have = flagLine.flags()
        deepEqual(have, want)
      })
    }
  })
})

suite("replaceValueNotation()", () => {
  const tests = {
    "string[=\"all\"]": "string",
    "string": "string",
  }
  for (const [give, want] of Object.entries(tests)) {
    test(give, () => {
      const have = replaceValueNotation(give)
      deepEqual(have, want)
    })
  }
})

suite("mergeFlags()", () => {
  test("merge flags", () => {
    const give = [["-b", "--beam"], ["-d", "--detached"], ["--no-detached"]]
    const want = [["-b", "--beam"], ["-d", "--detached", "--no-detached"]]
    const have = mergeFlags(give)
    deepEqual(have, want)
  })
})

suite("isNegatedFlagsGroup()", () => {
  const tests = [
    { give: [], want: false },
    { give: ["--no-push"], want: true },
    { give: ["--no-push", "--no-sync"], want: true },
    { give: ["--push"], want: false },
    { give: ["--push", "--no-push"], want: false },
  ]
  for (const { give, want } of tests) {
    test(JSON.stringify(give), () => {
      const have = isNegatedFlagsGroup(give)
      deepEqual(have, want)
    })
  }
})

suite("getPositiveFlagName()", () => {
  const tests = {
    "--no-push": "--push",
    "--no-detached string": "--detached",
    "--no-auto-resolve": "--auto-resolve",
  }
  for (const [give, want] of Object.entries(tests)) {
    test(give, () => {
      const have = getPositiveFlagName(give)
      deepEqual(have, want)
    })
  }
})

suite("matchesFlag()", () => {
  test("exact match", () => {
    deepEqual(matchesFlag("--push", "--push"), true)
  })

  test("flag with value type", () => {
    deepEqual(matchesFlag("--message string", "--message"), true)
  })

  test("different flags", () => {
    deepEqual(matchesFlag("--push", "--sync"), false)
  })

  test("prefix but not value type", () => {
    deepEqual(matchesFlag("--pushall", "--push"), false)
  })
})

suite("findGroupWithPositiveFlag()", () => {
  test("finds matching group", () => {
    const groups = [
      ["-b", "--beam"],
      ["-d", "--detached"],
      ["-p", "--push"],
    ]
    const have = findGroupWithPositiveFlag(groups, "--detached")
    deepEqual(have, ["-d", "--detached"])
  })

  test("finds group with flag with value", () => {
    const groups = [
      ["-m", "--message string"],
      ["-p", "--push"],
    ]
    const have = findGroupWithPositiveFlag(groups, "--message")
    deepEqual(have, ["-m", "--message string"])
  })

  test("no matching group", () => {
    const groups = [
      ["-b", "--beam"],
      ["-d", "--detached"],
    ]
    const have = findGroupWithPositiveFlag(groups, "--push")
    deepEqual(have, undefined)
  })

  test("empty groups", () => {
    const groups: string[][] = []
    const have = findGroupWithPositiveFlag(groups, "--push")
    deepEqual(have, undefined)
  })
})
