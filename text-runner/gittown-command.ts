import { deepEqual } from "node:assert/strict"
import { exec } from "node:child_process"
import { promisify } from "node:util"
import * as textRunner from "text-runner"

const execAsync = promisify(exec)

export async function gittownCommand(action: textRunner.actions.Args) {
  const doc = action.document

  // verify the command summary
  const summary = findCommandSummary(doc)
  const command = extractCommand(summary)
  const actualArgs = await commandArgs(command)
  const documentedArgs = extractArgs(summary)
  const summaryJSON = JSON.stringify(documentedArgs, null, 2)
  const actualJSON = JSON.stringify(actualArgs, null, 2)
  if (summaryJSON !== actualJSON) {
    action.log(`ACTUAL:\n${actualJSON}`)
    action.log(`SUMMARY:\n${summaryJSON}`)
    deepEqual(documentedArgs, actualArgs)
  }

  // verify the command options
  const options = findOptions(doc)
  const optionsJSON = JSON.stringify(options, null, 2)
  if (summaryJSON !== optionsJSON) {
    action.log(`ACTUAL:\n${actualJSON}`)
    action.log(`BODY:\n${optionsJSON}`)
    deepEqual(summaryJSON, optionsJSON)
  }
}

export function extractCommand(text: string): string {
  const match = text.match(/^git town ([^<[(]+?)(?:\s+-|\s+<|\s+\[|\s+\(|$)/)
  return match?.[1]?.trim() ?? ""
}

function findCommandSummary(doc: textRunner.ast.NodeList): string {
  const fences = doc.nodesOfTypes("fence")
  if (fences.length === 0) {
    throw new Error("no command summary found")
  }
  const fence = fences[0]
  const fenceNodes = doc.nodesFor(fence)
  const fenceText = fenceNodes.text()
  return fenceText
}

function findOptions(doc: textRunner.ast.NodeList): string[][] {
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
  return []
}

async function commandArgs(command: string): Promise<string[][]> {
  const output = await commandHelp(command)
  return parseCommandHelpOutput(output)
}

async function commandHelp(command: string): Promise<string> {
  const result = await execAsync(`git town ${command} --help`)
  return result.stdout
}

export function extractArgs(text: string): string[][] {
  const args: string[][] = []
  // Match all optional arguments in square brackets: [-p | --prototype] or [(-m | --message) <text>]
  const matches = text.matchAll(/\[([^\]]+)\]/g)
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

export function parseCommandHelpOutput(help: string): string[][] {
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
