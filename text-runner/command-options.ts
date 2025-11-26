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
    if (node.type === "h2_open") {
      if (insideOptions) {
        // here we run into the next h2 heading after options --> done parsing options
        return result
      }
      if (isOptionsNode(node, doc)) {
        insideOptions = true
        continue
      }
    }
    if (insideOptions) {
      if (isFlagNode(node, doc)) {
        const h4Nodes = doc.nodesFor(node)
        result.push(texts(h4Nodes))
      }
    }
  }
  return []
}

function isFlagNode(node: textRunner.ast.Node, doc: textRunner.ast.NodeList): boolean {
  return node.type === "h4_open"
}

function isOptionsNode(node: textRunner.ast.Node, doc: textRunner.ast.NodeList): boolean {
  const h2Nodes = doc.nodesFor(node)
  const h2Text = h2Nodes.text()
  return h2Text === "Options"
}

function texts(nodes: textRunner.ast.NodeList): string[] {
  let result: string[] = []
  for (const node of nodes) {
    if (node.type === "text") {
      result.push(node.content)
    }
  }
  return removeNegatedFlag(result)
}

function removeNegatedFlag(flags: string[]): string[] {
  return flags.filter((flag) => !flag.startsWith("--no-"))
}
