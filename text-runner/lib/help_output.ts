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

export function isNegatedFlagsGroup(flags: string[]): boolean {
  return flags.length > 0 && flags.every(flag => flag.startsWith("--no-"))
}

export function getPositiveFlagName(negatedFlag: string): string {
  const baseName = negatedFlag.substring(5).split(" ")[0]
  return "--" + baseName
}

export function matchesPositiveFlag(flag: string, positiveFlag: string): boolean {
  return flag === positiveFlag || flag.startsWith(positiveFlag + " ")
}

export function findGroupWithPositiveFlag(result: string[][], positiveFlag: string): string[] | undefined {
  return result.find(group => group.some(flag => matchesPositiveFlag(flag, positiveFlag)))
}

function extractNegated(flags: string[][]): { normal: string[][]; negated: string[][] } {
  const negated: string[][] = []
  const normal: string[][] = []

  for (const currentFlags of flags) {
    if (isNegatedFlagsGroup(currentFlags)) {
      negated.push(currentFlags)
    } else {
      normal.push([...currentFlags])
    }
  }

  return { normal, negated }
}

export function mergeFlags(flags: string[][]): string[][] {
  const { normal, negated } = extractNegated(flags)

  // merge negated flags with their positive counterparts
  for (const negatedFlags of negated) {
    const positiveFlag = getPositiveFlagName(negatedFlags[0])
    const targetGroup = findGroupWithPositiveFlag(normal, positiveFlag)
    if (targetGroup) {
      targetGroup.push(...negatedFlags)
    } else {
      normal.push([...negatedFlags])
    }
  }

  return normal
}
