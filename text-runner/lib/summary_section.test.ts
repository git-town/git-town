import { deepEqual } from "node:assert/strict"
import { suite, test } from "node:test"
import { SummarySection } from "./summary_section.ts"

suite("SummarySection", () => {
  suite(".args()", () => {
    const tests = [
      {
        desc: "append command",
        give:
          "git town append <branch-name> [--auto-resolve] [-b | --beam] [-c | --commit] [-d | --detached] [--dry-run] [-h | --help] [(-m | --message) <message>] [--propose] [-p | --prototype] [--push] [--stash] [--sync] [-v | --verbose]",
        want: [
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
        ],
      },
      {
        desc: "completions command",
        give: "git town completions (bash|fish|powershell|zsh) [--no-descriptions] [-h | --help]",
        want: [
          ["--no-descriptions"],
          ["-h", "--help"],
        ],
      },
      {
        desc: "config get-parent command",
        give: "git town config get-parent [<branch-name>] [-v | --verbose] [-h | --help]",
        want: [
          ["-v", "--verbose"],
          ["-h", "--help"],
        ],
      },
    ]
    for (const { desc, give, want } of tests) {
      test(desc, () => {
        const summarySection = new SummarySection(give)
        const have = summarySection.args()
        deepEqual(have, want)
      })
    }
  })

  suite(".command()", () => {
    const tests = {
      "git town append": "append",
      "git town config get-parent": "config get-parent",
      "git town sync": "sync",
      "git town ship": "ship",
      "git town append <branch-name> [-p | --prototype]": "append",
      "git town sync --all": "sync",
      "git town ship --all": "ship",
      "git town completions (bash|fish|powershell|zsh) [--no-descriptions] [-h | --help]": "completions",
      "git town config get-parent [<branch-name>] [-v | --verbose] [-h | --help]": "config get-parent",
    }
    for (const [give, want] of Object.entries(tests)) {
      test(give, () => {
        const summarySection = new SummarySection(give)
        const have = summarySection.command().name
        deepEqual(have, want)
      })
    }
  })
})
