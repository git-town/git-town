import { deepEqual } from "node:assert/strict"
import { exec } from "node:child_process"
import { promisify } from "node:util"
import * as textRunner from "text-runner"

const execAsync = promisify(exec)

/** verifies a MD page that describes a Git Town command */
export async function gitTownCommand(action: textRunner.actions.Args) {
  const doc = new Document(action.document)

  // determine the Git Town command that this page describes
  const summarySection = doc.summarySection()
  const command = summarySection.command()

  // get the actual arguments of this Git Town command
  const actualArgs = await command.actualArgs()
  const actualJSON = JSON.stringify(actualArgs, null, 2)

  // get the arguments described by the command summary
  const summaryArgs = summarySection.args()
  const summaryJSON = JSON.stringify(summaryArgs, null, 2)

  // ensure the summary documents the arguments correct
  action.log(`ACTUAL:\n${actualJSON}`)
  action.log(`SUMMARY SECTION:\n${summaryJSON}`)
  deepEqual(summaryArgs, actualArgs)

  // get the arguments described by the "## Options" section
  const optionsArgs = doc.argsInOptions()
  const optionsJSON = JSON.stringify(optionsArgs, null, 2)

  // ensure the options section documents the arguments correct
  action.log(`ACTUAL:\n${actualJSON}`)
  action.log(`OPTIONS SECTION:\n${optionsJSON}`)
  deepEqual(optionsArgs, actualArgs)
}

/** Document contains the AST for an entire webpage describing a Git Town command */
class Document {
  nodes: textRunner.ast.NodeList

  constructor(nodes: textRunner.ast.NodeList) {
    this.nodes = nodes
  }

  /** provides the text of the ```command-summary block at the beginning of this page */
  summarySection(): SummarySection {
    const fences = this.nodes.nodesOfTypes("fence")
    if (fences.length === 0) {
      throw new Error("no fenced blocks found")
    }
    // the first fenced block contains the summary
    const summaryBlock = fences[0]
    const summaryNodes = this.nodes.nodesFor(summaryBlock)
    return new SummarySection(summaryNodes.text())
  }

  /** provides the arguments documented in the "## Options" section of this page */
  argsInOptions(): string[][] {
    let result: string[][] = []
    let insideOptions = false
    for (const node of this.nodes) {
      if (isH2(node)) {
        if (insideOptions) {
          // here we run into the next h2 heading after options --> done parsing options
          return result
        }
        if (this.isOptionsHeading(node)) {
          insideOptions = true
        }
        continue
      }
      if (insideOptions) {
        if (isFlagHeading(node)) {
          const flagNodes = this.nodes.nodesFor(node)
          result.push(texts(flagNodes))
        }
      }
    }
    return result
  }

  /** indicates whether the given node is the "## Options" heading */
  isOptionsHeading(node: textRunner.ast.Node): boolean {
    const nodes = this.nodes.nodesFor(node)
    const text = nodes.text()
    return text === "Options"
  }
}

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
      const variations = normalizedArgText.split("|").map((v) => v.trim())
      result.push(variations)
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

/** GitTownCommand represents a specific Git Town command, like "append" or "sync" */
class GitTownCommand {
  name: string

  constructor(name: string) {
    this.name = name
  }

  /** provides the actual arguments that this Git Town command accepts, determined by calling it with --help and parsing the output */
  async actualArgs(): Promise<string[][]> {
    const result = await execAsync(`git town ${this.name} --help`)
    const output = new HelpOutput(result.stdout)
    return output.flags()
  }
}

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
    }

    // Group negated flags with their positive counterparts
    const groupsToRemove = new Set<number>()

    for (let i = 0; i < result.length; i++) {
      const group = result[i]

      // Check if this is a negated flag (single flag starting with --no-)
      if (group.length !== 1 || !group[0].startsWith('--no-')) {
        continue
      }

      const negatedFlag = group[0]
      const positiveFlagName = negatedFlag.substring(5) // Remove '--no-'

      // Find the group containing the positive counterpart
      for (let j = 0; j < result.length; j++) {
        if (i === j || groupsToRemove.has(j)) {
          continue
        }

        const targetGroup = result[j]
        const hasPositiveFlag = targetGroup.some(flag => flag === `--${positiveFlagName}`)

        if (hasPositiveFlag) {
          // Add the negated flag to this group
          targetGroup.push(negatedFlag)
          groupsToRemove.add(i)
          break
        }
      }
    }

    // Filter out the groups that were merged
    return result.filter((_, index) => !groupsToRemove.has(index))
  }
}

function isFlagHeading(node: textRunner.ast.Node): boolean {
  return node.type === "h4_open"
}

function isH2(node: textRunner.ast.Node): boolean {
  return node.type === "h2_open"
}

function texts(nodes: textRunner.ast.NodeList): string[] {
  let result: string[] = []
  for (const node of nodes) {
    if (node.type === "text") {
      result.push(node.content)
    }
  }
  return standardizeArgument(removeNegatedFlag(result))
}

export function removeNegatedFlag(flags: string[]): string[] {
  if (flags.length < 2) {
    return flags
  }
  return flags.filter(flag => !flag.startsWith("--no-"))
}

export function standardizeArgument(texts: string[]): string[] {
  const result: string[] = []
  for (const text of texts) {
    if (text.startsWith("--")) {
      const parts = text.split(" ")
      if (parts.length > 1) {
        parts[1] = "string"
      }
      result.push(parts.join(" "))
    } else if (text.startsWith("-")) {
      const parts = text.split(" ")
      result.push(parts[0])
    }
  }
  return result
}
