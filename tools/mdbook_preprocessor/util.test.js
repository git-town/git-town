// @ts-check
/// <reference types="node" />

import assert from "node:assert/strict"
import { describe, it } from "node:test"

import { extractCommand, tokenize } from "./util.js"

describe("tokenize", () => {
  it("should split string by spaces", () => {
    const have = "git town append"
    const want = ["git", "town", "append"]
    assert.deepEqual(tokenize(have), want)
  })

  it("should not split tokens wrapped by various brackets", () => {
    for (const [left, right] of [["(", ")"], ["<", ">"], ["[", "]"]]) {
      const have = `git town append ${left}-p | --prototype${right}`
      const want = ["git", "town", "append", `${left}-p | --prototype${right}`]
      assert.deepEqual(tokenize(have), want)
    }
  })

  it("should not split tokens wrapped by nested brackets", () => {
    const have = "git town hack [<branch name>]"
    const want = ["git", "town", "hack", "[<branch name>]"]
    assert.deepEqual(tokenize(have), want)
  })
})
