import { deepEqual } from "node:assert/strict"
import * as textRunner from "text-runner"
import { extractArgs } from "./command-summary.ts"

export function commandOptions(action: textRunner.actions.Args) {
  const doc = action.document
  const summary = findCommandSummary(doc)
  const options = findOptions(doc)
  const summaryJSON = JSON.stringify(summary, null, 2)
  const optionsJSON = JSON.stringify(options, null, 2)
  if (summaryJSON !== optionsJSON) {
    action.log(`SUMMARY:\n${summaryJSON}`)
    action.log(`BODY:\n${optionsJSON}`)
    deepEqual(summaryJSON, optionsJSON)
  }
}

function findCommandSummary(doc: textRunner.ast.NodeList): string[][] {
  const fences = doc.nodesOfTypes("fence")
  if (fences.length === 0) {
    throw new Error("no command summary found")
  }
  const fence = fences[0]
  const fenceNodes = doc.nodesFor(fence)
  const fenceText = fenceNodes.text()
  return extractArgs(fenceText)
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
