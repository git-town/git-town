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
    return result
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
    return flagsPart.split(/,\s+/).map(flag => flag.replace(/\[="[^"]*"\]/, ""))
  }
}
