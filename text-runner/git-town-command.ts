import { deepEqual } from "node:assert/strict"
import { exec } from "node:child_process"
import { promisify } from "node:util"
import * as textRunner from "text-runner"
import { SummarySection } from "./lib/summary_section"

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

function isFlagHeading(node: textRunner.ast.Node): boolean {
  return node.type === "h4_open"
}

function isH2(node: textRunner.ast.Node): boolean {
  return node.type === "h2_open"
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
