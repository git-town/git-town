import { equal } from "node:assert/strict";
import { test } from "node:test";
import { extractCommand } from "./command-summary.ts";

test("extractCommand extracts command name from valid git town command", () => {
  equal(extractCommand("git town append"), "append");
  equal(extractCommand("git town sync"), "sync");
  equal(extractCommand("git town ship"), "ship");
});

test("extractCommand extracts command with additional arguments", () => {
  equal(
    extractCommand("git town append <branch-name> [-p | --prototype]"),
    "append",
  );
  equal(
    extractCommand("git town sync --all"),
    "sync",
  );
});

test("extractCommand returns empty string when no match", () => {
  equal(extractCommand("not a git town command"), "");
  equal(extractCommand("git town"), "");
  equal(extractCommand(""), "");
  equal(extractCommand("git town "), "");
});

test("extractCommand handles commands with special characters", () => {
  // \w+ matches word characters, so it stops at the hyphen
  equal(extractCommand("git town append-branch"), "append");
  // \w+ matches digits as well
  equal(extractCommand("git town 123"), "123");
});
