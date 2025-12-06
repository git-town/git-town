import { GitTownCommand } from "./command"

/** SummarySection contains the text of the ```command-summary block of a Document*/
export class SummarySection {
  text: string

  constructor(text: string) {
    this.text = text
  }

  /** provides the arguments that this summary section describes for its Git Town command */
  args(): string[][] {
    const args: string[][] = []
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
      const variations = normalizedArgText.split("|").map((v) => v.trim())
      args.push(variations)
    }
    return args
  }

  /** provides the name of the Git Town command described by this summary section */
  command(): GitTownCommand {
    const match = this.text.match(/^git town ([^<[(]+?)(?:\s+-|\s+<|\s+\[|\s+\(|$)/)
    const commandName = match?.[1]?.trim() || ""
    return new GitTownCommand(commandName)
  }
}
