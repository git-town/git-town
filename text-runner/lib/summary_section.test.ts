import { deepEqual } from "node:assert/strict"
import { suite, test } from "node:test"
import { isNegatable, splitNegation, splitNegations, SummarySection, variationName } from "./summary_section.ts"

suite("SummarySection", () => {
  suite(".args()", () => {
    const tests = [
      {
        desc: "append command",
        give:
          "git town append <branch-name> [--(no)-auto-resolve] [-b | --beam] [-c | --commit] [-d | --(no)-detached] [--dry-run] [-h | --help] [(-m | --message) <message>] [--propose] [-p | --prototype] [--(no)-push] [--(no)-stash] [--(no)-sync] [-v | --verbose]",
        want: [
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

suite("negations", () => {
  suite("isNegatable", () => {
    const tests = {
      "--(no)-detach": true,
      "--beam": false,
    }
    for (const [give, want] of Object.entries(tests)) {
      test(give, () => {
        const have = isNegatable(give)
        deepEqual(have, want)
      })
    }
  })

  suite("variationName", () => {
    const tests = {
      "--(no)-detach": "detach",
    }
    for (const [give, want] of Object.entries(tests)) {
      test(give, () => {
        const have = variationName(give)
        deepEqual(have, want)
      })
    }
  })

  suite("splitNegation", () => {
    const tests = {
      "--(no)-detach": ["--detach", "--no-detach"],
    }
    for (const [give, want] of Object.entries(tests)) {
      test(give, () => {
        const have = splitNegation(give)
        deepEqual(have, want)
      })
    }
  })

  suite("splitNegations", () => {
    const tests = [
      {
        desc: "negatable",
        give: ["-d", "--(no)-detach"],
        want: ["-d", "--detach", "--no-detach"],
      },
      {
        desc: "not negatable",
        give: ["-b", "--beam"],
        want: ["-b", "--beam"],
      },
    ]
    for (const { desc, give, want } of tests) {
      test(desc, () => {
        const have = splitNegations(give)
        deepEqual(have, want)
      })
    }
  })
})
