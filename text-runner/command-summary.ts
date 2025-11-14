import * as textRunner from "text-runner";

export function commandSummary(action: textRunner.actions.Args) {
  // git town append <branch-name> [-p | --prototype] [-d | --detached] [-c | --commit] [-m | --message <message>] [--propose] [--dry-run] [-v | --verbose]
  const text = action.region.text();
  const command = extractCommand(text);
  const args = extractArgs(text);
}

export function extractCommand(text: string): string {
  const match = text.match(/^git town (\w+)/);
  return match?.[1] ?? "";
}

export function extractArgs(text: string): string[] {
  // TODO: implement this
}
