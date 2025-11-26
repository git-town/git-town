import * as textRunner from "text-runner"
import { extractArgs } from "./command-summary.ts"

export function commandOptions(action: textRunner.actions.Args) {
  const doc = action.document
  const summary = findCommandSummary(doc)
  const options = findOptions(doc)
}

function findCommandSummary(doc: textRunner.ast.NodeList): string[][] {
  const fences = doc.nodesOfTypes("fence")
  if (fences.length === 0) {
    throw new Error("no command summary found")
  }
  const fence = fences[0]
  const fenceNodes = doc.nodesFor(fence)
  const fenceText = fenceNodes.text()
  const summaryFlags = extractArgs(fenceText)
  return summaryFlags
}

function findOptions(doc: textRunner.ast.NodeList): string[] {
  let insideOptions = false
  for (const node of doc) {
    if (node.type === "h2_open") {
      if (isOptionsNode(node, doc)) {
        console.log("found options")
        insideOptions = true
        continue
      }
      // here we run into the next h2 heading after options
      insideOptions = false
      continue
    }
  }
  return []
}

function isOptionsNode(node: textRunner.ast.Node, doc: textRunner.ast.NodeList): boolean {
  const h2Nodes = doc.nodesFor(node)
  const h2Text = h2Nodes.text()
  return h2Text === "Options"
}
