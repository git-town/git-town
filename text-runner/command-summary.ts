import * as textRunner from "text-runner";

export function commandSummary(action: textRunner.actions.Args) {
  console.log("This is the implementation of the 'command-summary' action.");
  console.log("Text inside the semantic document region:", action.region.text());
  console.log("For more information see");
  console.log("https://github.com/kevgo/text-runner/blob/main/documentation/user-defined-actions.md");
}
