import { deepEqual } from "node:assert/strict"
import * as textRunner from "text-runner"
import { Document } from "./lib/document.ts"

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
