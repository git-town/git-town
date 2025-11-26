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
  return extractArgs(fenceText)
}

function findOptions(doc: textRunner.ast.NodeList): string[][] {
  let result: string[][] = []
  let insideOptions = false
  for (const node of doc) {
    if (node.type === "h2_open") {
      if (insideOptions) {
        console.log("done parsing options")
        // here we run into the next h2 heading after options --> done parsing options
        return result
      }
      if (isOptionsNode(node, doc)) {
        console.log("found options")
        insideOptions = true
        continue
      }
    }
  }
  return []
}

function isOptionsNode(node: textRunner.ast.Node, doc: textRunner.ast.NodeList): boolean {
  const h2Nodes = doc.nodesFor(node)
  const h2Text = h2Nodes.text()
  return h2Text === "Options"
}
