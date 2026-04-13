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
        const inner = groupMatch[1]
        const rest = groupMatch[2]
        // [(--non)-interactive] is stored as (--non) + -interactive; keep the (non)/(no)
        // template so splitNegations can expand to --interactive and --non-interactive.
        if (inner === "--no" && rest.startsWith("-")) {
          argText = `--(no)-${rest.slice(1)}`
        } else if (inner === "--non" && rest.startsWith("-")) {
          argText = `--(non)-${rest.slice(1)}`
        } else {
          argText = inner + rest
        }
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
    const negationKind = negationKindForVariation(variation)
    if (negationKind !== null) {
      result.push(...splitNegation(variation, negationKind))
    } else {
      result.push(variation)
    }
  }
  return result
}

function negationKindForVariation(variation: string): "no" | "non" | null {
  if (variation.startsWith("--(non)-")) {
    return "non"
  }
  if (variation.startsWith("--(no)-")) {
    return "no"
  }
  return null
}

export function isNegatable(variation: string): boolean {
  return negationKindForVariation(variation) !== null
}

export function splitNegation(variation: string, negation: string): string[] {
  const result: string[] = []
  const name = variationName(variation)
  result.push(`--${name}`)
  result.push(`--${negation}-${name}`)
  return result
}

export function variationName(variation: string): string {
  if (variation.startsWith("--(non)-")) {
    return variation.slice("--(non)-".length)
  }
  if (variation.startsWith("--(no)-")) {
    return variation.slice("--(no)-".length)
  }
  throw new Error(`variation is not negatable: ${variation}`)
}
