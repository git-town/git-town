/** HelpOutput is the output of a Git Town command executed with "--help" */
export class HelpOutput {
  text: string

  constructor(text: string) {
    this.text = text
  }

  /** provides the CLI flags described in this help output */
  flags(): string[][] {
    const result: string[][] = []
    for (const flagLine of this.lines().flagLines()) {
      result.push(flagLine.flags())
    }
    return mergeFlags(result)
  }

  lines(): Lines {
    const result: Line[] = []
    for (const line of this.text.split("\n")) {
      result.push(new Line(line))
    }
    return new Lines(result)
  }
}

/** Lines represents all lines in the help output for a Git Town command */
class Lines {
  lines: Line[]

  constructor(lines: Line[]) {
    this.lines = lines
  }

  /** flagsLines provides the lines of the "Flags:" section */
  flagLines(): FlagLine[] {
    let result: FlagLine[] = []
    let inFlagsSection = false
    for (const line of this.lines) {
      if (line.isStartOfFlagsSection()) {
        inFlagsSection = true
        continue
      }
      if (!inFlagsSection) {
        continue
      }
      if (line.isEndOfFlagsSection()) {
        break
      }
      result.push(new FlagLine(line.text))
    }
    return result
  }
}

/** Line is a line in the help output of a Git Town command */
class Line {
  text: string
  constructor(text: string) {
    this.text = text
  }

  isStartOfFlagsSection(): boolean {
    return this.text.includes("Flags:")
  }

  isEndOfFlagsSection(): boolean {
    return this.text.trim() === ""
  }
}

/** FlagLine is a line in the "Flags:" section of the help output of a Git Town command */
export class FlagLine {
  text: string

  constructor(text: string) {
    this.text = text
  }

  /** flags provides the flags that this FlagLine defines */
  flags(): string[] {
    // Parse flag line - format: "  -b, --beam             description"
    // The description starts after 2 or more spaces
    const match = this.text.match(/^\s+(.+?)\s{2,}/)
    if (!match) {
      return []
    }
    const flagsPart = match[1].trim()
    // Remove default value notation like [="all"]
    return flagsPart.split(/,\s+/).map(replaceValueNotation)
  }
}

export function replaceValueNotation(flag: string): string {
  return flag.replace(/\[="[^"]*"\]/, "")
}

export function mergeFlags(flags: string[][]): string[][] {
  const result: string[][] = []

  for (const currentFlags of flags) {
    // Check if this group contains only negated flags (flags starting with --no-)
    const allNegated = currentFlags.every(flag => flag.startsWith('--no-'))

    if (allNegated && currentFlags.length > 0) {
      // Extract the base flag name (without --no- prefix and without value types)
      const negatedFlag = currentFlags[0]
      const baseName = negatedFlag.substring(5).split(' ')[0] // Remove '--no-' and any value type
      const positiveFlag = '--' + baseName

      // Find in result which group contains the positive flag
      let foundInResult = false
      for (const resultGroup of result) {
        // Check if any flag in the result group matches the positive flag (with or without value)
        if (resultGroup.some(flag => flag === positiveFlag || flag.startsWith(positiveFlag + ' '))) {
          // Add all negated flags from current group to this result group
          resultGroup.push(...currentFlags)
          foundInResult = true
          break
        }
      }

      if (!foundInResult) {
        // No matching positive flag found, add as new group
        result.push([...currentFlags])
      }
    } else {
      // Not a negated-only group, add as is
      result.push([...currentFlags])
    }
  }

  return result
}
