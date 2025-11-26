import { deepEqual, equal } from "node:assert/strict"
import { suite, test } from "node:test"
import { removeNegatedFlag, standardizeArgument } from "../text-runner/command-options.ts"

suite("commandOptions", () => {
  test("removeNegatedFlag", () => {
    const tests = [
      // remove the negated flag
      [
        ["-d", "--detached", "--no-detached"],
        ["-d", "--detached"],
      ],
      // pass through flags without negation
      [
        ["-d", "--detached"],
        ["-d", "--detached"],
      ],
    ]
    for (const [give, want] of tests) {
      test(`${give} -> ${want}`, () => {
        const have = removeNegatedFlag(give)
        deepEqual(have, want)
      })
    }
  })

  test("standardizeArgument", () => {
    const tests = [
      // standardize the argument
      [
        ["-m <msg>", "--message <msg>"],
        ["-m", "--message string"],
      ],
      // work without arguments
      [
        ["-p", "--prototype"],
        ["-p", "--prototype"],
      ],
    ]
    for (const [give, want] of tests) {
      test(`${give} -> ${want}`, () => {
        const have = standardizeArgument(give)
        deepEqual(have, want)
      })
    }
  })
})
