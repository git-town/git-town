import { GitTownCommand } from "./command.ts"

/** SummarySection contains the text of the ```command-summary block of a Document*/
export class SummarySection {
  text: string

  constructor(text: string) {
    this.text = text
  }

  /** provides the arguments that this summary section describes for its Git Town command */
  args(): string[][] {
    const result: string[][] = []
    // Match all optional arguments in square brackets: [-p | --prototype] or [(-m | --message) <text>]
    const matches = this.text.matchAll(/\[([^\]]+)\]/g)
    for (const match of matches) {
      let argText = match[1]

      // Check if this contains grouped arguments in parentheses
      const groupMatch = argText.match(/^\(([^)]+)\)(.*)/)
      if (groupMatch) {
        // Extract the content inside parentheses and any content after (like <message>)
        argText = groupMatch[1] + groupMatch[2]
      }

      if (!argText.trim().startsWith("-")) {
        // this element doesn't contain a flag (doesn't start with -)
        continue
      }
      const normalizedArgText = argText.replace(/<.+?>/g, "string")
      // Split by | to get the different variations of the flag
      const variations = normalizedArgText.split("|").map(v => v.trim())
      // expand --(no)-foo into --foo and --no-foo
      const expanded = splitNegations(variations)
      result.push(expanded)
    }
    return result
  }

  /** provides the name of the Git Town command described by this summary section */
  command(): GitTownCommand {
    const match = this.text.match(/^git town ([^<[(]+?)(?:\s+-|\s+<|\s+\[|\s+\(|$)/)
    const commandName = match?.[1]?.trim() || ""
    return new GitTownCommand(commandName)
  }
}

export function splitNegations(variations: string[]): string[] {
  const result: string[] = []
  for (const variation of variations) {
    if (isNegatable(variation)) {
      result.push(...splitNegation(variation))
    } else {
      result.push(variation)
    }
  }
  return result
}

export function isNegatable(variation: string): boolean {
  return variation.startsWith("--(no)-")
}

export function splitNegation(variation: string): string[] {
  const result: string[] = []
  const name = variationName(variation)
  result.push(`--${name}`)
  result.push(`--no-${name}`)
  return result
}

export function variationName(variation: string): string {
  return variation.substring(7)
}
