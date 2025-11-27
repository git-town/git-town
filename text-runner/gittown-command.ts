import { deepEqual } from "node:assert/strict"
import { exec } from "node:child_process"
import { promisify } from "node:util"
import * as textRunner from "text-runner"

const execAsync = promisify(exec)

/** verifies a MD page that describes a Git Town command */
export async function gittownCommand(action: textRunner.actions.Args) {
  const doc = new Document(action.document)

  // determine the Git Town command that this page describes
  const summarySection = doc.summarySection()
  const command = summarySection.command()

  // get the actual arguments of this Git Town command
  const actualArgs = await command.loadArgs()
  const actualJSON = JSON.stringify(actualArgs, null, 2)

  // get the arguments described by the command summary
  const summaryArgs = summarySection.args()
  const summaryJSON = JSON.stringify(summaryArgs, null, 2)

  // ensure the summary documents the arguments correct
  if (summaryJSON !== actualJSON) {
    action.log(`ACTUAL:\n${actualJSON}`)
    action.log(`SUMMARY SECTION:\n${summaryJSON}`)
    deepEqual(summaryArgs, actualArgs)
  }

  // get the arguments described by the "## Options" section
  const optionsArgs = findArgsInOptionsSection(action.document)
  const optionsJSON = JSON.stringify(optionsArgs, null, 2)

  // ensure the options section documents the arguments correct
  if (summaryJSON !== optionsJSON) {
    action.log(`ACTUAL:\n${actualJSON}`)
    action.log(`OPTIONS SECTION:\n${optionsJSON}`)
    deepEqual(optionsJSON, actualJSON)
  }
}

/** Document represents the AST for the entire document of a page describing a Git Town command */
class Document {
  nodes: textRunner.ast.NodeList

  constructor(nodes: textRunner.ast.NodeList) {
    this.nodes = nodes
  }

  /** provides the text of the command summary section */
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
}

/** SummarySection is the text in the ```command-summary``` block at the beginning of a page describing a Git Town command */
export class SummarySection {
  text: string

  constructor(text: string) {
    this.text = text
  }

  /** provides the arguments described in this summary section */
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

  /** provides the name of the Git Town command described by the given summary text */
  command(): GitTownCommand {
    const match = this.text.match(/^git town ([^<[(]+?)(?:\s+-|\s+<|\s+\[|\s+\(|$)/)
    const commandName = match?.[1]?.trim() || ""
    return new GitTownCommand(commandName)
  }
}

/** GitTownCommand represents a specific Git Town command, like "append" or "sync" */
export class GitTownCommand {
  name: string

  constructor(name: string) {
    this.name = name
  }

  /** provides the actual arguments of the command, as reported by calling the command with --help */
  async loadArgs(): Promise<string[][]> {
    const output = await this.runCommandHelp(this.name)
    return this.parseHelpOutput(output)
  }

  parseHelpOutput(help: string): string[][] {
    const result: string[][] = []
    const lines = help.split("\n")
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
    return result
  }

  /** calls the command with "--help" on the CLI and provides the output */
  async runCommandHelp(command: string): Promise<string> {
    const result = await execAsync(`git town ${command} --help`)
    return result.stdout
  }
}

/** provides the options documented in the page body, under the "## Options" tag */
function findArgsInOptionsSection(doc: textRunner.ast.NodeList): string[][] {
  let result: string[][] = []
  let insideOptions = false
  for (const node of doc) {
    if (isH2(node)) {
      if (insideOptions) {
        // here we run into the next h2 heading after options --> done parsing options
        return result
      }
      if (isOptionsHeading(node, doc)) {
        insideOptions = true
      }
      continue
    }
    if (insideOptions) {
      if (isFlagHeading(node, doc)) {
        const flagNodes = doc.nodesFor(node)
        result.push(texts(flagNodes))
      }
    }
  }
  return result
}

function isFlagHeading(node: textRunner.ast.Node, doc: textRunner.ast.NodeList): boolean {
  return node.type === "h4_open"
}

function isH2(node: textRunner.ast.Node): boolean {
  return node.type === "h2_open"
}

function isOptionsHeading(node: textRunner.ast.Node, doc: textRunner.ast.NodeList): boolean {
  const nodes = doc.nodesFor(node)
  const text = nodes.text()
  return text === "Options"
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
  return flags.filter((flag) => !flag.startsWith("--no-"))
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
