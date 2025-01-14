// @ts-check
/// <reference types="node" />

/**
 * This file is a linter that verifies that the `printConfig` function
 * references every field in the `NormalConfigData` struct.
 */

import fs from "node:fs/promises";

const NORMAL_CONFIG_DATA_PATH = "internal/config/configdomain/normal_config_data.go";
const PRINT_CONFIG_PATH = "internal/cmd/config/root.go";
const ACCEPTABLE_MISSING_FIELDS = [
  "Aliases",
  "BranchTypeOverrides",
];

const normalConfigDataFile = await fs.readFile(NORMAL_CONFIG_DATA_PATH, "utf8");
const printConfigFile = await fs.readFile(PRINT_CONFIG_PATH, "utf8");

const normalConfigDataFieldsMatch = normalConfigDataFile.match(/type NormalConfigData struct {([^}]*)}/s);
if (!normalConfigDataFieldsMatch) {
  throw new Error("Failed to find NormalConfigData struct");
}
const normalConfigDataFields = normalConfigDataFieldsMatch[1].trim().split("\n").map((line) =>
  line.trim().split(/\s+/)[0]
).filter(field => !field.startsWith("//"));

const printConfigFunctionMatch = printConfigFile.match(/func printConfig\(.*?\) {([^}]*)}/s);
if (!printConfigFunctionMatch) {
  throw new Error("Failed to find printConfig function");
}
const printConfigFunctionBody = printConfigFunctionMatch[1];

const missingFields = normalConfigDataFields.filter((field) =>
  !printConfigFunctionBody.includes(field) && !ACCEPTABLE_MISSING_FIELDS.includes(field)
);
if (missingFields.length > 0) {
  console.error("Missing fields in printConfig function:");
  console.log(missingFields.join("\n"));
  process.exit(1);
}
