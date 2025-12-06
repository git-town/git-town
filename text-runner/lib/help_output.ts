/** HelpOutput is the output of a Git Town command executed with "--help" */
export class HelpOutput {
  text: string

  constructor(text: string) {
    this.text = text
  }

  /** provides the content of the "Flags:" section of this help output as a list of flag variations */
  flags(): string[][] {
    const result: string[][] = []
    const lines = this.text.split("\n")
    let inFlagsSection = false
    for (const line of lines) {
      if (line.includes("Flags:")) {
        inFlagsSection = true
        continue
      }
      if (!inFlagsSection) {
        continue
      }
      // Stop if we hit an empty line (end of flags section)
      if (line.trim() === "") {
        break
      }
      const flags = this.flagLine(line)
      result.push(...flags)
    }
    return result
  }

  flagLine(line: string): string[][] {
    const result: string[][] = []
    // Parse flag line - format: "  -b, --beam             description"
    // The description starts after 2 or more spaces
    const match = line.match(/^\s+(.+?)\s{2,}/)
    if (match) {
      const flagsPart = match[1].trim()
      const flags = flagsPart.split(/,\s+/).map((flag) => {
        // Remove default value notation like [="all"]
        return flag.replace(/\[="[^"]*"\]/, "")
      })
      if (flags.length > 0) {
        result.push(flags)
      }
    }
    return result
  }
}
