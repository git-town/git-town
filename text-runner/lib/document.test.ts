import { deepEqual } from "node:assert/strict"
import { suite, test } from "node:test"
import { removeNegatedFlag, standardizeArgument } from "./document.ts"

suite("removeNegatedFlag", () => {
  const tests = [
    {
      desc: "has negated flag",
      give: ["-d", "--detached", "--no-detached"],
      want: ["-d", "--detached"],
    },
    {
      desc: "no negated flag",
      give: ["-d", "--detached"],
      want: ["-d", "--detached"],
    },
    {
      desc: "allows a single negated flag",
      give: ["--no-detached"],
      want: ["--no-detached"],
    },
  ]
  for (const { desc, give, want } of tests) {
    test(desc, () => {
      const have = removeNegatedFlag(give)
      deepEqual(have, want)
    })
  }
})

suite("standardizeArgument", () => {
  const tests = [
    {
      desc: "has argument",
      give: ["-m <msg>", "--message <msg>"],
      want: ["-m", "--message string"],
    },
    {
      desc: "no argument",
      give: ["-p", "--prototype"],
      want: ["-p", "--prototype"],
    },
  ]
  for (const { desc, give, want } of tests) {
    test(desc, () => {
      const have = standardizeArgument(give)
      deepEqual(have, want)
    })
  }
})
