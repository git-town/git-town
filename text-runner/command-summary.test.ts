import { deepEqual, equal } from "node:assert/strict";
import { suite, test } from "node:test";
import { extractArgs, extractCommand } from "./command-summary.ts";

suite("extractCommand", () => {
  const tests = {
    "git town append": "append",
    "git town sync": "sync",
    "git town ship": "ship",
    "git town append <branch-name> [-p | --prototype]": "append",
    "git town sync --all": "sync",
    "git town ship --all": "ship",
  };
  for (const [give, want] of Object.entries(tests)) {
    test(`${give} -> ${want}`, () => {
      equal(extractCommand(give), want);
    });
  }
});

suite("extractArgs", () => {
  const tests = {
    "git town append <branch-name> [-p | --prototype] [-d | --detached] [-c | --commit] [-m | --message <message>] [--propose] [--dry-run] [-v | --verbose]":
      [
        ["-p", "--prototype"],
        ["-d", "--detached"],
        ["-c", "--commit"],
        ["-m", "--message <message>"],
        ["--propose"],
        ["--dry-run"],
        ["-v", "--verbose"],
      ],
  };
  for (const [give, want] of Object.entries(tests)) {
    test(`${give} -> ${want}`, () => {
      deepEqual(extractArgs(give), want);
    });
  }
});
