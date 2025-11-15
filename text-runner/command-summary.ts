import { exec } from "node:child_process"
import { promisify } from "node:util"
import * as textRunner from "text-runner"

const execAsync = promisify(exec)

export async function commandSummary(action: textRunner.actions.Args) {
  const text = action.region.text()
  const command = extractCommand(text)
  const documentedArgs = extractArgs(text)
  const actualArgs = await commandArgs(command)
}

export function extractCommand(text: string): string {
  const match = text.match(/^git town (\w+)/)
  return match?.[1] ?? ""
}

export function extractArgs(text: string): string[][] {
  const args: string[][] = []
  // Match all optional arguments in square brackets: [-p | --prototype]
  const matches = text.matchAll(/\[([^\]]+)\]/g)
  for (const match of matches) {
    const argText = match[1]
    // Split by | to get the different variations of the flag
    const variations = argText.split("|").map((v) => v.trim())
    args.push(variations)
  }
  return args
}

async function commandArgs(command: string): Promise<string[][]> {
  // run the command with --help and parse the output
  const output = await commandHelp(command)
  return parseCommandHelpOutput(output)
}

async function commandHelp(command: string): Promise<string> {
  const result = await execAsync(`git town ${command} --help`)
  return result.stdout
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
      const flags = flagsPart.split(/,\s+/)
      if (flags.length > 0) {
        result.push(flags)
      }
    }
  }
  return result
}
