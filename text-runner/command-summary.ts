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

export function extractArgs(text: string): string[][] {
  const args: string[][] = [];
  // Match all optional arguments in square brackets: [-p | --prototype]
  const matches = text.matchAll(/\[([^\]]+)\]/g);

  for (const match of matches) {
    const argText = match[1];
    // Split by | to get the different variations of the flag
    const variations = argText.split("|").map((v) => v.trim());
    args.push(variations);
  }

  return args;
}
