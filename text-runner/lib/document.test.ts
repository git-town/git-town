import { deepEqual } from "node:assert/strict"
import { suite, test } from "node:test"
import { standardizeArgument } from "./document.ts"

suite("Document", () => {
  suite("standardizeArgument()", () => {
    const tests = [
      {
        desc: "string argument",
        give: ["-m <msg>", "--message <msg>"],
        want: ["-m", "--message string"],
      },
      {
        desc: "int argument",
        give: ["-d int", "--down int"],
        want: ["-d", "--down int"],
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
})
